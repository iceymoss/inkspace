package models

import (
	"errors"
	"fmt"

	"github.com/iceymoss/inkspace/internal/seo"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// AllocateArticleSlug fills a newly inserted article's NULL slug inside its
// creation transaction. A global unique index arbitrates concurrent writers;
// duplicate-key races advance the deterministic suffix, never overwrite a row.
// Unscoped occupancy reads include soft-deleted rows. The unique index, rather than
// the preliminary snapshot read, is the authority when another writer wins.
func AllocateArticleSlug(tx *gorm.DB, article *Article) error {
	base := seo.SlugBase(article.Title, article.ID)
	for attempt := uint(1); attempt < 1000000; attempt++ {
		candidate := seo.SlugCandidate(base, attempt)
		var occupied uint
		if err := tx.Unscoped().Model(&Article{}).
			Select(ArticleColumns.ID).
			Where(ArticleColumns.Slug+" = ?", candidate).
			Limit(1).
			Scan(&occupied).Error; err != nil {
			return err
		}
		if occupied != 0 {
			continue
		}
		// Slug is otherwise create-only in the GORM model. This narrowly scoped
		// SQL write can fill NULL once, and cannot rewrite an existing slug.
		query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ? AND (%s IS NULL OR %s = '')",
			(Article{}).TableName(), ArticleColumns.Slug, ArticleColumns.ID, ArticleColumns.Slug, ArticleColumns.Slug)
		result := tx.Exec(query, candidate, article.ID)
		if result.Error != nil {
			var duplicate *mysqlDriver.MySQLError
			if errors.As(result.Error, &duplicate) && duplicate.Number == 1062 {
				continue
			}
			return result.Error
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("article %d already has a slug", article.ID)
		}
		article.Slug = candidate
		return nil
	}
	return fmt.Errorf("slug allocation exhausted for article %d", article.ID)
}
