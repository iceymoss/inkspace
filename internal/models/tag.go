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
	Name         string         `gorm:"uniqueIndex;size:50;not null" json:"name" binding:"required"`
	Slug         string         `gorm:"uniqueIndex;size:50" json:"slug"`
	Color        string         `gorm:"size:20;default:'#409eff'" json:"color"`
	ArticleCount int            `gorm:"default:0;not null" json:"article_count"`
}

type TagRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Slug  string `json:"slug" binding:"max=50"`
	Color string `json:"color" binding:"max=20"`
}

type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		Color:     t.Color,
		CreatedAt: t.CreatedAt,
	}
}

