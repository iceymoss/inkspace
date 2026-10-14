package handler

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestArticleDetailReadingTimeWithMySQLAndRedis(t *testing.T) {
	dsn := os.Getenv("INKSPACE_READING_TIME_DETAIL_MYSQL_DSN")
	redisAddr := os.Getenv("INKSPACE_READING_TIME_TEST_REDIS_ADDR")
	if dsn == "" || redisAddr == "" {
		t.Skip("set INKSPACE_READING_TIME_DETAIL_MYSQL_DSN and INKSPACE_READING_TIME_TEST_REDIS_ADDR for isolated integration tests")
	}
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || parsed.DBName != "inkspace_reading_time_detail_test" {
		t.Fatalf("detail integration tests require the inkspace_reading_time_detail_test database: %v", err)
	}
	if os.Getenv("INKSPACE_READING_TIME_TEST_ALLOW_DROP") != "1" {
		t.Fatal("set INKSPACE_READING_TIME_TEST_ALLOW_DROP=1 to acknowledge dedicated test services")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Article{}); err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 1})
	if err := client.Ping(database.Ctx).Err(); err != nil {
		t.Fatal(err)
	}
	if keys, err := client.DBSize(database.Ctx).Result(); err != nil || keys != 0 {
		t.Fatalf("dedicated Redis DB 1 must start empty: keys=%d, error=%v", keys, err)
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
			t.Errorf("get MySQL connection: %v", err)
		} else if err := sqlDB.Close(); err != nil {
			t.Errorf("close MySQL: %v", err)
		}
	})

	user := models.User{Username: "detail-author-" + strconv.FormatInt(time.Now().UnixNano(), 10), Email: fmt.Sprintf("detail-%d@reading.test", time.Now().UnixNano()), Password: "hash"}
	category := models.Category{Name: "Detail-" + strconv.FormatInt(time.Now().UnixNano(), 10), Slug: "detail-reading-" + strconv.FormatInt(time.Now().UnixNano(), 10), Logo: "reading.png"}
	legacy := models.Article{}
	saved := models.Article{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		for _, article := range []*models.Article{&legacy, &saved} {
			if article.ID == 0 {
				continue
			}
			if err := db.Exec("DELETE FROM article_tags WHERE article_id = ?", article.ID).Error; err != nil {
				t.Errorf("remove detail article tags: %v", err)
			}
			if err := db.Unscoped().Delete(article).Error; err != nil {
				t.Errorf("remove detail article fixture: %v", err)
			}
		}
		for _, record := range []interface{}{&category, &user} {
			if err := db.Unscoped().Delete(record).Error; err != nil {
				t.Errorf("remove detail fixture: %v", err)
			}
		}
	})
	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}
	content := strings.Repeat("文", 301)
	legacy = models.Article{Title: "历史文章", AuthorID: user.ID, CategoryID: category.ID, Content: content, ReadingTime: 0, Status: 1}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}
	key := fmt.Sprintf("article:%d", legacy.ID)
	if err := client.Del(database.Ctx, key).Err(); err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/articles/:id", NewArticleHandler().GetDetail)
	getDetail := func(id uint) int {
		t.Helper()
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest("GET", fmt.Sprintf("/articles/%d", id), nil))
		var result struct {
			Code int `json:"code"`
			Data struct {
				ReadingTime int `json:"reading_time"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatalf("decode detail: %v", err)
		}
		if recorder.Code != 200 || result.Code != 0 {
			t.Fatalf("detail status %d, body %s", recorder.Code, recorder.Body.String())
		}
		return result.Data.ReadingTime
	}
	assertOriginal := func(wantCache []byte) {
		t.Helper()
		var original models.Article
		if err := db.First(&original, legacy.ID).Error; err != nil {
			t.Fatal(err)
		}
		if original.ReadingTime != 0 {
			t.Fatalf("legacy MySQL reading_time = %d, want 0", original.ReadingTime)
		}
		cached, err := client.Get(database.Ctx, key).Bytes()
		if err != nil {
			t.Fatal(err)
		}
		if !json.Valid(cached) || !strings.Contains(string(cached), `"reading_time":0`) {
			t.Fatalf("cache does not preserve the original zero: %s", cached)
		}
		if wantCache != nil && string(cached) != string(wantCache) {
			t.Fatalf("cache changed on detail hit: before %s, after %s", wantCache, cached)
		}
	}

	if got := getDetail(legacy.ID); got != 2 {
		t.Fatalf("uncached legacy detail reading_time = %d, want 2", got)
	}
	assertOriginal(nil)
	originalCache, err := client.Get(database.Ctx, key).Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if got := getDetail(legacy.ID); got != 2 {
		t.Fatalf("cached legacy detail reading_time = %d, want 2", got)
	}
	assertOriginal(originalCache)

	saved = models.Article{Title: "已有时长", AuthorID: user.ID, CategoryID: category.ID, Content: content, ReadingTime: 5, Status: 1}
	if err := db.Create(&saved).Error; err != nil {
		t.Fatal(err)
	}
	if got := getDetail(saved.ID); got != 5 {
		t.Fatalf("saved reading_time = %d, want 5", got)
	}
	// Detail views increment view_count asynchronously; wait before closing the test database.
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		var viewed []models.Article
		if err := db.Find(&viewed, []uint{legacy.ID, saved.ID}).Error; err != nil {
			t.Fatal(err)
		}
		if len(viewed) == 2 && viewed[0].ViewCount+viewed[1].ViewCount == 3 {
			return
		}
	}
	t.Fatal("timed out waiting for asynchronous view increments")
}
