package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name" binding:"required"`
	Slug        string         `gorm:"uniqueIndex;size:50" json:"slug"`
	Description string         `gorm:"size:200" json:"description"`
	Sort        int            `gorm:"default:0" json:"sort"`
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"max=50"`
	Description string `json:"description" binding:"max=200"`
	Sort        int    `json:"sort"`
}

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Sort:        c.Sort,
		CreatedAt:   c.CreatedAt,
	}
}

