package handler

import (
	"fmt"
	"io"
	"mysite/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

