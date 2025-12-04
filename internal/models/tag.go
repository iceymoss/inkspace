package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"uniqueIndex:idx_user_tag_name;size:50;not null" json:"name" binding:"required"`
	Slug         string         `gorm:"size:50" json:"slug"`
	Color        string         `gorm:"size:20;default:'#409eff'" json:"color"`
	ArticleCount int            `gorm:"default:0;not null" json:"article_count"`
	UserID       *uint          `gorm:"uniqueIndex:idx_user_tag_name;index:idx_user_id" json:"user_id"` // 标签归属用户（NULL表示公共标签）
}

type TagRequest struct {
	Name   string `json:"name" binding:"required,max=50"`
	Slug   string `json:"slug" binding:"max=50"`
	Color  string `json:"color" binding:"max=20"`
	UserID *uint  `json:"user_id"` // 可选，NULL表示公共标签
}

type TagResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Color        string    `json:"color"`
	UserID       *uint     `json:"user_id"`
	ArticleCount int       `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
}

func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		Color:        t.Color,
		UserID:       t.UserID,
		ArticleCount: t.ArticleCount,
		CreatedAt:    t.CreatedAt,
	}
}

