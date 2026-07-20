package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	DocStatusDraft     = 0
	DocStatusPublished = 1
)

type Doc struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	WorkspaceID uint           `gorm:"index;not null" json:"workspace_id"`
	CatalogID   *uint          `gorm:"index" json:"catalog_id"`
	ArticleID   *uint          `gorm:"index" json:"article_id"`
	OwnerID     uint           `gorm:"index;not null" json:"owner_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Content     string         `gorm:"type:longtext" json:"content"`
	ContentHTML string         `gorm:"type:longtext" json:"content_html"`
	Status      int            `gorm:"index;default:0;not null" json:"status"`
	WordCount   int            `gorm:"default:0" json:"word_count"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	Sort        int            `gorm:"default:0" json:"sort"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type DocCreateRequest struct {
	WorkspaceID uint   `json:"workspace_id" binding:"required"`
	CatalogID   *uint  `json:"catalog_id"`
	Title       string `json:"title" binding:"required,max=200"`
	Content     string `json:"content"`
	Sort        int    `json:"sort"`
}

type DocSaveRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content"`
}

type DocAutosaveRequest struct {
	Content string `json:"content"`
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

type DocResponse struct {
	ID          uint       `json:"id"`
	WorkspaceID uint       `json:"workspace_id"`
	CatalogID   *uint      `json:"catalog_id"`
	ArticleID   *uint      `json:"article_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary,omitempty"`
	Content     string     `json:"content"`
	ContentHTML string     `json:"content_html"`
	Status      int        `json:"status"`
	WordCount   int        `json:"word_count"`
	ViewCount   int        `json:"view_count"`
	Sort        int        `json:"sort"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type DocSearchResponse struct {
	ID        uint      `json:"id"`
	CatalogID *uint     `json:"catalog_id"`
	ArticleID *uint     `json:"article_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Status    int       `json:"status"`
	WordCount int       `json:"word_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Doc) ToResponse() *DocResponse {
	return &DocResponse{
		ID: d.ID, WorkspaceID: d.WorkspaceID, CatalogID: d.CatalogID, ArticleID: d.ArticleID, Title: d.Title,
		Content: d.Content, ContentHTML: d.ContentHTML, Status: d.Status, WordCount: d.WordCount,
		ViewCount: d.ViewCount, Sort: d.Sort, PublishedAt: d.PublishedAt,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}
