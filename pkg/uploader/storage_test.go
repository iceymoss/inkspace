package uploader

import (
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/config"
)

func TestNewStorageWithConfig(t *testing.T) {
	storage, err := NewStorageWithConfig(config.UploadConfig{StorageType: "local", SavePath: t.TempDir()})
	if err != nil {
		t.Fatalf("NewStorageWithConfig(local) error = %v", err)
	}
	if _, ok := storage.(*LocalStorage); !ok {
		t.Fatalf("NewStorageWithConfig(local) returned %T", storage)
	}
	storage, err = NewStorageWithConfig(config.UploadConfig{
		StorageType: "cos",
		TencentCOS:  config.TencentCOSConfig{BucketURL: "https://bucket.example.com"},
	})
	if err != nil {
		t.Fatalf("NewStorageWithConfig(cos) error = %v", err)
	}
	if _, ok := storage.(*TencentCOSStorage); !ok {
		t.Fatalf("NewStorageWithConfig(cos) returned %T", storage)
	}

	if _, err := NewStorageWithConfig(config.UploadConfig{StorageType: "unknown"}); err == nil {
		t.Fatal("NewStorageWithConfig(unknown) error = nil")
	}
}

func TestUploadProviderDoesNotFallbackOnInvalidStorage(t *testing.T) {
	previous := config.AppConfig
	config.AppConfig = &config.Config{Upload: config.UploadConfig{StorageType: "invalid"}}
	t.Cleanup(func() { config.AppConfig = previous })

	provider := (&UploadProvider{}).NewUploadProvider()
	if _, err := provider.Upload(&UploadInput{Reader: strings.NewReader("data"), Size: 4}, "object.txt"); err == nil {
		t.Fatal("Upload() error = nil")
	}
}
