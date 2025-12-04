package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"uniqueIndex;size:50;not null" json:"name" binding:"required"`
	Slug         string         `gorm:"uniqueIndex;size:50" json:"slug"`
	Description  string         `gorm:"size:200" json:"description"`
	Logo         string         `gorm:"size:255;not null" json:"logo" binding:"required"` // 分类Logo（必填）
	Cover        string         `gorm:"size:255" json:"cover"`
	Sort         int            `gorm:"default:0;index:idx_sort" json:"sort"`
	ArticleCount int            `gorm:"default:0;not null" json:"article_count"`
	ParentID     *uint          `gorm:"index:idx_parent_id" json:"parent_id"` // 支持二级分类
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"max=50"`
	Description string `json:"description" binding:"max=200"`
	Logo        string `json:"logo" binding:"required"` // Logo必填
	Sort        int    `json:"sort"`
}

type CategoryResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Logo         string    `json:"logo"`
	ArticleCount int       `json:"article_count"`
	Sort         int       `json:"sort"`
	CreatedAt    time.Time `json:"created_at"`
}

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		Logo:         c.Logo,
		ArticleCount: c.ArticleCount,
		Sort:         c.Sort,
		CreatedAt:    c.CreatedAt,
	}
}

