package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment 附件表
type Attachment struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	UserID           uint           `gorm:"index:idx_user_id" json:"user_id"` // 上传者ID（可为空）
	User             *User          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	FileName         string         `gorm:"size:255;not null" json:"file_name"` // 原始文件名
	FilePath         string         `gorm:"size:500;not null" json:"file_path"` // 兼容字段，私有文件保存 object key
	FileSize         int64          `gorm:"not null" json:"file_size"`          // 文件大小(字节)
	FileType         string         `gorm:"size:50;not null;index" json:"file_type"`
	MimeType         string         `gorm:"size:100;not null" json:"mime_type"`
	Extension        string         `gorm:"size:20;not null" json:"extension"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	StorageType      string         `gorm:"size:20;default:'local'" json:"storage_type"`
	URL              string         `gorm:"size:500" json:"url"`
	UsageCount       int            `gorm:"default:0" json:"usage_count"`
	ObjectKey        string         `gorm:"size:500;index" json:"-"`
	Checksum         string         `gorm:"size:64;index" json:"checksum"`
	PreviewStatus    string         `gorm:"size:20;index;default:'none'" json:"preview_status"`
	PreviewObjectKey string         `gorm:"size:500" json:"-"`
	PreviewMimeType  string         `gorm:"size:100" json:"preview_mime_type"`
	PreviewError     string         `gorm:"size:500" json:"-"`
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
