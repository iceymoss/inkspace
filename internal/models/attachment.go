package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment 附件表
type Attachment struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"index:idx_user_id" json:"user_id"` // 上传者ID（可为空）
	User       *User          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	FileName   string         `gorm:"size:255;not null" json:"file_name"` // 原始文件名
	FilePath   string         `gorm:"size:500;not null" json:"file_path"` // 存储路径
	FileSize   int64          `gorm:"not null" json:"file_size"` // 文件大小(字节)
	FileType   string         `gorm:"size:50;not null;index" json:"file_type"` // 文件类型: image, video, audio, document, other
	MimeType   string         `gorm:"size:100;not null" json:"mime_type"` // MIME类型
	Extension  string         `gorm:"size:20;not null" json:"extension"` // 文件扩展名
	Width      int            `json:"width"` // 图片宽度
	Height     int            `json:"height"` // 图片高度
	StorageType string        `gorm:"size:20;default:'local'" json:"storage_type"` // 存储类型: local, oss, cos, qiniu
	URL        string         `gorm:"size:500" json:"url"` // 访问URL
	UsageCount int            `gorm:"default:0" json:"usage_count"` // 使用次数
}

type AttachmentRequest struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
	MimeType string `json:"mime_type"`
}

type AttachmentResponse struct {
	ID         uint      `json:"id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	MimeType   string    `json:"mime_type"`
	Extension  string    `json:"extension"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	URL        string    `json:"url"`
	UsageCount int       `json:"usage_count"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *Attachment) ToResponse() *AttachmentResponse {
	return &AttachmentResponse{
		ID:         a.ID,
		FileName:   a.FileName,
		FilePath:   a.FilePath,
		FileSize:   a.FileSize,
		FileType:   a.FileType,
		MimeType:   a.MimeType,
		Extension:  a.Extension,
		Width:      a.Width,
		Height:     a.Height,
		URL:        a.URL,
		UsageCount: a.UsageCount,
		CreatedAt:  a.CreatedAt,
	}
}

