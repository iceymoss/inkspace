package uploader

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStorageLifecycle(t *testing.T) {
	storage, err := NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	key := "documents/hello.txt"
	content := "hello storage"

	if err := storage.Put(ctx, key, strings.NewReader(content), int64(len(content)), "text/plain"); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	info, err := storage.Stat(ctx, key)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Key != key || info.Size != int64(len(content)) || !strings.HasPrefix(info.ContentType, "text/plain") {
		t.Fatalf("Stat() = %#v", info)
	}
	reader, err := storage.Open(ctx, key)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	got, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != content {
		t.Fatalf("Open() content = %q, want %q", got, content)
	}
	if err := storage.Delete(ctx, key); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := storage.Stat(ctx, key); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Stat() after Delete error = %v, want os.ErrNotExist", err)
	}
}

func TestLocalStorageRejectsEscapingPaths(t *testing.T) {
	root := t.TempDir()
	storage, err := NewLocalStorage(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, key := range []string{"../outside.txt", filepath.Join(root, "absolute.txt")} {
		if err := storage.Put(ctx, key, strings.NewReader("bad"), 3, "text/plain"); err == nil {
			t.Errorf("Put(%q) error = nil, want path validation error", key)
		}
	}

	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	if err := storage.Put(ctx, "escape/outside.txt", strings.NewReader("bad"), 3, "text/plain"); err == nil {
		t.Fatal("Put() through escaping symlink error = nil")
	}
}

func TestLocalStorageFailedPutPreservesExistingObject(t *testing.T) {
	storage, err := NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := storage.Put(ctx, "object.txt", strings.NewReader("original"), 8, "text/plain"); err != nil {
		t.Fatal(err)
	}
	reader := io.MultiReader(strings.NewReader("partial"), errorReader{})
	if err := storage.Put(ctx, "object.txt", reader, 8, "text/plain"); err == nil {
		t.Fatal("Put() error = nil")
	}
	object, err := storage.Open(ctx, "object.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer object.Close()
	content, err := io.ReadAll(object)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "original" {
		t.Fatalf("content = %q, want original", content)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }
