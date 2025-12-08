package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// CustomTime 自定义时间类型，支持多种时间格式解析
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON反序列化，支持多种时间格式
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return nil
	}

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			ct.Time = t
			return nil
		}
	}

	// 如果所有格式都失败，返回错误
	return &time.ParseError{
		Value: s,
		Layout: "multiple formats",
		Message: "unable to parse time",
	}
}

// MarshalJSON 自定义JSON序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + ct.Time.Format(time.RFC3339) + `"`), nil
}

// AdPlacement 广告投放表
type AdPlacement struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	PositionID    uint      `gorm:"index:idx_position;not null" json:"position_id"`    // 广告位置ID
	AdvertisementID uint    `gorm:"index:idx_advertisement;not null" json:"advertisement_id"` // 广告消息ID
	StartTime     *time.Time `gorm:"type:datetime" json:"start_time"`                 // 投放开始时间
	EndTime       *time.Time `gorm:"type:datetime" json:"end_time"`                   // 投放结束时间
	Status        int        `gorm:"default:1;comment:1-启用,0-禁用" json:"status"`     // 状态：1-启用，0-禁用
	Sort          int        `gorm:"default:0" json:"sort"`                           // 排序（同一位置多个广告时使用）
	ClickCount    int        `gorm:"default:0" json:"click_count"`                    // 点击次数
	ViewCount     int        `gorm:"default:0" json:"view_count"`                      // 展示次数

	// 关联关系
	Position      *AdPosition   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Advertisement *Advertisement `gorm:"foreignKey:AdvertisementID" json:"advertisement,omitempty"`
}

func (AdPlacement) TableName() string {
	return "ad_placements"
}

type AdPlacementResponse struct {
	ID              uint                   `json:"id"`
	PositionID      uint                   `json:"position_id"`
	AdvertisementID uint                   `json:"advertisement_id"`
	StartTime       *string                `json:"start_time"`
	EndTime         *string                 `json:"end_time"`
	Status          int                    `json:"status"`
	Sort            int                    `json:"sort"`
	ClickCount      int                    `json:"click_count"`
	ViewCount       int                    `json:"view_count"`
	Position        *AdPositionResponse    `json:"position,omitempty"`
	Advertisement   *AdvertisementResponse `json:"advertisement,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

func (p *AdPlacement) ToResponse() *AdPlacementResponse {
	resp := &AdPlacementResponse{
		ID:              p.ID,
		PositionID:      p.PositionID,
		AdvertisementID: p.AdvertisementID,
		Status:          p.Status,
		Sort:            p.Sort,
		ClickCount:      p.ClickCount,
		ViewCount:       p.ViewCount,
		CreatedAt:       p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if p.StartTime != nil {
		startTimeStr := p.StartTime.Format("2006-01-02 15:04:05")
		resp.StartTime = &startTimeStr
	}
	if p.EndTime != nil {
		endTimeStr := p.EndTime.Format("2006-01-02 15:04:05")
		resp.EndTime = &endTimeStr
	}

	if p.Position != nil {
		resp.Position = p.Position.ToResponse()
	}
	if p.Advertisement != nil {
		resp.Advertisement = p.Advertisement.ToResponse()
	}

	return resp
}

// AdDisplayResponse 前端展示用的广告响应（简化版）
type AdDisplayResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

// AdPlacementRequest 广告投放请求
type AdPlacementRequest struct {
	PositionID      uint        `json:"position_id" binding:"required"`
	AdvertisementID uint        `json:"advertisement_id" binding:"required"`
	StartTime       *CustomTime `json:"start_time"`
	EndTime         *CustomTime `json:"end_time"`
	Status          int         `json:"status"`
	Sort            int         `json:"sort"`
}

