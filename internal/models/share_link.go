package models

import (
	"time"

	"gorm.io/gorm"
)

type ShareLink struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Token     string         `gorm:"uniqueIndex;size:32;not null" json:"token"`
	DocID     uint           `gorm:"index;not null" json:"doc_id"`
	OwnerID   uint           `gorm:"index;not null" json:"owner_id"`
	ExpiresAt *time.Time     `json:"expires_at"`
	Enabled   bool           `gorm:"default:true;not null" json:"enabled"`
	ViewCount int            `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ShareLinkCreateRequest struct {
	Permanent bool       `json:"permanent"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type ShareLinkUpdateRequest struct {
	Enabled   *bool      `json:"enabled"`
	Permanent bool       `json:"permanent"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type ShareLinkResponse struct {
	ID        uint       `json:"id"`
	Token     string     `json:"token"`
	DocID     uint       `json:"doc_id"`
	ExpiresAt *time.Time `json:"expires_at"`
	Enabled   bool       `json:"enabled"`
	ViewCount int        `json:"view_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *ShareLink) ToResponse() *ShareLinkResponse {
	return &ShareLinkResponse{
		ID: s.ID, Token: s.Token, DocID: s.DocID, ExpiresAt: s.ExpiresAt,
		Enabled: s.Enabled, ViewCount: s.ViewCount, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
