package uploader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// TencentCOSUploader 腾讯云COS上传实现
type TencentCOSUploader struct {
	client *cos.Client
	domain string // CDN域名或存储桶域名
}

// NewTencentCOSUploader 创建腾讯云COS上传器
func NewTencentCOSUploader() *TencentCOSUploader {
	cfg := config.AppConfig.Upload.TencentCOS
	u, _ := url.Parse(cfg.BucketURL)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &TencentCOSUploader{
		client: client,
		domain: cfg.Domain,
	}
}

// Upload 上传文件到腾讯云COS
func (u *TencentCOSUploader) Upload(input *UploadInput, dstPath string) (string, error) {
	src, err := input.Open()
	if err != nil {
		return "", fmt.Errorf("open source file failed: %w", err)
	}
	defer src.Close()

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // 增加超时时间
	defer cancel()

	// 上传文件
	// cos SDK Put 方法接受 io.Reader
	// 需要注意的是，如果 src 是文件流，SDK 会自动获取大小
	// 如果 input.Size 已知，最好传递给 option (这里暂略，SDK handle file well)
	_, err = u.client.Object.Put(ctx, dstPath, src, nil)
	opt := &cos.ObjectPutOptions{
		// 图片
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "image/png", // todo: 设置Content-Type
		},
	}
	_, err = u.client.Object.Put(context.Background(), dstPath, src, opt)
	if err != nil {
		return "", fmt.Errorf("upload to cos failed: %w", err)
	}

	// 构造返回URL
	var fileURL string
	if u.domain != "" {
		// 简单处理：直接拼接
		// 需确保 domain 结尾无 /，dstPath 开头无 / (或处理之)
		// 这里假设配置正确
		fileURL = fmt.Sprintf("%s/%s", u.domain, dstPath)
	} else {
		// 使用默认的Bucket URL
		fileURL = u.client.Object.GetObjectURL(dstPath).String()
	}

	return fileURL, nil
}
