package database

import (
	"context"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/testutil"

	"gorm.io/gorm"
)

func TestArticlePublicationBackfillMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.Exec(`CREATE TABLE articles (id bigint unsigned PRIMARY KEY, title varchar(200),
		slug varchar(120), status int, deleted_at datetime(3), created_at datetime(3),
		updated_at datetime(3) ON UPDATE CURRENT_TIMESTAMP(3), publish_at datetime(3))`).Error; err != nil {
		t.Fatal(err)
	}
	// Store MySQL legacy zero/partial-zero/impossible dates on one connection.
	// This deliberately does not depend on the test server's global SQL mode.
	if err := db.Connection(func(conn *gorm.DB) error {
		if err := conn.Exec("SET SESSION sql_mode = 'ALLOW_INVALID_DATES'").Error; err != nil {
			return err
		}
		return conn.Exec(`INSERT INTO articles VALUES
			(1,'Missing',NULL,1,NULL,'2020-01-02','2021-03-04',NULL),
			(2,'Zero','zero',0,NULL,'2020-01-02','2021-03-04','0000-00-00'),
			(3,'Partial','partial',2,NULL,'2020-01-02','2021-03-04','2020-00-01'),
			(4,'Impossible','impossible',1,'2022-01-01','2020-01-02','2021-03-04','2020-02-31'),
			(5,'Valid','valid',1,NULL,'2020-01-02','2021-03-04','2019-01-01'),
			(6,'Future','future',1,NULL,'2020-01-02','2021-03-04','2099-01-01')`).Error
	}); err != nil {
		t.Fatal(err)
	}
	stats, err := MigrateArticleSlugs(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	if stats.PublishScanned != 6 || stats.PublishFilled != 4 || stats.PublishPreserved != 2 {
		t.Fatalf("publication stats: %+v", stats)
	}
	var rows []struct {
		ID          uint
		Publication string
		Modified    string
	}
	if err := db.Raw("SELECT id, CAST(publish_at AS CHAR) AS publication, CAST(updated_at AS CHAR) AS modified FROM articles ORDER BY id").Scan(&rows).Error; err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		want := "2020-01-02 00:00:00.000"
		if i == 4 {
			want = "2019-01-01 00:00:00.000"
		}
		if i == 5 {
			want = "2099-01-01 00:00:00.000"
		}
		if row.Publication != want || row.Modified != "2021-03-04 00:00:00.000" {
			t.Fatalf("row %+v, want publication %s and unchanged content timestamp", row, want)
		}
	}
	stats, err = MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.PublishFilled != 0 || stats.PublishPreserved != 6 {
		t.Fatalf("rerun: %+v %v", stats, err)
	}
}

func TestArticlePublicationReadinessMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	if err := CheckArticlePublicationSchema(db); err == nil {
		t.Fatal("startup accepted schema before explicit publication migration")
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	if err := CheckArticlePublicationSchema(db); err != nil {
		t.Fatal(err)
	}
}

func TestArticlePublicationBackfillRejectsInvalidCreationMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.Exec(`CREATE TABLE articles (id bigint unsigned PRIMARY KEY, title varchar(200), slug varchar(120),
		created_at datetime(3), publish_at datetime(3), updated_at datetime(3) DEFAULT '2020-01-01', status int, deleted_at datetime(3))`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO articles (id,title,slug,created_at,publish_at) VALUES (1,'Repairable','repairable','2020-01-01',NULL),(2,'Invalid','invalid',NULL,NULL)").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err == nil {
		t.Fatal("invalid fallback creation timestamp accepted")
	}
	var missing int64
	if err := db.Raw("SELECT COUNT(*) FROM articles WHERE publish_at IS NULL").Scan(&missing).Error; err != nil {
		t.Fatal(err)
	}
	if missing != 2 {
		t.Fatal("publication backfill wrote data before validating all fallbacks")
	}
	if err := db.Exec("UPDATE articles SET created_at='2020-02-01' WHERE id=2").Error; err != nil {
		t.Fatal(err)
	}
	stats, err := MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.PublishFilled != 2 {
		t.Fatalf("manual correction rerun: %+v %v", stats, err)
	}
}
