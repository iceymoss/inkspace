package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/iceymoss/inkspace/internal/models"
)

func TestPublicArticleCollectionsMySQL(t *testing.T) {
	db, author := publicationMySQL(t)
	published := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	modified := published.Add(time.Hour)
	for id := uint(1); id <= 55; id++ {
		article := models.Article{
			ID: id, Slug: fmt.Sprintf("post-%d", id), Title: "Public post", Content: "Feed fallback body",
			AuthorID: author.ID, Status: 1, PublishAt: &published, CreatedAt: published, UpdatedAt: modified,
			IsTop: id == 1,
		}
		if err := db.Create(&article).Error; err != nil {
			t.Fatal(err)
		}
	}
	// IDs 53–55 must never appear, even though the reader is also the author.
	if err := db.Model(&models.Article{}).Where("id = 53").UpdateColumn("status", 0).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&models.Article{}).Where("id = 54").UpdateColumn("status", 2).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&models.Article{}, 55).Error; err != nil {
		t.Fatal(err)
	}
	// First-publication order, not creation/update/pinning, is the feed order.
	later := published.Add(2 * time.Hour)
	if err := db.Model(&models.Article{}).Where("id = 1").UpdateColumn("publish_at", later).Error; err != nil {
		t.Fatal(err)
	}
	s := NewArticleService()
	feed, err := s.GetRecentPublicArticles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) != 50 || feed[0].ID != 1 || feed[1].ID != 52 || feed[49].ID != 4 {
		t.Fatalf("unexpected feed order/limit: %+v", feed)
	}
	if feed[0].Author == nil || feed[0].Author.Nickname != "Public author" || feed[0].Content != "Feed fallback body" {
		t.Fatalf("feed lacks author or fallback content: %+v", feed[0])
	}
	if feed[0].Author.Email != "" || feed[0].Author.Password != "" {
		t.Fatal("feed projection loaded private author data")
	}
	entries, err := s.GetPublicSitemapArticles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 52 {
		t.Fatalf("sitemap count=%d, want 52", len(entries))
	}
	for i, entry := range entries {
		if entry.ID != uint(i+1) || entry.Slug != fmt.Sprintf("post-%d", i+1) || !entry.LastModified().Equal(modified) {
			t.Fatalf("sitemap entry %d: %+v", i, entry)
		}
	}
	// No collection cache may retain an article after unpublication.
	if err := db.Model(&models.Article{}).Where("id = 1").UpdateColumn("status", 2).Error; err != nil {
		t.Fatal(err)
	}
	feed, err = s.GetRecentPublicArticles(context.Background())
	if err != nil || len(feed) != 50 || feed[0].ID != 52 {
		t.Fatalf("feed visibility stale: %v %v", feed, err)
	}
	entries, err = s.GetPublicSitemapArticles(context.Background())
	if err != nil || len(entries) != 51 || entries[0].ID != 2 {
		t.Fatalf("sitemap visibility stale: %v %v", entries, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.GetRecentPublicArticles(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("feed ignored context cancellation: %v", err)
	}
	if _, err := s.GetPublicSitemapArticles(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("sitemap ignored context cancellation: %v", err)
	}
	if err := db.Migrator().DropTable(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetRecentPublicArticles(context.Background()); err == nil {
		t.Fatal("feed query failure became empty success")
	}
	if _, err := s.GetPublicSitemapArticles(context.Background()); err == nil {
		t.Fatal("sitemap query failure became empty success")
	}
}

func TestPublicSitemapLastModified(t *testing.T) {
	created := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	published, updated := created.Add(time.Hour), created.Add(2*time.Hour)
	for _, tt := range []struct {
		name  string
		entry PublicSitemapArticle
		want  time.Time
	}{
		{name: "content time", entry: PublicSitemapArticle{CreatedAt: created, PublishAt: &published, UpdatedAt: updated}, want: updated},
		{name: "first publication fallback", entry: PublicSitemapArticle{CreatedAt: created, PublishAt: &published}, want: published},
		{name: "creation fallback", entry: PublicSitemapArticle{CreatedAt: created}, want: created},
		{name: "zero publication fallback", entry: PublicSitemapArticle{CreatedAt: created, PublishAt: new(time.Time)}, want: created},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.LastModified(); !got.Equal(tt.want) {
				t.Fatalf("lastmod=%v, want %v", got, tt.want)
			}
		})
	}
}
