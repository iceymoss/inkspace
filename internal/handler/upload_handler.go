package handler

import (
	"fmt"
	"io"
	"mysite/internal/database"
	"mysite/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-redis/redis/v8"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传图片
// POST /api/upload/image
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 验证文件大小（限制5MB）
	if file.Size > 5*1024*1024 {
		utils.BadRequest(c, "文件大小不能超过5MB")
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		utils.BadRequest(c, "只支持上传 jpg, jpeg, png, gif, webp 格式的图片")
		return
	}

	// 创建上传目录
	uploadDir := "uploads/images"
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, dateDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		utils.InternalServerError(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(fullDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		utils.InternalServerError(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		utils.InternalServerError(c, "保存文件失败")
		return
	}

	// 返回文件访问URL
	fileURL := fmt.Sprintf("/%s/%s/%s", uploadDir, dateDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// UploadAvatar 上传头像（限制更严格）
// POST /api/upload/avatar
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 验证文件大小（限制2MB）
	if file.Size > 2*1024*1024 {
		utils.BadRequest(c, "头像大小不能超过2MB")
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		utils.BadRequest(c, "只支持上传 jpg, jpeg, png, webp 格式的头像")
		return
	}

	// 创建上传目录
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.InternalServerError(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		utils.InternalServerError(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		utils.InternalServerError(c, "保存文件失败")
		return
	}

	// 返回文件访问URL
	fileURL := fmt.Sprintf("/%s/%s", uploadDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// UploadPhoto 上传摄影作品原图（不压缩，保留原图质量）
// POST /api/upload/photo
func (h *UploadHandler) UploadPhoto(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 验证文件大小（摄影作品限制20MB，保留高质量）
	if file.Size > 20*1024*1024 {
		utils.BadRequest(c, "摄影作品文件大小不能超过20MB")
		return
	}

	// 验证文件类型（只支持高质量图片格式）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExts[ext] {
		utils.BadRequest(c, "摄影作品只支持 jpg, jpeg, png 格式")
		return
	}

	// 创建专门的摄影作品上传目录
	uploadDir := "uploads/photos"
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, dateDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		utils.InternalServerError(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(fullDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		utils.InternalServerError(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 直接复制文件内容（保留原图质量，不进行任何压缩）
	if _, err := io.Copy(dst, src); err != nil {
		utils.InternalServerError(c, "保存文件失败")
		return
	}

	// 返回文件访问URL
	fileURL := fmt.Sprintf("/%s/%s/%s", uploadDir, dateDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": file.Filename,
		"size":     file.Size,
		"type":     "photo", // 标记为摄影作品原图
	})
}

// UploadMarkdownImage 上传Markdown编辑器图片（公开API，但限制每天100张）
// POST /api/upload/markdown-image
func (h *UploadHandler) UploadMarkdownImage(c *gin.Context) {
	// 获取客户端IP作为标识
	clientIP := c.ClientIP()
	
	// 获取当前日期作为Redis key的一部分
	today := time.Now().Format("2006-01-02")
	limitKey := fmt.Sprintf("upload:markdown-image:%s:%s", today, clientIP)
	
	// 检查今天是否已上传超过100张
	ctx := c.Request.Context()
	count, err := database.RDB.Get(ctx, limitKey).Int()
	if err != nil && err != redis.Nil {
		utils.InternalServerError(c, "检查上传限制失败")
		return
	}
	
	if count >= 100 {
		utils.Error(c, 429, "今天上传图片数量已达上限（100张），请明天再试")
		return
	}
	
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 验证文件大小（限制5MB）
	if file.Size > 5*1024*1024 {
		utils.BadRequest(c, "文件大小不能超过5MB")
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		utils.BadRequest(c, "只支持上传 jpg, jpeg, png, gif, webp 格式的图片")
		return
	}

	// 创建上传目录
	uploadDir := "uploads/images"
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, dateDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		utils.InternalServerError(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(fullDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		utils.InternalServerError(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		utils.InternalServerError(c, "保存文件失败")
		return
	}

	// 增加上传计数（设置24小时过期）
	expiration := 24 * time.Hour
	if err := database.RDB.Incr(ctx, limitKey).Err(); err != nil {
		// 如果计数失败，记录日志但不影响上传
		fmt.Printf("Failed to increment upload count: %v\n", err)
	} else {
		// 设置过期时间（如果key不存在则设置，存在则不更新）
		database.RDB.Expire(ctx, limitKey, expiration)
	}

	// 返回文件访问URL
	fileURL := fmt.Sprintf("/%s/%s/%s", uploadDir, dateDir, filename)

	utils.Success(c, gin.H{
		"url":      fileURL,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
