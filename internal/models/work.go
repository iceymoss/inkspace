package models

import (
	"time"

	"gorm.io/gorm"
)

type Work struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:200;not null" json:"title" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Cover       string         `gorm:"size:255" json:"cover"`
	Images      string         `gorm:"type:text" json:"images"` // JSON array of image URLs
	Link        string         `gorm:"size:255" json:"link"`
	GithubURL   string         `gorm:"size:255" json:"github_url"`
	DemoURL     string         `gorm:"size:255" json:"demo_url"`
	TechStack   string         `gorm:"size:500" json:"tech_stack"` // 技术栈，逗号分隔
	Sort        int            `gorm:"default:0;index:idx_status_sort" json:"sort"`
	ViewCount   int            `gorm:"default:0;not null" json:"view_count"`
	Status      int            `gorm:"default:1;index:idx_status_sort" json:"status"` // 1: published, 0: draft
}

type WorkRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Description string   `json:"description"`
	Cover       string   `json:"cover"`
	Images      []string `json:"images"`
	Link        string   `json:"link" binding:"max=255"`
	Sort        int      `json:"sort"`
	Status      int      `json:"status"`
}

type WorkResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	Images      []string  `json:"images"`
	Link        string    `json:"link"`
	Sort        int       `json:"sort"`
	ViewCount   int       `json:"view_count"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

