package database

import (
	"context"
	"fmt"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"gorm.io/gorm"
)

// ArticleSlugMigrationStats records changes made by an explicit offline migration.
type ArticleSlugMigrationStats struct {
	Scanned, Filled                                 int
	PublishScanned, PublishFilled, PublishPreserved int
}

// MigrateArticleSlugs runs the reviewed slug and first-publication
// prepare/backfill/validate/index sequence.
// Stop all article writers first. The startup advisory lock serializes this
// command against application initialization, but is not an online-write lock.
// Completed rows survive a failure and are preserved on the next invocation.
func MigrateArticleSlugs(ctx context.Context, db *gorm.DB) (stats ArticleSlugMigrationStats, err error) {
	err = withInitializationLock(db.WithContext(ctx), func(conn *gorm.DB) error {
		conn = conn.Session(&gorm.Session{NewDB: true})
		if !conn.Migrator().HasTable(&models.Article{}) {
			return fmt.Errorf("articles table is missing")
		}
		if !conn.Migrator().HasColumn(&models.Article{}, models.ArticleColumns.Slug) {
			if err := conn.Exec("ALTER TABLE articles ADD COLUMN slug varchar(120) CHARACTER SET ascii COLLATE ascii_bin NULL DEFAULT NULL").Error; err != nil {
				return err
			}
		}
		// Validate every existing nonempty slug before changing anything. Include
		// deleted rows, and reserve existing slugs regardless of ID ordering.
		occupied := make(map[string]uint)
		var lastID uint
		for {
			var batch []models.Article
			if err := conn.Unscoped().Select("id", "slug").Where("id > ?", lastID).Order("id ASC").Limit(500).Find(&batch).Error; err != nil {
				return err
			}
			if len(batch) == 0 {
				break
			}
			for _, article := range batch {
				stats.Scanned++
				lastID = article.ID
				if article.Slug == "" {
					continue
				}
				if !seo.ValidSlug(article.Slug) {
					return fmt.Errorf("article %d has invalid slug %q; manual correction required", article.ID, article.Slug)
				}
				if previous, exists := occupied[article.Slug]; exists {
					return fmt.Errorf("duplicate slug %q on articles %d and %d", article.Slug, previous, article.ID)
				}
				occupied[article.Slug] = article.ID
			}
		}
		// Publication repair is part of the same explicit offline rollout, not a
		// potentially large startup task. Preflight all fallback dates first.
		if err := backfillArticlePublication(conn, &stats); err != nil {
			return err
		}
		lastID = 0
		for {
			var batch []models.Article
			if err := conn.Unscoped().Select("id", "title").Where("id > ? AND (slug IS NULL OR slug = '')", lastID).Order("id ASC").Limit(500).Find(&batch).Error; err != nil {
				return err
			}
			if len(batch) == 0 {
				break
			}
			for _, article := range batch {
				lastID = article.ID
				base := seo.SlugBase(article.Title, article.ID)
				candidate := base
				for attempt := uint(2); occupied[candidate] != 0; attempt++ {
					candidate = seo.SlugCandidate(base, attempt)
				}
				result := conn.Exec("UPDATE articles SET slug = ?, updated_at = updated_at WHERE id = ? AND (slug IS NULL OR slug = '')", candidate, article.ID)
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected != 1 {
					return fmt.Errorf("article %d changed during backfill; stop writers", article.ID)
				}
				occupied[candidate] = article.ID
				stats.Filled++
			}
		}
		var invalid int64
		if err := conn.Raw("SELECT COUNT(*) FROM articles WHERE slug IS NULL OR slug = ''").Scan(&invalid).Error; err != nil {
			return err
		}
		if invalid != 0 {
			return fmt.Errorf("%d articles still lack slugs", invalid)
		}
		if !conn.Migrator().HasIndex(&models.Article{}, "ux_articles_slug") {
			if err := conn.Exec("CREATE UNIQUE INDEX ux_articles_slug ON articles(slug)").Error; err != nil {
				return err
			}
		}
		if err := CheckArticleSlugSchema(conn); err != nil {
			return err
		}
		if err := markArticlePublicationReady(conn); err != nil {
			return err
		}
		return CheckArticlePublicationSchema(conn)
	})
	return stats, err
}

// CheckArticleSlugSchema performs a bounded startup check, not historical
// backfill. The explicit migration must finish before old installations start.
func CheckArticleSlugSchema(db *gorm.DB) error {
	var count int64
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.statistics
		WHERE table_schema = DATABASE() AND table_name = 'articles'
		AND index_name = 'ux_articles_slug' AND non_unique = 0 AND column_name = 'slug'
		AND seq_in_index = 1 AND sub_part IS NULL
		AND NOT EXISTS (SELECT 1 FROM information_schema.statistics s
		WHERE s.table_schema = DATABASE() AND s.table_name = 'articles'
		AND s.index_name = 'ux_articles_slug' AND s.seq_in_index > 1)`).Scan(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("article slug migration required: stop writers and run cmd/article-slugs before starting this version")
	}
	var missingID uint
	if err := db.Raw("SELECT id FROM articles WHERE slug IS NULL OR slug = '' LIMIT 1").Scan(&missingID).Error; err != nil {
		return err
	}
	if missingID != 0 {
		return fmt.Errorf("article %d lacks a slug; run cmd/article-slugs offline", missingID)
	}
	return nil
}
