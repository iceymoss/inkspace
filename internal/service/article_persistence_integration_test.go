package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func articlePersistenceTestServices(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INKSPACE_READING_TIME_TEST_MYSQL_DSN")
	redisAddr := os.Getenv("INKSPACE_READING_TIME_TEST_REDIS_ADDR")
	if dsn == "" || redisAddr == "" {
		t.Skip("set INKSPACE_READING_TIME_TEST_MYSQL_DSN and INKSPACE_READING_TIME_TEST_REDIS_ADDR for isolated integration tests")
	}
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || parsed.DBName != "inkspace_reading_time_test" {
		t.Fatalf("reading-time integration tests require exactly the inkspace_reading_time_test database: %v", err)
	}
	if os.Getenv("INKSPACE_READING_TIME_TEST_ALLOW_DROP") != "1" {
		t.Fatal("set INKSPACE_READING_TIME_TEST_ALLOW_DROP=1 to acknowledge dedicated test services")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated MySQL: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Article{}, &models.Workspace{}, &models.Doc{}, &models.WorkspaceMember{}); err != nil {
		t.Fatalf("migrate isolated MySQL: %v", err)
	}
	client := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 2})
	if err := client.Ping(database.Ctx).Err(); err != nil {
		t.Fatalf("connect to isolated Redis: %v", err)
	}
	if keys, err := client.DBSize(database.Ctx).Result(); err != nil || keys != 0 {
		t.Fatalf("dedicated Redis DB 2 must start empty: keys=%d, error=%v", keys, err)
	}
	previousDB, previousRDB, previousConfig := database.DB, database.RDB, config.AppConfig
	database.DB, database.RDB = db, client
	config.AppConfig = &config.Config{Cache: config.CacheConfig{ArticleExpire: 60}}
	t.Cleanup(func() {
		if err := client.FlushDB(database.Ctx).Err(); err != nil {
			t.Errorf("clear dedicated Redis database: %v", err)
		}
		database.DB, database.RDB, config.AppConfig = previousDB, previousRDB, previousConfig
		if err := client.Close(); err != nil {
			t.Errorf("close Redis: %v", err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			t.Errorf("retrieve MySQL connection: %v", err)
		} else if err := sqlDB.Close(); err != nil {
			t.Errorf("close MySQL: %v", err)
		}
	})
	return db
}

func TestArticleReadingTimeSavedWithContent(t *testing.T) {
	db := articlePersistenceTestServices(t)
	user := models.User{Username: "reading-author-" + strconv.FormatInt(time.Now().UnixNano(), 10), Email: fmt.Sprintf("reading-%d@author.test", time.Now().UnixNano()), Password: "hash"}
	admin := models.User{Username: "reading-admin-" + strconv.FormatInt(time.Now().UnixNano(), 10), Email: fmt.Sprintf("reading-%d@admin.test", time.Now().UnixNano()), Password: "hash", Role: "admin"}
	category := models.Category{Name: "Reading-" + strconv.FormatInt(time.Now().UnixNano(), 10), Slug: "reading-" + strconv.FormatInt(time.Now().UnixNano(), 10), Logo: "reading.png"}
	for _, record := range []interface{}{&user, &admin, &category} {
		if err := db.Create(record).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
	t.Cleanup(func() {
		for _, table := range []string{"article_tags", "articles", "categories", "users"} {
			var err error
			switch table {
			case "article_tags":
				err = db.Exec("DELETE FROM article_tags WHERE article_id IN (SELECT id FROM articles WHERE author_id IN (?, ?))", user.ID, admin.ID).Error
			case "articles":
				err = db.Unscoped().Where("author_id IN ?", []uint{user.ID, admin.ID}).Delete(&models.Article{}).Error
			case "categories":
				err = db.Unscoped().Delete(&category).Error
			case "users":
				err = db.Unscoped().Where("id IN ?", []uint{user.ID, admin.ID}).Delete(&models.User{}).Error
			}
			if err != nil {
				t.Errorf("remove article test fixtures from %s: %v", table, err)
			}
		}
	})

	for _, tc := range []struct {
		name   string
		actor  uint
		role   string
		before string
		after  string
		want   int
	}{
		{name: "author updates legacy article", actor: user.ID, role: "user", before: strings.Repeat("文", 301), after: strings.Repeat("文", 601), want: 3},
		{name: "administrator updates legacy article", actor: admin.ID, role: "admin", before: strings.Repeat("文", 301), after: "", want: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			legacy := models.Article{Title: "旧文章", AuthorID: user.ID, CategoryID: category.ID, Content: tc.before, ReadingTime: 0}
			if err := db.Create(&legacy).Error; err != nil {
				t.Fatal(err)
			}
			// Prime the detail cache with the old zero-valued record before editing.
			if _, err := NewArticleService().GetByID(legacy.ID); err != nil {
				t.Fatal(err)
			}
			updated, err := NewArticleService().Update(legacy.ID, &models.ArticleRequest{
				Title: legacy.Title, Content: tc.after, CategoryID: category.ID, Status: 1,
			}, tc.actor, tc.role)
			if err != nil {
				t.Fatal(err)
			}
			if updated.Content != tc.after || updated.ReadingTime != tc.want {
				t.Fatalf("updated content/reading time = %q/%d, want %q/%d", updated.Content, updated.ReadingTime, tc.after, tc.want)
			}
			cacheKey := fmt.Sprintf("article:%d", legacy.ID)
			if exists, err := database.RDB.Exists(database.Ctx, cacheKey).Result(); err != nil || exists != 0 {
				t.Fatalf("stale detail cache after edit: exists=%d, error=%v", exists, err)
			}
			persisted, err := NewArticleService().GetByID(legacy.ID)
			if err != nil {
				t.Fatal(err)
			}
			if persisted.Content != tc.after || persisted.ToResponse().ReadingTime != tc.want {
				t.Fatalf("detail content/reading time = %q/%d, want %q/%d", persisted.Content, persisted.ReadingTime, tc.after, tc.want)
			}
		})
	}

	for _, tc := range []struct {
		name    string
		content string
		want    int
	}{
		{name: "long article", content: strings.Repeat("文", 301), want: 2},
		{name: "empty article", content: "", want: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			created, err := NewArticleService().Create(&models.ArticleRequest{
				Title: tc.name, Content: tc.content, CategoryID: category.ID, Status: 1,
			}, user.ID)
			if err != nil {
				t.Fatal(err)
			}
			persisted, err := NewArticleService().GetByID(created.ID)
			if err != nil {
				t.Fatal(err)
			}
			if created.ReadingTime != tc.want || persisted.Content != tc.content || persisted.ToResponse().ReadingTime != tc.want {
				t.Fatalf("created %d and detail %q/%d, want %q/%d", created.ReadingTime, persisted.Content, persisted.ReadingTime, tc.content, tc.want)
			}
		})
	}
}

func TestPublishKnowledgeDocToBlogSavesArticleReadingTime(t *testing.T) {
	db := articlePersistenceTestServices(t)
	owner := models.User{Username: "knowledge-owner-" + strconv.FormatInt(time.Now().UnixNano(), 10), Email: fmt.Sprintf("knowledge-%d@owner.test", time.Now().UnixNano()), Password: "hash"}
	category := models.Category{Name: "Knowledge-" + strconv.FormatInt(time.Now().UnixNano(), 10), Slug: "knowledge-" + strconv.FormatInt(time.Now().UnixNano(), 10), Logo: "knowledge.png"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}
	workspace := models.Workspace{OwnerID: owner.ID, Name: "Workspace"}
	if err := db.Create(&workspace).Error; err != nil {
		t.Fatal(err)
	}
	original := strings.Repeat("文", 301)
	doc := models.Doc{WorkspaceID: workspace.ID, OwnerID: owner.ID, Title: "知识文章", Content: original, WordCount: 42, Revision: 1}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		var savedDoc models.Doc
		if err := db.First(&savedDoc, doc.ID).Error; err != nil {
			t.Errorf("retrieve document fixture: %v", err)
		} else if savedDoc.ArticleID != nil {
			if err := db.Exec("DELETE FROM article_tags WHERE article_id = ?", *savedDoc.ArticleID).Error; err != nil {
				t.Errorf("remove published article tags: %v", err)
			}
			if err := db.Unscoped().Delete(&models.Article{}, *savedDoc.ArticleID).Error; err != nil {
				t.Errorf("remove published article: %v", err)
			}
		}
		for _, record := range []interface{}{&doc, &workspace, &category, &owner} {
			if err := db.Unscoped().Delete(record).Error; err != nil {
				t.Errorf("remove knowledge test fixture: %v", err)
			}
		}
	})

	publish := func(wantContent string, wantMinutes int) {
		t.Helper()
		article, err := NewDocService().PublishToBlog(doc.ID, owner.ID, &models.DocPublishToBlogRequest{CategoryID: category.ID})
		if err != nil {
			t.Fatal(err)
		}
		persisted, err := NewArticleService().GetByID(article.ID)
		if err != nil {
			t.Fatal(err)
		}
		if persisted.Content != wantContent || persisted.ToResponse().ReadingTime != wantMinutes {
			t.Fatalf("published content/reading time = %q/%d, want %q/%d", persisted.Content, persisted.ReadingTime, wantContent, wantMinutes)
		}
		var savedDoc models.Doc
		if err := db.First(&savedDoc, doc.ID).Error; err != nil {
			t.Fatal(err)
		}
		if savedDoc.ArticleID == nil || *savedDoc.ArticleID != article.ID || savedDoc.WordCount != 42 {
			t.Fatalf("document link or statistics changed: %+v", savedDoc)
		}
	}
	publish(original, 2)
	updated := "新正文"
	if err := db.Model(&doc).Update("content", updated).Error; err != nil {
		t.Fatal(err)
	}
	publish(updated, 1)
}
