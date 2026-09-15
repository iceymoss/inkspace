package models

import "time"

type AdminAuditLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	AdminID       uint      `gorm:"index;not null" json:"admin_id"`
	AdminUsername string    `gorm:"size:50" json:"admin_username"`
	Action        string    `gorm:"size:32;index;not null" json:"action"`
	TargetType    string    `gorm:"size:32;index;not null" json:"target_type"`
	TargetID      uint      `gorm:"index;not null" json:"target_id"`
	TargetTitle   string    `gorm:"size:255" json:"target_title"`
	OwnerID       uint      `gorm:"index" json:"owner_id"`
	Reason        string    `gorm:"size:255" json:"reason"`
	IP            string    `gorm:"size:45" json:"ip"`
	CreatedAt     time.Time `json:"created_at"`
}
