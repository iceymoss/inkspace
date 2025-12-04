package models

import (
	"time"

	"gorm.io/gorm"
)

// ArticleLike 文章点赞记录表
type ArticleLike struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ArticleID uint           `gorm:"index:idx_article_user,priority:1;not null" json:"article_id"` // 文章ID
	Article   *Article       `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"`
	UserID    uint           `gorm:"index:idx_article_user,priority:2;index:idx_user_id" json:"user_id"` // 用户ID（0表示游客）
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	IP        string         `gorm:"size:50;index" json:"ip"` // IP地址（游客点赞）
}

// CommentLike 评论点赞记录表
type CommentLike struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CommentID uint           `gorm:"index:idx_comment_user,priority:1;not null" json:"comment_id"` // 评论ID
	Comment   *Comment       `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"comment,omitempty"`
	UserID    uint           `gorm:"index:idx_comment_user,priority:2;index:idx_user_id" json:"user_id"` // 用户ID（0表示游客）
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	IP        string         `gorm:"size:50;index" json:"ip"` // IP地址（游客点赞）
}

// 组合唯一索引
func (ArticleLike) TableName() string {
	return "article_likes"
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

// 检查是否已点赞的请求
type CheckLikeRequest struct {
	TargetID uint   `json:"target_id" binding:"required"`
	UserID   uint   `json:"user_id"`
	IP       string `json:"ip"`
}

