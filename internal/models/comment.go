package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ArticleID  uint           `gorm:"index:idx_article_id;not null" json:"article_id"`
	Article    *Article       `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"`
	UserID     uint           `gorm:"index:idx_user_id" json:"user_id"`
	User       *User          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Content    string         `gorm:"type:text;not null" json:"content" binding:"required"`
	ParentID   *uint          `gorm:"index:idx_parent_id" json:"parent_id"`
	Parent     *Comment       `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	RootID     *uint          `gorm:"index:idx_root_id" json:"root_id"` // 根评论ID，方便查询
	ReplyToID  *uint          `json:"reply_to_id"` // 回复的评论ID
	Nickname   string         `gorm:"size:50" json:"nickname"`
	Email      string         `gorm:"size:100" json:"email"`
	Website    string         `gorm:"size:200" json:"website"`
	IP         string         `gorm:"size:50" json:"ip"`
	UserAgent  string         `gorm:"size:255" json:"user_agent"`
	Status     int            `gorm:"default:1;index:idx_status_created" json:"status"` // 1: approved, 0: pending, -1: rejected
	LikeCount  int            `gorm:"default:0;not null" json:"like_count"`
	ReplyCount int            `gorm:"default:0;not null" json:"reply_count"`
}

type CommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,max=500"`
	ParentID  *uint  `json:"parent_id"`
	Nickname  string `json:"nickname" binding:"max=50"`
	Email     string `json:"email" binding:"email,max=100"`
	Website   string `json:"website" binding:"max=200"`
}

type CommentListQuery struct {
	Page      int  `form:"page,default=1"`
	PageSize  int  `form:"page_size,default=10"`
	ArticleID uint `form:"article_id"`
	Status    *int `form:"status"`
}

type CommentResponse struct {
	ID         uint              `json:"id"`
	ArticleID  uint              `json:"article_id"`
	UserID     uint              `json:"user_id"`
	User       *UserResponse     `json:"user,omitempty"`
	Content    string            `json:"content"`
	ParentID   *uint             `json:"parent_id"`
	RootID     *uint             `json:"root_id"`
	Nickname   string            `json:"nickname"`
	Email      string            `json:"email"`
	Website    string            `json:"website"`
	Status     int               `json:"status"`
	LikeCount  int               `json:"like_count"`
	ReplyCount int               `json:"reply_count"`
	CreatedAt  time.Time         `json:"created_at"`
	Replies    []CommentResponse `json:"replies,omitempty"`
}

func (c *Comment) ToResponse() *CommentResponse {
	resp := &CommentResponse{
		ID:         c.ID,
		ArticleID:  c.ArticleID,
		UserID:     c.UserID,
		Content:    c.Content,
		ParentID:   c.ParentID,
		RootID:     c.RootID,
		Nickname:   c.Nickname,
		Email:      c.Email,
		Website:    c.Website,
		Status:     c.Status,
		LikeCount:  c.LikeCount,
		ReplyCount: c.ReplyCount,
		CreatedAt:  c.CreatedAt,
	}

	if c.User != nil {
		resp.User = c.User.ToResponse()
	}

	return resp
}

