package models

import (
	"time"

	"gorm.io/gorm"
)

// Subscription 订阅表
type Subscription struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"` // 订阅邮箱
	Status    int            `gorm:"default:1;index" json:"status"` // 1: 已确认, 0: 未确认, -1: 已取消
	Token     string         `gorm:"uniqueIndex;size:64" json:"-"` // 确认/取消令牌
	IP        string         `gorm:"size:50" json:"ip"` // 订阅时的IP
	UserAgent string         `gorm:"size:255" json:"user_agent"` // User Agent
	ConfirmAt *time.Time     `gorm:"type:datetime(3)" json:"confirm_at"` // 确认时间
}

type SubscriptionRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SubscriptionResponse struct {
	ID        uint       `json:"id"`
	Email     string     `json:"email"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	ConfirmAt *time.Time `json:"confirm_at"`
}

func (s *Subscription) ToResponse() *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:        s.ID,
		Email:     s.Email,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		ConfirmAt: s.ConfirmAt,
	}
}

