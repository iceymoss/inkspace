package uploader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/iceymoss/inkspace/internal/config"
)

// LocalUploader 本地文件上传实现
type LocalUploader struct {
	BaseURL  string // 本地文件访问的基础URL
	SavePath string // 本地文件保存的根目录
}

// NewLocalUploader 创建本地上传器
func NewLocalUploader() *LocalUploader {
	savePath := config.AppConfig.Upload.SavePath
	if savePath == "" {
		savePath = "uploads"
	}
	return &LocalUploader{
		BaseURL:  "/" + savePath,
		SavePath: savePath,
	}
}

// Upload 上传文件到本地
func (u *LocalUploader) Upload(input *UploadInput, dstPath string) (string, error) {
	// 拼接完整保存路径
	fullPath := filepath.Join(u.SavePath, dstPath)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directory failed: %w", err)
	}

	// 打开源文件
	src, err := input.Open()
	if err != nil {
		return "", fmt.Errorf("open source file failed: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create destination file failed: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copy file content failed: %w", err)
	}

	// 返回访问URL
	// 统一使用正斜杠
	// windows下 filepath.Join 会使用反斜杠，需要替换
	urlPath := filepath.ToSlash(filepath.Join(u.BaseURL, dstPath))
	return urlPath, nil
}
