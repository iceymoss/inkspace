package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/testutil"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// publicationMySQL uses real, randomly isolated MySQL and a disposable RESP
// cache invalidation stub. Persistence is never mocked and no Redis DB is used.
func publicationMySQL(t *testing.T) (*gorm.DB, *models.User) {
	t.Helper()
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Article{}); err != nil {
		t.Fatal(err)
	}
	if _, err := database.MigrateArticleSlugs(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	author := &models.User{Username: "publisher", Email: "publisher@example.test", Password: "hash", Nickname: "Public author"}
	if err := db.Create(author).Error; err != nil {
		t.Fatal(err)
	}
	oldDB, oldRedis := database.DB, database.RDB
	client := publicationCacheStub(t)
	database.DB, database.RDB = db, client
	t.Cleanup(func() {
		database.DB, database.RDB = oldDB, oldRedis
		if err := client.Close(); err != nil {
			t.Error(err)
		}
	})
	return db, author
}

func publicationCacheStub(t *testing.T) *redis.Client {
	t.Helper()
	var workers sync.WaitGroup
	t.Cleanup(workers.Wait)
	return redis.NewClient(&redis.Options{Dialer: func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		workers.Add(1)
		go func() {
			defer workers.Done()
			defer server.Close()
			r := bufio.NewReader(server)
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					return
				}
				count, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "*")))
				if err != nil {
					return
				}
				args := make([]string, count)
				for i := range args {
					line, err = r.ReadString('\n')
					if err != nil {
						return
					}
					size, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "$")))
					if err != nil {
						return
					}
					data := make([]byte, size+2)
					if _, err := io.ReadFull(r, data); err != nil {
						return
					}
					args[i] = string(data[:size])
				}
				response := "-ERR test cache only supports invalidation\r\n"
				switch strings.ToLower(args[0]) {
				case "scan":
					response = "*2\r\n$1\r\n0\r\n*0\r\n"
				case "del":
					response = ":0\r\n"
				}
				if _, err := io.WriteString(server, response); err != nil {
					return
				}
			}
		}()
		return client, nil
	}})
}

func TestArticlePublicationLifecycleMySQL(t *testing.T) {
	_, author := publicationMySQL(t)
	s := NewArticleService()
	for _, initialStatus := range []int{0, 1, 2} {
		t.Run(fmt.Sprintf("initial status %d", initialStatus), func(t *testing.T) {
			req := &models.ArticleRequest{Title: "Publication", Content: "Original body", Status: initialStatus}
			before := time.Now().Add(-time.Second)
			article, err := s.Create(req, author.ID)
			if err != nil {
				t.Fatal(err)
			}
			if article.Status != initialStatus || (article.PublishAt != nil) != (initialStatus == 1) {
				t.Fatalf("created article: %+v", article)
			}
			if initialStatus != 1 {
				req.Status = 1
				article, err = s.Update(article.ID, req, author.ID, "user")
				if err != nil {
					t.Fatal(err)
				}
			}
			if article.PublishAt == nil || article.PublishAt.Before(before) || article.PublishAt.After(time.Now()) {
				t.Fatalf("first publication time: %v", article.PublishAt)
			}
			// Observe persisted millisecond precision before comparing returned times.
			article, err = s.GetPublicByID(context.Background(), article.ID)
			if err != nil {
				t.Fatal(err)
			}
			first, previousUpdate := *article.PublishAt, article.UpdatedAt
			for _, status := range []int{0, 2, 1, 1} {
				time.Sleep(5 * time.Millisecond)
				req.Status, req.Content = status, fmt.Sprintf("Edited body %d", status)
				article, err = s.Update(article.ID, req, author.ID, "user")
				if err != nil {
					t.Fatal(err)
				}
				if article.PublishAt == nil || !article.PublishAt.Equal(first) {
					t.Fatalf("status %d reset first publication: %v, want %v", status, article.PublishAt, first)
				}
				if !article.UpdatedAt.After(previousUpdate) {
					t.Fatalf("content/state change did not advance updated_at: %v", article.UpdatedAt)
				}
				previousUpdate = article.UpdatedAt
			}
		})
	}
}

func TestConcurrentFirstPublicationMySQL(t *testing.T) {
	_, author := publicationMySQL(t)
	s := NewArticleService()
	req := &models.ArticleRequest{Title: "Concurrent publication", Content: "Body", Status: 0}
	article, err := s.Create(req, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	req.Status = 1
	const n = 6
	start := make(chan struct{})
	results := make(chan *models.Article, n)
	errors := make(chan error, n)
	var workers sync.WaitGroup
	for range n {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			updated, err := s.Update(article.ID, req, author.ID, "user")
			results <- updated
			errors <- err
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	persisted, err := s.GetPublicByID(context.Background(), article.ID)
	if err != nil || persisted.PublishAt == nil {
		t.Fatalf("first publication: %+v %v", persisted, err)
	}
	for updated := range results {
		if updated.PublishAt == nil || !updated.PublishAt.Equal(*persisted.PublishAt) {
			t.Fatalf("concurrent publisher changed first publication: %v, want %v", updated.PublishAt, persisted.PublishAt)
		}
	}
}

func TestHistoricalPublicationSurvivesRepublishMySQL(t *testing.T) {
	db, author := publicationMySQL(t)
	s := NewArticleService()
	article, err := s.Create(&models.ArticleRequest{Title: "Historical private post", Content: "Body", Status: 2}, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	created := time.Date(2018, 2, 3, 4, 5, 6, 0, time.UTC)
	if err := db.Model(article).UpdateColumn("created_at", created).Error; err != nil {
		t.Fatal(err)
	}
	stats, err := database.MigrateArticleSlugs(context.Background(), db)
	if err != nil || stats.PublishFilled != 1 {
		t.Fatalf("historical repair: %+v %v", stats, err)
	}
	updated, err := s.Update(article.ID, &models.ArticleRequest{Title: article.Title, Content: "Edited", Status: 1}, author.ID, "user")
	if err != nil {
		t.Fatal(err)
	}
	if updated.PublishAt == nil || !updated.PublishAt.Equal(created) {
		t.Fatalf("republishing reset repaired history: %v, want %v", updated.PublishAt, created)
	}
}

func TestArticleInteractionsPreserveContentTimeMySQL(t *testing.T) {
	testArticleInteractionsPreserveContentTime(t, false)
}

func TestArticleInteractionsPreserveLegacyAutomaticTimeMySQL(t *testing.T) {
	testArticleInteractionsPreserveContentTime(t, true)
}

func testArticleInteractionsPreserveContentTime(t *testing.T, legacy bool) {
	t.Helper()
	db, author := publicationMySQL(t)
	// Finish all association-driven schema migration BEFORE installing the legacy
	// column, so AutoMigrate cannot silently erase the behavior under test.
	if err := db.AutoMigrate(&models.ArticleFavorite{}, &models.Like{}, &models.Comment{}, &models.Setting{}, &models.Work{}, &models.UserFollow{}); err != nil {
		t.Fatal(err)
	}
	if legacy {
		// This DDL is ONLY on SEOTestMySQL's random, opt-in test schema.
		if err := db.Exec("ALTER TABLE articles MODIFY updated_at datetime(3) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3)").Error; err != nil {
			t.Fatal(err)
		}
		var extra string
		if err := db.Raw(`SELECT extra FROM information_schema.columns WHERE table_schema=DATABASE()
			AND table_name='articles' AND column_name='updated_at'`).Scan(&extra).Error; err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.ToLower(extra), "on update current_timestamp(3)") {
			t.Fatalf("legacy automatic timestamp fixture missing: %q", extra)
		}
	}
	s := NewArticleService()
	article, err := s.Create(&models.ArticleRequest{Title: "Interaction timestamps", Content: "Body", Status: 1}, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	// Use a historical content time to detect any automatic timestamp mutation.
	modified := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := db.Model(article).UpdateColumn("updated_at", modified).Error; err != nil {
		t.Fatal(err)
	}
	article, err = s.GetPublicByID(context.Background(), article.ID)
	if err != nil {
		t.Fatal(err)
	}
	first := *article.PublishAt
	check := func(t *testing.T, views, likes, favorites, comments int) {
		t.Helper()
		got, err := s.GetPublicByID(context.Background(), article.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ViewCount != views || got.LikeCount != likes || got.FavoriteCount != favorites || got.CommentCount != comments {
			t.Fatalf("interaction counters: %+v", got)
		}
		if !got.UpdatedAt.Equal(modified) || got.PublishAt == nil || !got.PublishAt.Equal(first) {
			t.Fatalf("interaction changed publication/content time: %+v", got)
		}
		entries, err := s.GetPublicSitemapArticles(context.Background())
		if err != nil || len(entries) != 1 || !entries[0].LastModified().Equal(modified) {
			t.Fatalf("interaction changed sitemap lastmod: %+v %v", entries, err)
		}
	}
	if err := s.RecordPublicArticleView(context.Background(), article.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.IncrementViewCount(article.ID); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 0)
	if err := NewFavoriteService().AddFavorite(author.ID, article.ID); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 1, 0)
	if err := NewFavoriteService().RemoveFavorite(author.ID, article.ID); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 0)
	if err := NewLikeService().LikeArticle(author.ID, article.ID); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 1, 0, 0)
	if err := NewLikeService().LikeArticle(author.ID, article.ID); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 0)
	comment, err := NewCommentService().Create(&models.CommentRequest{ArticleID: &article.ID, Content: "A comment"}, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 1)
	if err := NewCommentService().Delete(comment.ID, author.ID, "user"); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 0)
	// Moderation has independent counter writes in both directions.
	comment, err = NewCommentService().Create(&models.CommentRequest{ArticleID: &article.ID, Content: "Moderated comment"}, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 1)
	for _, status := range []int{0, 1, -1} {
		if err := NewCommentService().UpdateStatus(comment.ID, status); err != nil {
			t.Fatal(err)
		}
		want := 0
		if status == 1 {
			want = 1
		}
		check(t, 2, 0, 0, want)
	}
	// Legacy database counter-maintenance entry points also must not advance
	// editorial time. No user/reply counter is involved in this detached value.
	detached := &models.Comment{ArticleID: &article.ID}
	if err := database.AfterCreateComment(db, detached); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 1)
	if err := database.AfterDeleteComment(db, detached); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 0)
	// Force the maintenance sync to change both counters, not a no-op UPDATE.
	if err := db.Model(article).UpdateColumns(map[string]any{"comment_count": 8, "favorite_count": 8, "updated_at": modified}).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.SyncCounters(); err != nil {
		t.Fatal(err)
	}
	check(t, 2, 0, 0, 1) // rejected, undeleted comment is included by legacy sync.
}
