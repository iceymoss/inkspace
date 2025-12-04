package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:200;not null" json:"title" binding:"required"`
	Content     string         `gorm:"type:longtext;not null" json:"content" binding:"required"`
	Summary     string         `gorm:"size:500" json:"summary"`
	Cover       string         `gorm:"size:255" json:"cover"`
	CategoryID  uint           `gorm:"index" json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag          `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	AuthorID    uint           `gorm:"index;not null" json:"author_id"`
	Author      *User          `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	Status      int            `gorm:"default:1;index" json:"status"` // 1: published, 0: draft
	IsTop       bool           `gorm:"default:false" json:"is_top"`
	IsRecommend bool           `gorm:"default:false" json:"is_recommend"`
}

type ArticleRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary" binding:"max=500"`
	Cover       string   `json:"cover"`
	CategoryID  uint     `json:"category_id"`
	TagIDs      []uint   `json:"tag_ids"`
	Status      int      `json:"status"`
	IsTop       bool     `json:"is_top"`
	IsRecommend bool     `json:"is_recommend"`
}

type ArticleListQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	Keyword    string `form:"keyword"`
	Status     *int   `form:"status"`
}

type ArticleResponse struct {
	ID          uint              `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Summary     string            `json:"summary"`
	Cover       string            `json:"cover"`
	CategoryID  uint              `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Tags        []TagResponse     `json:"tags,omitempty"`
	Author      *UserResponse     `json:"author,omitempty"`
	ViewCount   int               `json:"view_count"`
	LikeCount   int               `json:"like_count"`
	Status      int               `json:"status"`
	IsTop       bool              `json:"is_top"`
	IsRecommend bool              `json:"is_recommend"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (a *Article) ToResponse() *ArticleResponse {
	resp := &ArticleResponse{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Summary:     a.Summary,
		Cover:       a.Cover,
		CategoryID:  a.CategoryID,
		ViewCount:   a.ViewCount,
		LikeCount:   a.LikeCount,
		Status:      a.Status,
		IsTop:       a.IsTop,
		IsRecommend: a.IsRecommend,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	if a.Category != nil {
		resp.Category = a.Category.ToResponse()
	}

	if len(a.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(a.Tags))
		for i, tag := range a.Tags {
			resp.Tags[i] = *tag.ToResponse()
		}
	}

	if a.Author != nil {
		resp.Author = a.Author.ToResponse()
	}

	return resp
}

