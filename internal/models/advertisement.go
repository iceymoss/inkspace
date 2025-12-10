package models

import (
	"time"

	"gorm.io/gorm"
)

// Advertisement 广告消息表
type Advertisement struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Title       string `gorm:"type:varchar(200);not null" json:"title"`       // 广告标题
	Image       string `gorm:"type:varchar(500)" json:"image"`                // 广告图片URL
	Link        string `gorm:"type:varchar(500)" json:"link"`                // 广告链接
	Description string `gorm:"type:text" json:"description"`                 // 广告描述
	Status      int    `gorm:"default:1;comment:1-启用,0-禁用" json:"status"` // 状态：1-启用，0-禁用
	ClickCount  int    `gorm:"default:0" json:"click_count"`                  // 点击次数
	ViewCount   int    `gorm:"default:0" json:"view_count"`                   // 展示次数
}

func (Advertisement) TableName() string {
	return "advertisements"
}

type AdvertisementResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	ClickCount  int    `json:"click_count"`
	ViewCount   int    `json:"view_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (a *Advertisement) ToResponse() *AdvertisementResponse {
	return &AdvertisementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Image:       a.Image,
		Link:        a.Link,
		Description: a.Description,
		Status:      a.Status,
		ClickCount:  a.ClickCount,
		ViewCount:   a.ViewCount,
		CreatedAt:   a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// AdvertisementRequest 广告消息请求
type AdvertisementRequest struct {
	Title       string `json:"title" binding:"required"`
	Image       string `json:"image"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

