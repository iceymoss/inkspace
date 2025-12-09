package models

import (
	"time"

	"gorm.io/gorm"
)

// Link 友情链接表
type Link struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name" binding:"required"` // 链接名称
	URL         string         `gorm:"size:255;not null" json:"url" binding:"required,url"` // 链接地址
	Logo        string         `gorm:"size:255" json:"logo"` // 网站Logo
	Description string         `gorm:"size:200" json:"description"` // 描述
	Email       string         `gorm:"size:100" json:"email"` // 联系邮箱
	Sort        int            `gorm:"default:0;index:idx_status_sort" json:"sort"` // 排序
	Status      int            `gorm:"default:1;index:idx_status_sort" json:"status"` // 1: 显示, 0: 隐藏
}

type LinkRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	URL         string `json:"url" binding:"required,url"`
	Logo        string `json:"logo"`
	Description string `json:"description" binding:"max=200"`
	Email       string `json:"email" binding:"omitempty,email"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

type LinkResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Logo        string    `json:"logo"`
	Description string    `json:"description"`
	Sort        int       `json:"sort"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (l *Link) ToResponse() *LinkResponse {
	return &LinkResponse{
		ID:          l.ID,
		Name:        l.Name,
		URL:         l.URL,
		Logo:        l.Logo,
		Description: l.Description,
		Sort:        l.Sort,
		Status:      l.Status,
		CreatedAt:   l.CreatedAt,
	}
}

