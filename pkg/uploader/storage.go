package uploader

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
)

// ObjectInfo describes an object stored by a Storage implementation.
type ObjectInfo struct {
	Key         string
	Size        int64
	ContentType string
	ModTime     time.Time
	ETag        string
}

// Storage provides the object operations needed by upload and download services.
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Stat(ctx context.Context, key string) (*ObjectInfo, error)
}

// NewStorage creates the configured storage provider.
func NewStorage() (Storage, error) {
	if config.AppConfig == nil {
		return nil, fmt.Errorf("application config is not initialized")
	}
	return NewStorageWithConfig(config.AppConfig.Upload)
}

// NewStorageWithConfig creates a storage provider from an upload configuration.
func NewStorageWithConfig(cfg config.UploadConfig) (Storage, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.StorageType)) {
	case "", "local":
		return NewLocalStorage(cfg.SavePath)
	case "cos":
		return NewTencentCOSStorage(cfg.TencentCOS)
	default:
		return nil, fmt.Errorf("unsupported upload storage type %q", cfg.StorageType)
	}
}
