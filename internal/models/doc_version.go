package models

import "time"

type DocVersion struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	DocID        uint      `gorm:"index;not null;uniqueIndex:idx_doc_version" json:"doc_id"`
	Version      int       `gorm:"not null;uniqueIndex:idx_doc_version" json:"version"`
	Title        string    `gorm:"size:200" json:"title"`
	Content      string    `gorm:"type:longtext" json:"content"`
	Revision     uint64    `gorm:"default:1;not null" json:"revision"`
	Kind         string    `gorm:"size:20;default:'markdown'" json:"kind"`
	Language     string    `gorm:"size:50" json:"language"`
	AttachmentID *uint     `gorm:"index" json:"attachment_id"`
	OwnerID      uint      `gorm:"index" json:"owner_id"`
	Remark       string    `gorm:"size:100" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

type DocVersionResponse struct {
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content,omitempty"`
	Revision  uint64    `json:"revision"`
	Kind      string    `json:"kind"`
	Language  string    `json:"language"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (v *DocVersion) ToResponse() *DocVersionResponse {
	kind := v.Kind
	if kind == "" {
		kind = DocKindMarkdown
	}
	revision := v.Revision
	if revision == 0 {
		revision = 1
	}
	return &DocVersionResponse{
		Version: v.Version, Title: v.Title, Content: v.Content, Revision: revision,
		Kind: kind, Language: v.Language, Remark: v.Remark, CreatedAt: v.CreatedAt,
	}
}
