package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

// Read textual dates rather than time.Time: MySQL's parseTime conversion can
// normalize an invalid legacy date (e.g. February 31) into a plausible Go date.
type articlePublicationRow struct {
	ID                 uint
	Created, Published sql.NullString
}

func validArticleDate(raw sql.NullString) bool {
	if !raw.Valid {
		return false
	}
	at, err := time.Parse("2006-01-02 15:04:05.999999999", raw.String)
	return err == nil && at.Year() >= 1000
}

func publicationBatch(db *gorm.DB, after uint) ([]articlePublicationRow, error) {
	var batch []articlePublicationRow
	err := db.Raw(`SELECT id, CAST(created_at AS CHAR) AS created,
		CAST(publish_at AS CHAR) AS published FROM articles WHERE id > ? ORDER BY id ASC LIMIT 500`, after).
		Scan(&batch).Error
	return batch, err
}

func backfillArticlePublication(db *gorm.DB, stats *ArticleSlugMigrationStats) error {
	if !db.Migrator().HasColumn(&models.Article{}, models.ArticleColumns.CreatedAt) {
		return fmt.Errorf("articles.created_at is missing; publication backfill requires a valid historical fallback")
	}
	if !db.Migrator().HasColumn(&models.Article{}, models.ArticleColumns.PublishAt) {
		if err := db.Exec("ALTER TABLE articles ADD COLUMN publish_at datetime(3) NULL DEFAULT NULL").Error; err != nil {
			return err
		}
	}
	var lastID uint
	for {
		batch, err := publicationBatch(db, lastID)
		if err != nil {
			return err
		}
		if len(batch) == 0 {
			break
		}
		for _, row := range batch {
			lastID = row.ID
			stats.PublishScanned++
			if validArticleDate(row.Published) {
				stats.PublishPreserved++
				continue
			}
			if !validArticleDate(row.Created) {
				return fmt.Errorf("article %d has invalid publish_at and created_at; manual correction required", row.ID)
			}
		}
	}
	lastID = 0
	for {
		batch, err := publicationBatch(db, lastID)
		if err != nil {
			return err
		}
		if len(batch) == 0 {
			return nil
		}
		for _, row := range batch {
			lastID = row.ID
			if validArticleDate(row.Published) {
				continue
			}
			// Raw UPDATE avoids GORM timestamps and explicitly preserves MySQL
			// legacy ON UPDATE columns. Include deleted/draft/private history too.
			result := db.Exec(`UPDATE articles SET publish_at = created_at, updated_at = updated_at
				WHERE id = ? AND CAST(publish_at AS CHAR) <=> ?`, row.ID, row.Published)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return fmt.Errorf("article %d changed during publication backfill; stop writers", row.ID)
			}
			stats.PublishFilled++
		}
	}
}

// The dedicated nonunique index is also the completion marker for this offline
// data migration. Startup checks metadata only, never scans historical dates.
func markArticlePublicationReady(db *gorm.DB) error {
	if db.Migrator().HasIndex(&models.Article{}, "idx_articles_first_publication_v1") {
		return CheckArticlePublicationSchema(db)
	}
	return db.Exec("CREATE INDEX idx_articles_first_publication_v1 ON articles(status, deleted_at, publish_at DESC, id DESC)").Error
}

// CheckArticlePublicationSchema refuses pre-backfill deployments using the
// migration's supporting index as a bounded readiness marker. It does not perform
// historical scans or repair and must not be substituted for the offline command.
func CheckArticlePublicationSchema(db *gorm.DB) error {
	var columns int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.statistics
		WHERE table_schema = DATABASE() AND table_name = 'articles'
		AND index_name = 'idx_articles_first_publication_v1' AND non_unique = 1 AND sub_part IS NULL
		AND ((seq_in_index = 1 AND column_name = 'status' AND collation = 'A')
		OR (seq_in_index = 2 AND column_name = 'deleted_at' AND collation = 'A')
		OR (seq_in_index = 3 AND column_name = 'publish_at' AND collation = 'D')
		OR (seq_in_index = 4 AND column_name = 'id' AND collation = 'D'))
		AND NOT EXISTS (SELECT 1 FROM information_schema.statistics s
		WHERE s.table_schema = DATABASE() AND s.table_name = 'articles'
		AND s.index_name = 'idx_articles_first_publication_v1' AND s.seq_in_index > 4)`).Scan(&columns).Error; err != nil {
		return err
	}
	if columns != 4 {
		return fmt.Errorf("article publication migration required: stop writers and run cmd/article-slugs before starting this version")
	}
	return nil
}
