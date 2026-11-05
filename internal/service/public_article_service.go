package service

import (
	"context"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"gorm.io/gorm"
)

// PublicArticles is the shared visibility scope for HTML, slug API and discovery
// queries. It deliberately accepts no user-supplied list or authorization filters.
func PublicArticles(db *gorm.DB) *gorm.DB {
	return db.Model(&models.Article{}).
		Where(models.ArticleColumns.Status+" = ?", 1).
		Where(models.ArticleColumns.DeletedAt + " IS NULL")
}

// GetPublicBySlug reads fresh visibility state, never the internal ID cache.
// Unpublished, deleted and missing articles all return gorm.ErrRecordNotFound.
func (s *ArticleService) GetPublicBySlug(ctx context.Context, slug string) (*models.Article, error) {
	if !seo.ValidSlug(slug) {
		return nil, gorm.ErrRecordNotFound
	}
	return findPublicArticle(PublicArticles(database.DB.WithContext(ctx)).Where(models.ArticleColumns.Slug+" = ?", slug))
}

// GetPublicByID resolves legacy links without granting author/admin exceptions.
func (s *ArticleService) GetPublicByID(ctx context.Context, id uint) (*models.Article, error) {
	return findPublicArticle(PublicArticles(database.DB.WithContext(ctx)).Where(models.ArticleColumns.ID+" = ?", id))
}

func findPublicArticle(query *gorm.DB) (*models.Article, error) {
	var article models.Article
	if err := query.Preload("Author").Preload("Category").Preload("Tags").First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// RecordPublicArticleView preserves the numeric API's interaction behavior while
// excluding articles that have just become nonpublic. It never updates updated_at.
func (s *ArticleService) RecordPublicArticleView(ctx context.Context, id uint) error {
	return PublicArticles(database.DB.WithContext(ctx)).
		Where(models.ArticleColumns.ID+" = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
