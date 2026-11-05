package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/testutil"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestFeedsMySQLHTTP(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Article{}, &models.Setting{}); err != nil {
		t.Fatal(err)
	}
	oldDB, oldConfig := database.DB, config.AppConfig
	database.DB = db
	config.AppConfig = &config.Config{PublicSiteURL: "https://canonical.example", JWT: config.JWTConfig{Secret: "feed-http-test", ExpireHours: 1}}
	t.Cleanup(func() { database.DB, config.AppConfig = oldDB, oldConfig })
	author := models.User{ID: 1, Username: "author", Nickname: "作者 & 中文", Email: "private-author@example.org", Password: "private-password"}
	if err := db.Create(&author).Error; err != nil {
		t.Fatal(err)
	}
	for _, setting := range []models.Setting{
		{Key: models.SettingSiteName, Value: "公开站点", IsPublic: true},
		{Key: models.SettingSiteDescription, Value: "private-site-description", IsPublic: false},
	} {
		if err := db.Create(&setting).Error; err != nil {
			t.Fatal(err)
		}
	}
	published := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	for id := 1; id <= 55; id++ {
		first := published
		if id == 1 {
			first = first.Add(time.Hour)
		}
		article := models.Article{
			ID: uint(id), Slug: fmt.Sprintf("post-%d", id), Title: fmt.Sprintf("Title %d", id),
			Content: "**正文** & 中文", AuthorID: 1, Status: 1, PublishAt: &first,
			CreatedAt: published.Add(-24 * time.Hour), UpdatedAt: published.Add(2 * time.Hour),
		}
		if err := db.Create(&article).Error; err != nil {
			t.Fatal(err)
		}
	}
	for id, slug := range []string{"private", "draft", "deleted"} {
		article := models.Article{ID: uint(100 + id), Slug: slug, Title: "SECRET " + slug, Content: "secret body", Status: 1, AuthorID: 1, PublishAt: &published}
		if slug == "private" {
			article.Status = 2
		}
		if slug == "deleted" {
			article.DeletedAt = gorm.DeletedAt{Time: published, Valid: true}
		}
		if err := db.Create(&article).Error; err != nil {
			t.Fatal(err)
		}
		if slug == "draft" {
			if err := db.Model(&models.Article{}).Where(models.ArticleColumns.ID+" = ?", article.ID).UpdateColumn(models.ArticleColumns.Status, 0).Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	gin.SetMode(gin.TestMode)
	r := SetupUserRouter(fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}})
	token, err := utils.GenerateToken(1, "author", "admin")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/feed.xml", "/atom.xml"} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path+"?status=0&author_id=1&show_all=true&page_size=1000&sort_order=asc", nil)
		req.Host = "attacker.invalid"
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
		entries := decodedFeedEntries(t, w.Body.Bytes())
		if len(entries) != 50 || entries[0].URL != "https://canonical.example/blog/post-1" || entries[1].URL != "https://canonical.example/blog/post-55" || entries[49].URL != "https://canonical.example/blog/post-7" || entries[0].Author != author.Nickname {
			t.Fatalf("limit/order/public author: %+v", entries)
		}
		for _, secret := range []string{"SECRET", "private-author", "private-password", "private-site-description", "attacker.invalid"} {
			if strings.Contains(w.Body.String(), secret) {
				t.Fatalf("leaked %s", secret)
			}
		}
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			conditional := httptest.NewRecorder()
			req := httptest.NewRequest(method, path, nil)
			req.Header.Set("If-None-Match", "W/"+w.Header().Get("ETag"))
			r.ServeHTTP(conditional, req)
			if conditional.Code != 304 || conditional.Body.Len() != 0 {
				t.Fatal("conditional", conditional.Code, conditional.Body.String())
			}
		}
		head := httptest.NewRecorder()
		r.ServeHTTP(head, httptest.NewRequest(http.MethodHead, path, nil))
		if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("ETag") != w.Header().Get("ETag") || head.Header().Get("Content-Length") != w.Header().Get("Content-Length") {
			t.Fatal("HEAD", head.Code, head.Header())
		}
	}
	// A status change immediately changes even a conditional representation.
	before := httptest.NewRecorder()
	r.ServeHTTP(before, httptest.NewRequest(http.MethodGet, "/feed.xml", nil))
	if err := db.Model(&models.Article{}).Where(models.ArticleColumns.ID+" = ?", 1).UpdateColumn(models.ArticleColumns.Status, 2).Error; err != nil {
		t.Fatal(err)
	}
	after := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/feed.xml", nil)
	req.Header.Set("If-None-Match", before.Header().Get("ETag"))
	r.ServeHTTP(after, req)
	if after.Code != 200 || strings.Contains(after.Body.String(), "https://canonical.example/blog/post-1<") {
		t.Fatal("unpublished article retained", after.Code, after.Body.String())
	}
	// Force an isolated database query failure; never touch a developer schema.
	if err := db.Migrator().DropTable(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/feed.xml", "/atom.xml"} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, path, nil)
			req.Header.Set("If-None-Match", "*")
			r.ServeHTTP(w, req)
			if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("ETag") != "" || w.Body.Len() != 0 {
				t.Fatalf("DB failure: %d %v %s", w.Code, w.Header(), w.Body.String())
			}
		}
	}
}
