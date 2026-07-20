package models

import (
	"time"

	"gorm.io/gorm"
)

type Catalog struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	WorkspaceID uint           `gorm:"index;not null" json:"workspace_id"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	OwnerID     uint           `gorm:"index;not null" json:"owner_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Sort        int            `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CatalogCreateRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=100"`
	Sort     int    `json:"sort"`
}

type CatalogUpdateRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

type CatalogMoveRequest struct {
	ParentID *uint `json:"parent_id"`
	Sort     int   `json:"sort"`
}

type CatalogResponse struct {
	ID          uint               `json:"id"`
	WorkspaceID uint               `json:"workspace_id"`
	ParentID    *uint              `json:"parent_id"`
	Name        string             `json:"name"`
	Sort        int                `json:"sort"`
	Children    []*CatalogResponse `json:"children"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (c *Catalog) ToResponse() *CatalogResponse {
	return &CatalogResponse{
		ID: c.ID, WorkspaceID: c.WorkspaceID, ParentID: c.ParentID, Name: c.Name,
		Sort: c.Sort, Children: make([]*CatalogResponse, 0), CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}
