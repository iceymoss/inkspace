package models

import (
	"time"

	"gorm.io/gorm"
)

// ArticleFavorite 文章收藏表
type ArticleFavorite struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"index:idx_user_article,priority:1;not null" json:"user_id"` // 用户ID
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	ArticleID uint           `gorm:"index:idx_user_article,priority:2;index:idx_article_id;not null" json:"article_id"` // 文章ID
	Article   *Article       `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"`
}

// 组合唯一索引：一个用户对同一篇文章只能收藏一次
func (ArticleFavorite) TableName() string {
	return "article_favorites"
}

type FavoriteRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
}

type FavoriteListQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type FavoriteResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	ArticleID uint             `json:"article_id"`
	Article   *ArticleResponse `json:"article,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

func (f *ArticleFavorite) ToResponse() *FavoriteResponse {
	resp := &FavoriteResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		ArticleID: f.ArticleID,
		CreatedAt: f.CreatedAt,
	}

	if f.Article != nil {
		resp.Article = f.Article.ToResponse()
	}

	return resp
}

// Favorite 统一收藏表（支持文章和作品）
type Favorite struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint  `gorm:"not null;index:idx_user_target" json:"user_id"`
	ArticleID *uint `gorm:"index:idx_user_target;index:idx_article" json:"article_id,omitempty"` // 可为空
	WorkID    *uint `gorm:"index:idx_user_target;index:idx_work" json:"work_id,omitempty"`       // 可为空

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Work    *Work    `gorm:"foreignKey:WorkID" json:"work,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

