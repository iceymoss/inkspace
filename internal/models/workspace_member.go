package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	WorkspaceRoleOwner  = "owner"
	WorkspaceRoleAdmin  = "admin"
	WorkspaceRoleEditor = "editor"
	WorkspaceRoleViewer = "viewer"
)

type WorkspaceMember struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	WorkspaceID uint           `gorm:"uniqueIndex:idx_workspace_user;not null" json:"workspace_id"`
	UserID      uint           `gorm:"uniqueIndex:idx_workspace_user;not null" json:"user_id"`
	Role        string         `gorm:"size:20;index;not null" json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WorkspaceMemberCreateRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Role     string `json:"role" binding:"required,oneof=admin editor viewer"`
}

type WorkspaceMemberUpdateRequest struct {
	Role string `json:"role" binding:"required,oneof=admin editor viewer"`
}

type WorkspaceMemberResponse struct {
	WorkspaceID uint      `json:"workspace_id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
