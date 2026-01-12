package uploader

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// CompressImage 压缩图片
// file: 源文件
// maxSize: 触发压缩的大小阈值 (字节)
// quality: 压缩质量 (1-100, 仅对JPEG有效)
// 返回: 压缩后的临时文件路径 (调用者需负责删除), 是否进行了压缩, 错误
func CompressImage(file *multipart.FileHeader, maxSize int64, quality int) (string, bool, error) {
	// 如果文件大小未超过阈值，直接返回
	if file.Size <= maxSize {
		return "", false, nil
	}

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return "", false, err
	}
	defer src.Close()

	// 解码图片
	img, format, err := image.Decode(src)
	if err != nil {
		// 如果无法解码（非图片或不支持的格式），则不压缩
		return "", false, nil
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "compressed-*"+filepath.Ext(file.Filename))
	if err != nil {
		return "", false, err
	}
	tmpPath := tmpFile.Name()

	// 压缩并保存
	// 注意：这里简单的使用 resize 库进行尺寸调整作为压缩示例，或者直接降低 JPEG 质量
	// 为了保持尺寸不变只降低质量/大小，我们重新编码

	// 如果图片非常大，可以考虑缩小尺寸
	// 这里简单策略：如果宽度大于 1920，则缩放到 1920
	bounds := img.Bounds()
	if bounds.Dx() > 1920 {
		img = resize.Resize(1920, 0, img, resize.Lanczos3)
	}

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: quality})
	case "png":
		// PNG 是无损的，通常无法通过降低质量参数压缩，除非改变调色板
		// 这里简单使用默认编码
		err = png.Encode(tmpFile, img)
	default:
		// 其他格式原样保存或暂不支持压缩
		// 重置文件指针并复制
		src.Seek(0, 0)
		_, err = io.Copy(tmpFile, src)
	}

	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", false, err
	}

	tmpFile.Close()
	return tmpPath, true, nil
}

// GetFileContentType 获取文件的 Content-Type
// 如果 fileHeader 提供了 Header["Content-Type"]，则直接使用
// 否则尝试从文件名推断或读取文件头（如果是 multipart.File）进行嗅探
func GetFileContentType(file *multipart.FileHeader) (string, error) {
	// 1. 优先尝试从 Header 获取
	if ct := file.Header.Get("Content-Type"); ct != "" {
		return ct, nil
	}

	// 2. 如果 Header 中没有，尝试打开文件读取前512字节进行嗅探
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	// DetectContentType 需要至少 512 字节的数据，如果不够也没关系
	contentType := http.DetectContentType(buffer[:n])

	// http.DetectContentType 可能会返回 generic 的 application/octet-stream
	// 如果是这种情况，或者我们需要更精确的类型，可以结合文件扩展名作为兜底（虽然不安全，但在无法探测时的补充）

	return contentType, nil
}

// IsImage 检查 Content-Type 是否为图片
func IsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}
