package models

import (
	"time"

	"gorm.io/gorm"
)

type Workspace struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	OwnerID     uint           `gorm:"index;not null" json:"owner_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Icon        string         `gorm:"size:255" json:"icon"`
	Sort        int            `gorm:"default:0" json:"sort"`
	DocCount    int            `gorm:"default:0;not null" json:"doc_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WorkspaceRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=255"`
	Sort        int    `json:"sort"`
}

type WorkspaceResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Sort        int       `json:"sort"`
	DocCount    int       `json:"doc_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (w *Workspace) ToResponse() *WorkspaceResponse {
	return &WorkspaceResponse{
		ID: w.ID, Name: w.Name, Description: w.Description, Icon: w.Icon,
		Sort: w.Sort, DocCount: w.DocCount, CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
	}
}
