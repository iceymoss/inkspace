package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Title        string         `gorm:"size:200;not null;index:idx_title" json:"title" binding:"required"`
	Content      string         `gorm:"type:longtext;not null" json:"content" binding:"required"`
	ContentHTML  string         `gorm:"type:longtext" json:"content_html"`
	Summary      string         `gorm:"size:500" json:"summary"`
	Cover        string         `gorm:"size:255" json:"cover"`
	CategoryID   uint           `gorm:"index:idx_category_id" json:"category_id"`
	Category     *Category      `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Tags         []Tag          `gorm:"many2many:article_tags;constraint:OnDelete:CASCADE" json:"tags,omitempty"`
	AuthorID     uint           `gorm:"index:idx_author_id;not null" json:"author_id"`
	Author       *User          `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author,omitempty"`
	ViewCount    int            `gorm:"default:0;not null" json:"view_count"`
	LikeCount     int            `gorm:"default:0;not null" json:"like_count"`
	CommentCount  int            `gorm:"default:0;not null" json:"comment_count"`
	FavoriteCount int            `gorm:"default:0;not null" json:"favorite_count"` // 收藏数
	WordCount     int            `gorm:"default:0;not null" json:"word_count"`
	ReadingTime  int            `gorm:"default:0;not null" json:"reading_time"` // 阅读时间（分钟）
	Status       int            `gorm:"default:1;index:idx_status;index:idx_top_status_created" json:"status"` // 1: published, 0: draft
	IsTop        bool           `gorm:"default:false;index:idx_top_status_created" json:"is_top"`
	IsRecommend  bool           `gorm:"default:false" json:"is_recommend"`
	IsOriginal   bool           `gorm:"default:true" json:"is_original"`
	SourceURL    string         `gorm:"size:255" json:"source_url"`
	PublishAt    *time.Time     `gorm:"type:datetime(3)" json:"publish_at"`
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
	AuthorID   uint   `form:"author_id"`
}

type ArticleResponse struct {
	ID            uint              `json:"id"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	Summary       string            `json:"summary"`
	Cover         string            `json:"cover"`
	CategoryID    uint              `json:"category_id"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Tags          []TagResponse     `json:"tags,omitempty"`
	AuthorID      uint              `json:"author_id"`
	Author        *UserResponse     `json:"author,omitempty"`
	ViewCount     int               `json:"view_count"`
	LikeCount     int               `json:"like_count"`
	CommentCount  int               `json:"comment_count"`
	FavoriteCount int               `json:"favorite_count"`
	WordCount     int               `json:"word_count"`
	ReadingTime   int               `json:"reading_time"`
	Status        int               `json:"status"`
	IsTop         bool              `json:"is_top"`
	IsRecommend   bool              `json:"is_recommend"`
	IsOriginal    bool              `json:"is_original"`
	SourceURL     string            `json:"source_url"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func (a *Article) ToResponse() *ArticleResponse {
	resp := &ArticleResponse{
		ID:            a.ID,
		Title:         a.Title,
		Content:       a.Content,
		Summary:       a.Summary,
		Cover:         a.Cover,
		CategoryID:    a.CategoryID,
		AuthorID:      a.AuthorID,
		ViewCount:     a.ViewCount,
		LikeCount:     a.LikeCount,
		CommentCount:  a.CommentCount,
		FavoriteCount: a.FavoriteCount,
		WordCount:     a.WordCount,
		ReadingTime:   a.ReadingTime,
		Status:        a.Status,
		IsTop:         a.IsTop,
		IsRecommend:   a.IsRecommend,
		IsOriginal:    a.IsOriginal,
		SourceURL:     a.SourceURL,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
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

