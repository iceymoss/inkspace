package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知表
type Notification struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"index:idx_user_status;not null" json:"user_id"` // 接收者ID
	User       *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	FromUserID uint           `gorm:"index:idx_from_user_id" json:"from_user_id"` // 发送者ID
	FromUser   *User          `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE" json:"from_user,omitempty"`
	Type       string         `gorm:"size:50;not null;index" json:"type"` // 类型: comment, reply, like, system, mention
	Title      string         `gorm:"size:200;not null" json:"title"` // 标题
	Content    string         `gorm:"type:text" json:"content"` // 内容
	TargetType string         `gorm:"size:50" json:"target_type"` // 目标类型: article, comment, work
	TargetID   uint           `gorm:"index" json:"target_id"` // 目标ID
	Link       string         `gorm:"size:255" json:"link"` // 跳转链接
	IsRead     bool           `gorm:"default:false;index:idx_user_status" json:"is_read"` // 是否已读
	ReadAt     *time.Time     `gorm:"type:datetime(3)" json:"read_at"` // 阅读时间
}

type NotificationRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	FromUserID uint   `json:"from_user_id"`
	Type       string `json:"type" binding:"required,oneof=comment reply like system mention"`
	Title      string `json:"title" binding:"required,max=200"`
	Content    string `json:"content"`
	TargetType string `json:"target_type"`
	TargetID   uint   `json:"target_id"`
	Link       string `json:"link"`
}

type NotificationListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Type     string `form:"type"`
	IsRead   *bool  `form:"is_read"`
}

type NotificationResponse struct {
	ID         uint           `json:"id"`
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	TargetType string         `json:"target_type"`
	TargetID   uint           `json:"target_id"`
	Link       string         `json:"link"`
	IsRead     bool           `json:"is_read"`
	FromUser   *UserResponse  `json:"from_user,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	ReadAt     *time.Time     `json:"read_at"`
}

func (n *Notification) ToResponse() *NotificationResponse {
	resp := &NotificationResponse{
		ID:         n.ID,
		Type:       n.Type,
		Title:      n.Title,
		Content:    n.Content,
		TargetType: n.TargetType,
		TargetID:   n.TargetID,
		Link:       n.Link,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt,
		ReadAt:     n.ReadAt,
	}

	if n.FromUser != nil {
		resp.FromUser = n.FromUser.ToResponse()
	}

	return resp
}

