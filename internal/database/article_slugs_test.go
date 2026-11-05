package database

import (
	"context"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/testutil"
)

func TestArticleSlugMigrationMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.Exec(`CREATE TABLE articles (id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
		title varchar(200), deleted_at datetime(3), status int, updated_at datetime(3), created_at datetime(3) DEFAULT '2019-01-01')`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO articles (id,title,status,deleted_at,updated_at) VALUES
		(1,'Same',1,NULL,'2020-01-01'),(2,'Same',0,NULL,'2020-01-01'),
		(3,'Same',2,NOW(),'2020-01-01'),(4,'😀',1,NULL,'2020-01-01')`).Error; err != nil {
		t.Fatal(err)
	}
	if err := CheckArticleSlugSchema(db); err == nil {
		t.Fatal("unmigrated schema accepted")
	}
	stats, err := MigrateArticleSlugs(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Scanned != 4 || stats.Filled != 4 {
		t.Fatalf("stats: %+v", stats)
	}
	var rows []models.Article
	if err := db.Unscoped().Order("id").Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"same", "same-2", "same-3", "article-4"} {
		if rows[i].Slug != want || rows[i].UpdatedAt.Year() != 2020 {
			t.Fatalf("row: %+v", rows[i])
		}
	}
	stats, err = MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.Filled != 0 {
		t.Fatalf("rerun: %+v %v", stats, err)
	}
	if err := db.Exec("INSERT INTO articles (title,slug) VALUES ('reused','same-3')").Error; err == nil {
		t.Fatal("deleted row's slug reused")
	}
	if err := db.Exec("INSERT INTO articles (id,title,slug) VALUES (10,'Different title','reserved'),(11,'Reserved',NULL)").Error; err != nil {
		t.Fatal(err)
	}
	stats, err = MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.Filled != 1 {
		t.Fatalf("new legacy rows: %+v %v", stats, err)
	}
	var preserved string
	if err := db.Raw("SELECT slug FROM articles WHERE id=10").Scan(&preserved).Error; err != nil {
		t.Fatal(err)
	}
	if preserved != "reserved" {
		t.Fatalf("existing slug overwritten: %q", preserved)
	}
	var slug string
	if err := db.Raw("SELECT slug FROM articles WHERE id=11").Scan(&slug).Error; err != nil {
		t.Fatal(err)
	}
	if slug != "reserved-2" {
		t.Fatalf("got %q", slug)
	}
	if err := db.Exec("UPDATE articles SET slug='invalid/slug' WHERE id=10").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err == nil {
		t.Fatal("invalid existing slug accepted")
	}
}

func TestArticleSlugMigrationConflictsMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.Exec(`CREATE TABLE articles (id bigint unsigned PRIMARY KEY AUTO_INCREMENT, title varchar(200),slug varchar(120),
		created_at datetime(3) DEFAULT '2019-01-01', updated_at datetime(3) DEFAULT '2020-01-01', status int, deleted_at datetime(3))`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO articles(title,slug) VALUES ('first','occupied'),('second','occupied'),('pending',NULL)").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err == nil {
		t.Fatal("duplicate existing slug accepted")
	}
	if db.Migrator().HasIndex(&models.Article{}, "ux_articles_slug") {
		t.Fatal("index created despite conflict")
	}
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM articles WHERE slug IS NULL").Scan(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("failed validation mutated pending rows")
	}
	if err := db.Exec("UPDATE articles SET slug='manual-fix' WHERE id=2").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateArticleSlugs(context.Background(), db); err != nil {
		t.Fatal(err)
	}
}
