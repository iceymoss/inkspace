package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知
type Notification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID     uint   `gorm:"not null;index" json:"user_id"`             // 接收通知的用户
	FromUserID uint   `gorm:"not null;index" json:"from_user_id"`        // 触发通知的用户
	Type       string `gorm:"type:varchar(50);not null" json:"type"`     // comment/like/favorite/follow/reply
	Content    string `gorm:"type:text" json:"content"`                  // 通知内容
	ArticleID  *uint  `gorm:"index" json:"article_id,omitempty"`         // 相关文章ID
	WorkID     *uint  `gorm:"index" json:"work_id,omitempty"`            // 相关作品ID
	CommentID  *uint  `gorm:"index" json:"comment_id,omitempty"`         // 相关评论ID
	IsRead     bool   `gorm:"default:false;index" json:"is_read"`        // 是否已读

	User     *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FromUser *User    `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	Article  *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Work     *Work    `gorm:"foreignKey:WorkID" json:"work,omitempty"`
	Comment  *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID         uint          `json:"id"`
	UserID     uint          `json:"user_id"`
	FromUserID uint          `json:"from_user_id"`
	FromUser   *UserResponse `json:"from_user,omitempty"`
	Type       string        `json:"type"`
	Content    string        `json:"content"`
	ArticleID  *uint         `json:"article_id,omitempty"`
	WorkID     *uint         `json:"work_id,omitempty"`
	CommentID  *uint         `json:"comment_id,omitempty"`
	IsRead     bool          `json:"is_read"`
	CreatedAt  time.Time     `json:"created_at"`
}

// ToResponse 转换为响应格式
func (n *Notification) ToResponse() *NotificationResponse {
	resp := &NotificationResponse{
		ID:         n.ID,
		UserID:     n.UserID,
		FromUserID: n.FromUserID,
		Type:       n.Type,
		Content:    n.Content,
		ArticleID:  n.ArticleID,
		WorkID:     n.WorkID,
		CommentID:  n.CommentID,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt,
	}

	if n.FromUser != nil {
		resp.FromUser = n.FromUser.ToResponse()
	}

	return resp
}
