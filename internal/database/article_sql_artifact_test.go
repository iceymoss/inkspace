package database

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/testutil"

	"gorm.io/gorm"
)

// runArticleSQL executes mysql-client DELIMITER artifacts only inside the
// explicitly opted-in random SEOTestMySQL schema, never application config.
func runArticleSQL(t *testing.T, db *gorm.DB, name string) error {
	t.Helper()
	data, err := os.ReadFile("../../scripts/article-slugs/" + name)
	if err != nil {
		t.Fatal(err)
	}
	delimiter := ";"
	var statement strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "DELIMITER ") {
			delimiter = strings.TrimPrefix(line, "DELIMITER ")
			continue
		}
		statement.WriteString(line + "\n")
		if strings.HasSuffix(line, delimiter) {
			sql := strings.TrimSuffix(strings.TrimSpace(statement.String()), delimiter)
			if err := db.Exec(sql).Error; err != nil {
				return err
			}
			statement.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(statement.String()) != "" {
		t.Fatal("unterminated SQL artifact")
	}
	return nil
}

func TestArticleSQLPhasesMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	var version string
	if err := db.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		t.Fatal(err)
	}
	t.Logf("SQL artifact acceptance server: MySQL %s", version)
	if err := db.Connection(func(conn *gorm.DB) error {
		// Connection passes an initialized GORM handle. Clone each statement so
		// an expected guard error cannot poison later statements on this session.
		conn = conn.Session(&gorm.Session{NewDB: true})
		if err := conn.Exec(`CREATE TABLE articles (id bigint unsigned PRIMARY KEY, title varchar(200), status int,
			deleted_at datetime(3), created_at datetime(3), updated_at datetime(3))`).Error; err != nil {
			return err
		}
		if err := conn.Exec("SET @inkspace_article_seo_offline = NULL").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "schema.sql"); err == nil {
			t.Fatal("SQL accepted missing offline opt-in")
		}
		var preparedColumns int64
		if err := conn.Raw(`SELECT COUNT(*) FROM information_schema.columns
			WHERE table_schema=DATABASE() AND table_name='articles' AND column_name IN ('slug','publish_at')`).Scan(&preparedColumns).Error; err != nil {
			return err
		}
		if preparedColumns != 0 {
			t.Fatal("unacknowledged prepare changed schema")
		}
		if err := conn.Exec("SET @inkspace_article_seo_offline = 1").Error; err != nil {
			return err
		}
		for range 2 {
			if err := runArticleSQL(t, conn, "schema.sql"); err != nil {
				return err
			}
		}
		if err := conn.Exec("INSERT INTO articles VALUES (1,'Existing',1,NULL,'2020-01-01','2021-01-01','reserved',NULL)").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
			t.Fatal("constraints installed before date repair")
		}
		if conn.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") || conn.Migrator().HasIndex("articles", "ux_articles_slug") {
			t.Fatal("pre-repair constraints installed")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	// The Go command owns its advisory-lock connection, like a separate operator
	// process between SQL phases. Never nest it inside a pinned Connection callback.
	stats, err := MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.Filled != 0 || stats.PublishFilled != 1 {
		t.Fatalf("repair: %+v %v", stats, err)
	}
	if err := db.Connection(func(conn *gorm.DB) error {
		conn = conn.Session(&gorm.Session{NewDB: true})
		if err := conn.Exec("SET @inkspace_article_seo_offline = NULL").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
			t.Fatal("finalize accepted missing offline opt-in")
		}
		// The successful Go repair has already installed both valid indexes.
		if !conn.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") || !conn.Migrator().HasIndex("articles", "ux_articles_slug") {
			t.Fatal("unacknowledged finalize changed existing constraints")
		}
		// Each SQL phase must explicitly opt in on its own pinned connection.
		if err := conn.Exec("SET @inkspace_article_seo_offline = 1").Error; err != nil {
			return err
		}
		for range 2 {
			if err := runArticleSQL(t, conn, "finalize.sql"); err != nil {
				return err
			}
		}
		var slug string
		if err := conn.Raw("SELECT slug FROM articles WHERE id=1").Scan(&slug).Error; err != nil {
			return err
		}
		if slug != "reserved" {
			t.Fatal("populated slug rewritten", slug)
		}
		if err := conn.Exec("DROP INDEX ux_articles_slug ON articles").Error; err != nil {
			return err
		}
		if err := conn.Exec("DROP INDEX idx_articles_first_publication_v1 ON articles").Error; err != nil {
			return err
		}
		for _, invalid := range []string{"", "Upper", "123", "../private", "reserved"} {
			if err := conn.Exec("INSERT INTO articles (id,title,slug,publish_at) VALUES (2,'Invalid',?,'2020-01-01')", invalid).Error; err != nil {
				return err
			}
			if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
				t.Fatalf("invalid/duplicate slug %q accepted", invalid)
			}
			if conn.Migrator().HasIndex("articles", "ux_articles_slug") || conn.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") {
				t.Fatal("invalid data installed constraints")
			}
			// Only this deliberately invalid fixture row is removed, never a real article.
			if err := conn.Exec("DELETE FROM articles WHERE id=2").Error; err != nil {
				return err
			}
		}
		if err := conn.Exec("SET SESSION sql_mode='ALLOW_INVALID_DATES'").Error; err != nil {
			return err
		}
		for _, invalid := range []string{"0000-00-00", "2020-00-01", "2020-02-31"} {
			if err := conn.Exec("UPDATE articles SET publish_at=? WHERE id=1", invalid).Error; err != nil {
				return err
			}
			if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
				t.Fatalf("invalid date %q accepted", invalid)
			}
			if conn.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") {
				t.Fatal("invalid date installed completion marker")
			}
		}
		if err := conn.Exec("UPDATE articles SET publish_at='2020-01-01' WHERE id=1").Error; err != nil {
			return err
		}
		if err := conn.Exec("CREATE INDEX ux_articles_slug ON articles(slug)").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
			t.Fatal("nonunique slug constraint accepted")
		}
		if conn.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") {
			t.Fatal("invalid unique constraint installed completion marker")
		}
		if err := conn.Exec("DROP INDEX ux_articles_slug ON articles").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "finalize.sql"); err != nil {
			return err
		}
		if err := conn.Exec("DROP INDEX idx_articles_first_publication_v1 ON articles").Error; err != nil {
			return err
		}
		if err := conn.Exec("CREATE INDEX idx_articles_first_publication_v1 ON articles(status)").Error; err != nil {
			return err
		}
		if err := runArticleSQL(t, conn, "finalize.sql"); err == nil {
			t.Fatal("invalid completion marker accepted")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	// The Go orchestration must also reject a malformed completion marker.
	if _, err := MigrateArticleSlugs(context.Background(), db); err == nil {
		t.Fatal("Go migration accepted invalid completion marker")
	}
	if err := db.Exec("DROP INDEX idx_articles_first_publication_v1 ON articles").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("DROP INDEX ux_articles_slug ON articles").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE INDEX ux_articles_slug ON articles(slug)").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err == nil {
		t.Fatal("Go migration accepted nonunique slug constraint")
	}
	if db.Migrator().HasIndex("articles", "idx_articles_first_publication_v1") {
		t.Fatal("Go migration installed completion marker before slug validation")
	}
}

func TestArticleSQLArtifactIsGuardedAndPhased(t *testing.T) {
	prepare, err := os.ReadFile("../../scripts/article-slugs/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	finalize, err := os.ReadFile("../../scripts/article-slugs/finalize.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range [][]byte{prepare, finalize} {
		if !strings.Contains(string(data), "@inkspace_article_seo_offline") || !strings.Contains(string(data), "SIGNAL SQLSTATE '45000'") {
			t.Fatal("missing explicit guards")
		}
	}
	if strings.Contains(string(prepare), "CREATE INDEX") || strings.Contains(string(prepare), "CREATE UNIQUE INDEX") {
		t.Fatal("prepare installs readiness constraints")
	}
	if strings.Index(string(finalize), "CREATE UNIQUE INDEX") > strings.Index(string(finalize), "CREATE INDEX idx_articles_first_publication_v1") {
		t.Fatal("completion marker installed before slug constraint")
	}
}
