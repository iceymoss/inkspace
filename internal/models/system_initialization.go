package models

import "time"

// SystemInitialization is a private, single-row record; it is not an editable setting.
type SystemInitialization struct {
	ID          uint   `gorm:"primarykey"`
	State       string `gorm:"size:32;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

const (
	InitializationPending      = "pending"
	InitializationLegacyClosed = "legacy_closed"
	InitializationCompleted    = "completed"
)

type AdminBootstrapRequest struct {
	Secret          string `json:"secret" binding:"required,max=256"`
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email,max=100"`
	Password        string `json:"password" binding:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type AdminBootstrapStatus struct {
	State            string `json:"state"`
	SecretConfigured bool   `json:"secret_configured"`
}
