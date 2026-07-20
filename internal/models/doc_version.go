package models

import "time"

type DocVersion struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	DocID     uint      `gorm:"index;not null;uniqueIndex:idx_doc_version" json:"doc_id"`
	Version   int       `gorm:"not null;uniqueIndex:idx_doc_version" json:"version"`
	Title     string    `gorm:"size:200" json:"title"`
	Content   string    `gorm:"type:longtext" json:"content"`
	OwnerID   uint      `gorm:"index" json:"owner_id"`
	Remark    string    `gorm:"size:100" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

type DocVersionResponse struct {
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content,omitempty"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (v *DocVersion) ToResponse() *DocVersionResponse {
	return &DocVersionResponse{
		Version: v.Version, Title: v.Title, Content: v.Content, Remark: v.Remark, CreatedAt: v.CreatedAt,
	}
}
