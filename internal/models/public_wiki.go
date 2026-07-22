package models

import "time"

type PublicWorkspaceResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Icon         string    `json:"icon"`
	DocCount     int64     `json:"doc_count"`
	AuthorID     uint      `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PublicWorkspaceSummary struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type PublicCatalogResponse struct {
	ID       uint                     `json:"id"`
	ParentID *uint                    `json:"parent_id"`
	Name     string                   `json:"name"`
	Sort     int                      `json:"sort"`
	Children []*PublicCatalogResponse `json:"children"`
}

type PublicDocTreeResponse struct {
	ID          uint       `json:"id"`
	CatalogID   *uint      `json:"catalog_id"`
	Title       string     `json:"title"`
	Sort        int        `json:"sort"`
	PublishedAt *time.Time `json:"published_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PublicWorkspaceTreeResponse struct {
	Workspace *PublicWorkspaceSummary  `json:"workspace"`
	Catalogs  []*PublicCatalogResponse `json:"catalogs"`
	Docs      []*PublicDocTreeResponse `json:"docs"`
}

type PublicDocResponse struct {
	ID          uint       `json:"id"`
	WorkspaceID uint       `json:"workspace_id"`
	CatalogID   *uint      `json:"catalog_id"`
	Title       string     `json:"title"`
	ContentHTML string     `json:"content_html"`
	ViewCount   int        `json:"view_count"`
	PublishedAt *time.Time `json:"published_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
