package models

import (
	"time"

	"gorm.io/gorm"
)

// AdPosition 广告位置表
type AdPosition struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(100);not null" json:"name"`        // 位置名称，如"博客列表右侧"
	Code        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // 位置代码，如"blog_right_ad"
	Description string `gorm:"type:text" json:"description"`                     // 位置描述
	MaxCount    int    `gorm:"default:4;comment:最大投放数量" json:"max_count"`  // 最大投放数量，默认4个
	Status      int    `gorm:"default:1;comment:1-启用,0-禁用" json:"status"`    // 状态：1-启用，0-禁用
	Sort        int    `gorm:"default:0" json:"sort"`                            // 排序
}

func (AdPosition) TableName() string {
	return "ad_positions"
}

type AdPositionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	MaxCount    int    `json:"max_count"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (p *AdPosition) ToResponse() *AdPositionResponse {
	return &AdPositionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		MaxCount:    p.MaxCount,
		Status:      p.Status,
		Sort:        p.Sort,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// AdPositionRequest 广告位置请求
type AdPositionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	MaxCount    int    `json:"max_count"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
}

