package service

import (
	"context"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

// GetRecentPublicArticles returns at most 50 fresh published, nondeleted articles
// ordered by first-publication time DESC, ID DESC (not pinning or edit time).
// The projection contains feed text/dates/slug and Author's public name only;
// callers must not serialize models.User as a private user response. This read
// has no view-count/cache side effects and preserves query/cancellation errors.
func (s *ArticleService) GetRecentPublicArticles(ctx context.Context) ([]*models.Article, error) {
	articles := make([]*models.Article, 0)
	err := PublicArticles(database.DB.WithContext(ctx)).
		Select([]string{
			models.ArticleColumns.ID,
			models.ArticleColumns.Slug,
			models.ArticleColumns.Title,
			models.ArticleColumns.Summary,
			models.ArticleColumns.Content,
			models.ArticleColumns.AuthorID,
			models.ArticleColumns.PublishAt,
			models.ArticleColumns.UpdatedAt,
			models.ArticleColumns.CreatedAt,
		}).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "nickname")
		}).
		Order(models.ArticleColumns.PublishAt + " DESC").
		Order(models.ArticleColumns.ID + " DESC").
		Limit(50).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// PublicSitemapArticle is the compact sitemap projection. ID order, unlike edit
// or publication order, does not move records between shards after an edit.
type PublicSitemapArticle struct {
	ID        uint
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
	PublishAt *time.Time
}

// LastModified returns the content/state timestamp, then first publication, then
// creation time. Zero means no reliable timestamp is available and may be omitted.
func (a PublicSitemapArticle) LastModified() time.Time {
	if !a.UpdatedAt.IsZero() {
		return a.UpdatedAt
	}
	if a.PublishAt != nil && !a.PublishAt.IsZero() {
		return *a.PublishAt
	}
	return a.CreatedAt
}

// GetPublicSitemapArticles returns every public article's compact projection in
// ascending ID order using one SELECT. len(result) is the matching count from the
// same statement snapshot; consumers may split it deterministically into shards.
// No full content, association, user/admin filters or cache are involved. Separate
// HTTP requests observe fresh visibility, not a cross-request snapshot; consumers
// must not promise immutable shard membership while articles are published/deleted.
func (s *ArticleService) GetPublicSitemapArticles(ctx context.Context) ([]PublicSitemapArticle, error) {
	articles := make([]PublicSitemapArticle, 0)
	err := PublicArticles(database.DB.WithContext(ctx)).
		Select([]string{
			models.ArticleColumns.ID,
			models.ArticleColumns.Slug,
			models.ArticleColumns.CreatedAt,
			models.ArticleColumns.UpdatedAt,
			models.ArticleColumns.PublishAt,
		}).
		Order(models.ArticleColumns.ID + " ASC").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
