package uploader

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/iceymoss/inkspace/internal/config"
)

// UploadInput 上传输入源
type UploadInput struct {
	FileHeader  *multipart.FileHeader // 来自HTTP请求的文件
	LocalPath   string                // 本地文件路径 (用于压缩后的文件)
	Reader      io.Reader             // 通用Reader
	Size        int64                 // 文件大小
	Name        string                // 文件名
	ContentType string                // 文件 MIME 类型
}

// Uploader 文件上传接口
type Uploader interface {
	// Upload 上传文件
	// input: 上传输入源
	// dstPath: 目标路径 (相对于存储根目录)
	// 返回: 访问URL, 错误
	Upload(input *UploadInput, dstPath string) (string, error)
}

// UploadProvider 创建Uploader的工厂
type UploadProvider struct{}

type unavailableUploader struct{ err error }

func (u *unavailableUploader) Upload(_ *UploadInput, _ string) (string, error) {
	return "", u.err
}

// NewUploadProvider 根据配置创建Uploader
func (f *UploadProvider) NewUploadProvider() Uploader {
	if config.AppConfig == nil {
		return NewLocalUploader()
	}
	storage, err := NewStorage()
	if err != nil {
		return &unavailableUploader{err: fmt.Errorf("initialize upload storage: %w", err)}
	}
	switch storage := storage.(type) {
	case *LocalStorage:
		savePath := config.AppConfig.Upload.SavePath
		if savePath == "" {
			savePath = "uploads"
		}
		return &LocalUploader{BaseURL: "/" + savePath, SavePath: savePath, storage: storage}
	case *TencentCOSStorage:
		return &TencentCOSUploader{storage: storage, domain: strings.TrimRight(config.AppConfig.Upload.TencentCOS.Domain, "/")}
	default:
		return NewLocalUploader()
	}
}

func (i *UploadInput) contentType() string {
	if i.ContentType != "" {
		return i.ContentType
	}
	if i.FileHeader != nil {
		if contentType := i.FileHeader.Header.Get("Content-Type"); contentType != "" {
			return contentType
		}
	}
	if contentType := mime.TypeByExtension(filepath.Ext(i.Name)); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

// NewUploadInputFromFileHeader 从 FileHeader 创建输入
func NewUploadInputFromFileHeader(fh *multipart.FileHeader) *UploadInput {
	return &UploadInput{
		FileHeader: fh,
		Size:       fh.Size,
		Name:       fh.Filename,
	}
}

// NewUploadInputFromLocalPath 从本地路径创建输入
func NewUploadInputFromLocalPath(path string, originalName string) (*UploadInput, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &UploadInput{
		LocalPath: path,
		Size:      info.Size(),
		Name:      originalName,
	}, nil
}

// Open 打开输入流
func (i *UploadInput) Open() (io.ReadCloser, error) {
	if i.FileHeader != nil {
		return i.FileHeader.Open()
	}
	if i.LocalPath != "" {
		return os.Open(i.LocalPath)
	}
	if i.Reader != nil {
		// Reader 无法 Close，只能返回 NopCloser
		return io.NopCloser(i.Reader), nil
	}
	return nil, os.ErrInvalid
}
