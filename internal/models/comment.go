package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ArticleID uint           `gorm:"index;not null" json:"article_id"`
	Article   *Article       `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	UserID    uint           `gorm:"index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content" binding:"required"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Email     string         `gorm:"size:100" json:"email"`
	Website   string         `gorm:"size:200" json:"website"`
	Status    int            `gorm:"default:1;index" json:"status"` // 1: approved, 0: pending
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
	ID        uint            `json:"id"`
	ArticleID uint            `json:"article_id"`
	UserID    uint            `json:"user_id"`
	User      *UserResponse   `json:"user,omitempty"`
	Content   string          `json:"content"`
	ParentID  *uint           `json:"parent_id"`
	Nickname  string          `json:"nickname"`
	Email     string          `json:"email"`
	Website   string          `json:"website"`
	Status    int             `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

func (c *Comment) ToResponse() *CommentResponse {
	resp := &CommentResponse{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		Content:   c.Content,
		ParentID:  c.ParentID,
		Nickname:  c.Nickname,
		Email:     c.Email,
		Website:   c.Website,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
	}

	if c.User != nil {
		resp.User = c.User.ToResponse()
	}

	return resp
}

