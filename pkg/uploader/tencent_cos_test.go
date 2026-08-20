package uploader

import (
	"context"
	"hash/crc64"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/iceymoss/inkspace/internal/config"
)

func TestTencentCOSStoragePutOnceWithContentType(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/images/photo.jpg" {
			t.Errorf("path = %q, want /images/photo.jpg", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "image/jpeg" {
			t.Errorf("Content-Type = %q, want image/jpeg", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ReadAll() error = %v", err)
		}
		if string(body) != "image-data" {
			t.Errorf("body = %q, want image-data", body)
		}
		checksum := crc64.Checksum(body, crc64.MakeTable(crc64.ECMA))
		w.Header().Set("x-cos-hash-crc64ecma", strconv.FormatUint(checksum, 10))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage, err := NewTencentCOSStorage(config.TencentCOSConfig{BucketURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), "images/photo.jpg", strings.NewReader("image-data"), 10, "image/jpeg"); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("request count = %d, want 1", got)
	}
}

func TestTencentCOSStoragePutUsesContext(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage, err := NewTencentCOSStorage(config.TencentCOSConfig{BucketURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := storage.Put(ctx, "object.txt", strings.NewReader("data"), 4, "text/plain"); err == nil {
		t.Fatal("Put() error = nil for canceled context")
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("request count = %d, want 0", got)
	}
}
