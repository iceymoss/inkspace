package uploader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// TencentCOSStorage stores objects in Tencent Cloud COS.
type TencentCOSStorage struct {
	client *cos.Client
}

// NewTencentCOSStorage creates a COS storage provider.
func NewTencentCOSStorage(cfg config.TencentCOSConfig) (*TencentCOSStorage, error) {
	bucketURL, err := url.Parse(cfg.BucketURL)
	if err != nil {
		return nil, fmt.Errorf("parse COS bucket URL: %w", err)
	}
	if bucketURL.Scheme == "" || bucketURL.Host == "" {
		return nil, fmt.Errorf("invalid COS bucket URL %q", cfg.BucketURL)
	}
	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{SecretID: cfg.SecretID, SecretKey: cfg.SecretKey},
	})
	return &TencentCOSStorage{client: client}, nil
}

func (s *TencentCOSStorage) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	options := &cos.ObjectPutOptions{ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
		ContentLength: size,
		ContentType:   contentType,
	}}
	if _, err := s.client.Object.Put(ctx, key, r, options); err != nil {
		return fmt.Errorf("put COS object: %w", err)
	}
	return nil
}

func (s *TencentCOSStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	response, err := s.client.Object.Get(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("open COS object: %w", err)
	}
	return response.Body, nil
}

func (s *TencentCOSStorage) Delete(ctx context.Context, key string) error {
	response, err := s.client.Object.Delete(ctx, key)
	if response != nil && response.Body != nil {
		response.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("delete COS object: %w", err)
	}
	return nil
}

func (s *TencentCOSStorage) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	response, err := s.client.Object.Head(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("stat COS object: %w", err)
	}
	defer response.Body.Close()
	size, err := strconv.ParseInt(response.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse COS object size: %w", err)
	}
	var modTime time.Time
	if value := response.Header.Get("Last-Modified"); value != "" {
		modTime, err = http.ParseTime(value)
		if err != nil {
			return nil, fmt.Errorf("parse COS object modification time: %w", err)
		}
	}
	return &ObjectInfo{
		Key:         key,
		Size:        size,
		ContentType: response.Header.Get("Content-Type"),
		ModTime:     modTime,
		ETag:        strings.Trim(response.Header.Get("ETag"), `"`),
	}, nil
}

// TencentCOSUploader preserves the existing URL-returning upload API.
type TencentCOSUploader struct {
	storage *TencentCOSStorage
	domain  string
	initErr error
}

func NewTencentCOSUploader() *TencentCOSUploader {
	var cfg config.TencentCOSConfig
	if config.AppConfig != nil {
		cfg = config.AppConfig.Upload.TencentCOS
	}
	storage, err := NewTencentCOSStorage(cfg)
	return &TencentCOSUploader{storage: storage, domain: strings.TrimRight(cfg.Domain, "/"), initErr: err}
}

func (u *TencentCOSUploader) Upload(input *UploadInput, dstPath string) (string, error) {
	return u.UploadContext(context.Background(), input, dstPath)
}

// UploadContext uploads a COS object using the caller's context.
func (u *TencentCOSUploader) UploadContext(ctx context.Context, input *UploadInput, dstPath string) (string, error) {
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
	if u.domain != "" {
		return u.domain + "/" + strings.TrimLeft(dstPath, "/"), nil
	}
	return u.storage.client.Object.GetObjectURL(dstPath).String(), nil
}
