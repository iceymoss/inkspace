package uploader

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/iceymoss/inkspace/internal/config"
)

// LocalStorage stores objects below a fixed local root directory.
type LocalStorage struct {
	root string
}

// NewLocalStorage creates local storage rooted at root.
func NewLocalStorage(root string) (*LocalStorage, error) {
	if root == "" {
		root = "uploads"
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve storage root: %w", err)
	}
	if err := os.MkdirAll(absRoot, 0755); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}
	resolvedRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve storage root symlinks: %w", err)
	}
	return &LocalStorage{root: resolvedRoot}, nil
}

func (s *LocalStorage) Put(ctx context.Context, key string, r io.Reader, _ int64, _ string) error {
	name, err := s.objectName(key)
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(s.root)
	if err != nil {
		return fmt.Errorf("open storage root: %w", err)
	}
	defer root.Close()
	if err := root.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return fmt.Errorf("create object directory: %w", err)
	}
	tempName := filepath.Join(filepath.Dir(name), ".inkspace-upload-"+uuid.NewString())
	temp, err := root.OpenFile(tempName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create temporary object: %w", err)
	}
	keepTemp := true
	defer func() {
		if keepTemp {
			_ = root.Remove(tempName)
		}
	}()
	_, copyErr := io.Copy(temp, &contextReader{ctx: ctx, reader: r})
	closeErr := temp.Close()
	if copyErr != nil {
		return fmt.Errorf("write object: %w", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close object: %w", closeErr)
	}
	if err := root.Rename(tempName, name); err != nil {
		return fmt.Errorf("replace object: %w", err)
	}
	keepTemp = false
	return nil
}

func (s *LocalStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	name, err := s.objectName(key)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(s.root)
	if err != nil {
		return nil, fmt.Errorf("open storage root: %w", err)
	}
	file, err := root.Open(name)
	root.Close()
	if err != nil {
		return nil, fmt.Errorf("open object: %w", err)
	}
	return file, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	name, err := s.objectName(key)
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(s.root)
	if err != nil {
		return fmt.Errorf("open storage root: %w", err)
	}
	defer root.Close()
	if err := root.Remove(name); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func (s *LocalStorage) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	name, err := s.objectName(key)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(s.root)
	if err != nil {
		return nil, fmt.Errorf("open storage root: %w", err)
	}
	defer root.Close()
	info, err := root.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("stat object: %w", err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("object %q is not a regular file", key)
	}
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		file, openErr := root.Open(name)
		if openErr != nil {
			return nil, fmt.Errorf("open object for content type: %w", openErr)
		}
		buffer := make([]byte, 512)
		n, readErr := file.Read(buffer)
		file.Close()
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("read object content type: %w", readErr)
		}
		contentType = http.DetectContentType(buffer[:n])
	}
	return &ObjectInfo{Key: key, Size: info.Size(), ContentType: contentType, ModTime: info.ModTime()}, nil
}

func (s *LocalStorage) objectName(key string) (string, error) {
	if key == "" || filepath.IsAbs(key) {
		return "", fmt.Errorf("invalid object key %q", key)
	}
	name := filepath.Clean(key)
	if !filepath.IsLocal(name) {
		return "", fmt.Errorf("object key %q escapes storage root", key)
	}
	return name, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

// LocalUploader preserves the existing URL-returning upload API.
type LocalUploader struct {
	BaseURL  string
	SavePath string
	storage  *LocalStorage
	initErr  error
}

func NewLocalUploader() *LocalUploader {
	savePath := "uploads"
	if config.AppConfig != nil && config.AppConfig.Upload.SavePath != "" {
		savePath = config.AppConfig.Upload.SavePath
	}
	storage, err := NewLocalStorage(savePath)
	return &LocalUploader{BaseURL: "/" + savePath, SavePath: savePath, storage: storage, initErr: err}
}

func (u *LocalUploader) Upload(input *UploadInput, dstPath string) (string, error) {
	return u.UploadContext(context.Background(), input, dstPath)
}

// UploadContext uploads a local object using the caller's context.
func (u *LocalUploader) UploadContext(ctx context.Context, input *UploadInput, dstPath string) (string, error) {
	if u.initErr != nil {
		return "", u.initErr
	}
	src, err := input.Open()
	if err != nil {
		return "", fmt.Errorf("open source file: %w", err)
	}
	defer src.Close()
	if err := u.storage.Put(ctx, dstPath, src, input.Size, input.contentType()); err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Join(u.BaseURL, dstPath)), nil
}
