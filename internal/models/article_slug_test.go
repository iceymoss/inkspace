package models

import (
	"errors"
	"sync"
	"testing"

	"github.com/iceymoss/inkspace/internal/seo"
	"github.com/iceymoss/inkspace/internal/testutil"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func TestArticleSlugAllocationMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&Article{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX ux_articles_slug ON articles(slug)").Error; err != nil {
		t.Fatal(err)
	}
	deleted := &Article{Title: "Concurrent", Content: "body", AuthorID: 1}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(deleted).Error; err != nil {
			return err
		}
		return AllocateArticleSlug(tx, deleted)
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(deleted).Error; err != nil {
		t.Fatal(err)
	}
	const workers = 12
	results := make(chan error, workers)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			var err error
			for range 20 {
				article := &Article{Title: "Concurrent", Content: "body", AuthorID: 1}
				err = db.Transaction(func(tx *gorm.DB) error {
					if err := tx.Create(article).Error; err != nil {
						return err
					}
					return AllocateArticleSlug(tx, article)
				})
				var mysqlErr *mysqlDriver.MySQLError
				if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1213 || mysqlErr.Number == 1205) {
					continue
				}
				break
			}
			results <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	var articles []Article
	if err := db.Unscoped().Order("slug").Find(&articles).Error; err != nil {
		t.Fatal(err)
	}
	if len(articles) != workers+1 {
		t.Fatalf("got %d rows", len(articles))
	}
	seen := make(map[string]bool)
	for _, article := range articles {
		if !seo.ValidSlug(article.Slug) || seen[article.Slug] {
			t.Fatalf("invalid/duplicate slug %q", article.Slug)
		}
		seen[article.Slug] = true
	}
	for attempt := uint(1); attempt <= workers+1; attempt++ {
		if !seen[seo.SlugCandidate("concurrent", attempt)] {
			t.Fatalf("missing deterministic suffix %d", attempt)
		}
	}
	var article Article
	if err := db.First(&article).Error; err != nil {
		t.Fatal(err)
	}
	oldSlug := article.Slug
	article.Title = "Changed title"
	article.Slug = "attempted-change"
	if err := db.Save(&article).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&article, article.ID).Error; err != nil {
		t.Fatal(err)
	}
	if article.Slug != oldSlug {
		t.Fatal("slug changed through ordinary save")
	}
	if err := db.Model(&article).Updates(map[string]interface{}{"slug": "also-rejected", "title": "Again"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&article, article.ID).Error; err != nil {
		t.Fatal(err)
	}
	if article.Slug != oldSlug {
		t.Fatal("slug changed through update")
	}
}
