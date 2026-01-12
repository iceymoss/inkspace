package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/utils"
	"github.com/iceymoss/inkspace/pkg/uploader"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploader uploader.Uploader
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploader: uploader.NewLocalUploader(),
	}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	h.handleUpload(c, "images", 5*1024*1024, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}, false)
}

// UploadAvatar 上传头像
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	h.handleUpload(c, "avatars", 2*1024*1024, []string{".jpg", ".jpeg", ".png", ".webp"}, false)
}

// UploadPhoto 上传摄影作品
func (h *UploadHandler) UploadPhoto(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	limitKey := fmt.Sprintf("upload:work-image:%s:%s", today, userID)

	ctx := c.Request.Context()
	count, err := database.RDB.Get(ctx, limitKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		utils.InternalServerError(c, "检查上传限制失败")
		return
	}

	if count >= 50 {
		utils.Error(c, 429, "今天上传图片数量已达上限（50张），请明天再试")
		return
	}

	// 摄影作品限制20MB，超过则压缩
	h.handleUpload(c, "photos", 20*1024*1024, []string{".jpg", ".jpeg", ".png"}, true)
}

// UploadMarkdownImage 上传Markdown插图
func (h *UploadHandler) UploadMarkdownImage(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	limitKey := fmt.Sprintf("upload:markdown-image:%s:%s", today, userID)

	ctx := c.Request.Context()
	count, err := database.RDB.Get(ctx, limitKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		utils.InternalServerError(c, "检查上传限制失败")
		return
	}

	if count >= 100 {
		utils.Error(c, 429, "今天上传图片数量已达上限（100张），请明天再试")
		return
	}

	success := h.handleUpload(c, "images", 5*1024*1024, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}, false)

	if success {
		expiration := 24 * time.Hour
		if err := database.RDB.Incr(ctx, limitKey).Err(); err != nil {
			fmt.Printf("Failed to increment upload count: %v\n", err)
		} else {
			database.RDB.Expire(ctx, limitKey, expiration)
		}
	}
}

// handleUpload 通用上传逻辑
func (h *UploadHandler) handleUpload(c *gin.Context, subDir string, maxSize int64, allowedExts []string, compressLarge bool) bool {
	// 1. 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return false
	}

	// 2. 验证类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	isAllowed := false
	for _, allowed := range allowedExts {
		if allowed == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		utils.BadRequest(c, fmt.Sprintf("不支持的文件格式: %s", ext))
		return false
	}

	contentType, err := uploader.GetFileContentType(file)
	if err != nil {
		utils.InternalServerError(c, "获取文件类型失败")
		return false
	}

	// 3. 准备 UploadInput
	var input *uploader.UploadInput
	var tempPathToDelete string

	// 压缩阈值 20MB
	const CompressThreshold = 20 * 1024 * 1024

	// 4. 大文件处理 (压缩)
	if compressLarge && file.Size > CompressThreshold {
		// 压缩
		compressedPath, compressed, err := uploader.CompressImage(file, CompressThreshold, 85)
		if err != nil {
			utils.InternalServerError(c, "图片处理失败")
			return false
		}

		if compressed {
			tempPathToDelete = compressedPath
			input, err = uploader.NewUploadInputFromLocalPath(compressedPath, file.Filename)
			if err != nil {
				utils.InternalServerError(c, "读取处理后的图片失败")
				if tempPathToDelete != "" {
					os.Remove(tempPathToDelete)
				}
				return false
			}
		} else {
			// 未压缩（可能是格式不支持压缩），使用原文件
			input = uploader.NewUploadInputFromFileHeader(file)
		}
	} else {
		// 5. 普通文件大小检查
		if file.Size > maxSize {
			utils.BadRequest(c, fmt.Sprintf("文件大小不能超过 %.2f MB", float64(maxSize)/1024/1024))
			return false
		}
		input = uploader.NewUploadInputFromFileHeader(file)
	}

	defer func() {
		if tempPathToDelete != "" {
			os.Remove(tempPathToDelete)
		}
	}()

	// 6. 生成目标路径
	// 格式: subDir/YYYY/MM/DD/uuid.ext
	dateDir := time.Now().Format("2006/01/02")
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := fmt.Sprintf("%s/%s/%s", subDir, dateDir, filename)

	// 针对 avatars 目录结构可能不同 (原代码是 uploads/avatars/uuid.ext)
	if subDir == "avatars" {
		dstPath = fmt.Sprintf("%s/%s", subDir, filename)
	}

	// 7. 执行上传
	url, err := h.uploader.Upload(input, dstPath)
	if err != nil {
		utils.InternalServerError(c, "文件上传失败")
		fmt.Printf("Upload failed: %v\n", err)
		return false
	}

	// 8. 返回结果
	utils.Success(c, gin.H{
		"url":          url,
		"filename":     file.Filename,
		"size":         input.Size,
		"type":         subDir, // 可选
		"content_type": contentType,
	})

	return true
}
