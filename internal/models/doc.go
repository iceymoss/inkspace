package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	DocStatusDraft     = 0
	DocStatusPublished = 1

	DocKindMarkdown = "markdown"
	DocKindText     = "text"
	DocKindCode     = "code"
	DocKindFile     = "file"
)

type Doc struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	WorkspaceID           uint           `gorm:"index;not null" json:"workspace_id"`
	CatalogID             *uint          `gorm:"index" json:"catalog_id"`
	ArticleID             *uint          `gorm:"index" json:"article_id"`
	OwnerID               uint           `gorm:"index;not null" json:"owner_id"`
	Title                 string         `gorm:"size:200;not null" json:"title"`
	Content               string         `gorm:"type:longtext" json:"content"`
	ContentHTML           string         `gorm:"type:longtext" json:"content_html"`
	Kind                  string         `gorm:"size:20;index;default:'markdown';not null" json:"kind"`
	Language              string         `gorm:"size:50" json:"language"`
	AttachmentID          *uint          `gorm:"index" json:"attachment_id"`
	PublishedAttachmentID *uint          `gorm:"index" json:"published_attachment_id"`
	PublishedRevision     uint64         `gorm:"default:0;not null" json:"published_revision"`
	Revision              uint64         `gorm:"default:1;not null" json:"revision"`
	Status                int            `gorm:"index;default:0;not null" json:"status"`
	WordCount             int            `gorm:"default:0" json:"word_count"`
	ViewCount             int            `gorm:"default:0" json:"view_count"`
	Sort                  int            `gorm:"default:0" json:"sort"`
	PublishedAt           *time.Time     `json:"published_at"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

type DocCreateRequest struct {
	WorkspaceID uint   `json:"workspace_id" binding:"required"`
	CatalogID   *uint  `json:"catalog_id"`
	Title       string `json:"title" binding:"required,max=200"`
	FileName    string `json:"file_name" binding:"omitempty,max=200"`
	Content     string `json:"content"`
	Kind        string `json:"kind" binding:"omitempty,oneof=markdown text code"`
	Language    string `json:"language" binding:"max=50"`
	Sort        int    `json:"sort"`
}

type DocSaveRequest struct {
	Title    string `json:"title" binding:"required,max=200"`
	Content  string `json:"content"`
	Revision uint64 `json:"revision" binding:"required"`
}

type DocAutosaveRequest struct {
	Content  string `json:"content"`
	Revision uint64 `json:"revision" binding:"required"`
}

type DocMoveRequest struct {
	CatalogID *uint `json:"catalog_id"`
	Sort      int   `json:"sort"`
}

type DocPublishRequest struct {
	Status *int `json:"status" binding:"omitempty,oneof=0 1"`
}

type DocPublishToBlogRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
	Summary    string `json:"summary" binding:"max=500"`
	Cover      string `json:"cover" binding:"max=255"`
}

type DocListQuery struct {
	CatalogID *uint `form:"catalog_id"`
}

type KnowledgeFileUploadRequest struct {
	CatalogID *uint  `form:"catalog_id"`
	Title     string `form:"title" binding:"omitempty,max=200"`
	RequestID string `form:"request_id" binding:"omitempty,max=100"`
}

type DocResponse struct {
	ID          uint       `json:"id"`
	WorkspaceID uint       `json:"workspace_id"`
	CatalogID   *uint      `json:"catalog_id"`
	ArticleID   *uint      `json:"article_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary,omitempty"`
	Content     string     `json:"content"`
	ContentHTML string     `json:"content_html"`
	Kind        string     `json:"kind"`
	Language    string     `json:"language"`
	Editable    bool       `json:"editable"`
	Revision    uint64     `json:"revision"`
	Status      int        `json:"status"`
	WordCount   int        `json:"word_count"`
	ViewCount   int        `json:"view_count"`
	Sort        int        `json:"sort"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type DocCapabilities struct {
	CanView          bool `json:"can_view"`
	CanEdit          bool `json:"can_edit"`
	CanReplace       bool `json:"can_replace"`
	CanDownload      bool `json:"can_download"`
	CanPublish       bool `json:"can_publish"`
	CanShare         bool `json:"can_share"`
	CanDelete        bool `json:"can_delete"`
	CanManageMembers bool `json:"can_manage_members"`
}

// DocAttachmentDetail intentionally excludes storage paths and URLs.
type DocAttachmentDetail struct {
	ID              uint   `json:"id"`
	FileName        string `json:"file_name"`
	FileSize        int64  `json:"file_size"`
	MimeType        string `json:"mime_type"`
	Extension       string `json:"extension"`
	Checksum        string `json:"checksum"`
	PreviewStatus   string `json:"preview_status"`
	PreviewMimeType string `json:"preview_mime_type"`
}

type DocDetailResponse struct {
	ID           uint                 `json:"id"`
	WorkspaceID  uint                 `json:"workspace_id"`
	CatalogID    *uint                `json:"catalog_id"`
	ArticleID    *uint                `json:"article_id"`
	Title        string               `json:"title"`
	Content      string               `json:"content,omitempty"`
	ContentHTML  string               `json:"content_html,omitempty"`
	Kind         string               `json:"kind"`
	Language     string               `json:"language"`
	Revision     uint64               `json:"revision"`
	Status       int                  `json:"status"`
	WordCount    int                  `json:"word_count"`
	ViewCount    int                  `json:"view_count"`
	PublishedAt  *time.Time           `json:"published_at"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Attachment   *DocAttachmentDetail `json:"attachment"`
	Preview      interface{}          `json:"preview"`
	Capabilities DocCapabilities      `json:"capabilities"`
}

type DocSearchResponse struct {
	ID        uint      `json:"id"`
	CatalogID *uint     `json:"catalog_id"`
	ArticleID *uint     `json:"article_id"`
	Title     string    `json:"title"`
	Kind      string    `json:"kind"`
	Language  string    `json:"language"`
	Editable  bool      `json:"editable"`
	Summary   string    `json:"summary"`
	Status    int       `json:"status"`
	WordCount int       `json:"word_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Doc) ToResponse() *DocResponse {
	kind := d.Kind
	if kind == "" {
		kind = DocKindMarkdown
	}
	revision := d.Revision
	if revision == 0 {
		revision = 1
	}
	return &DocResponse{
		ID: d.ID, WorkspaceID: d.WorkspaceID, CatalogID: d.CatalogID, ArticleID: d.ArticleID, Title: d.Title,
		Content: d.Content, ContentHTML: d.ContentHTML, Kind: kind, Language: d.Language, Revision: revision,
		Status: d.Status, WordCount: d.WordCount,
		ViewCount: d.ViewCount, Sort: d.Sort, PublishedAt: d.PublishedAt,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}
