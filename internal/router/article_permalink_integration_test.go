package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/testutil"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestPublicPermalinksMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Article{}, &models.Setting{}); err != nil {
		t.Fatal(err)
	}
	if _, err := database.MigrateArticleSlugs(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	old := database.DB
	oldConfig := config.AppConfig
	config.AppConfig = &config.Config{
		PublicSiteURL: "https://canonical.example",
		JWT:           config.JWTConfig{Secret: "public-contract-test", ExpireHours: 1},
	}
	database.DB = db
	t.Cleanup(func() { database.DB = old; config.AppConfig = oldConfig })
	for _, setting := range []models.Setting{
		{Key: "site_name", Value: "PRIVATE SITE NAME", IsPublic: false},
		{Key: "site_description", Value: "公开描述", IsPublic: true},
		{Key: "public_site_url", Value: "https://attacker.invalid", IsPublic: true},
		{Key: "site_logo", Value: "javascript:alert(1)", IsPublic: true},
	} {
		if err := db.Create(&setting).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Create(&models.User{ID: 1, Username: "author", Nickname: "作者 & 中文", Email: "PRIVATE EMAIL", Password: "PRIVATE PASSWORD"}).Error; err != nil {
		t.Fatal(err)
	}
	now := time.Now().Truncate(time.Millisecond)
	for _, a := range []models.Article{
		{ID: 1, Slug: "public-post", Title: "Public", Content: "Body", Status: 1, AuthorID: 1, PublishAt: &now},
		{ID: 2, Slug: "draft", Title: "Draft", Content: "secret", Status: 1, AuthorID: 1},
		{ID: 3, Slug: "private", Title: "Private", Content: "secret", Status: 2, AuthorID: 1},
		{ID: 4, Slug: "deleted", Title: "Deleted", Content: "secret", Status: 1, AuthorID: 1, DeletedAt: gorm.DeletedAt{Time: now, Valid: true}},
	} {
		if err := db.Create(&a).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Model(&models.Article{}).
		Where(models.ArticleColumns.ID+" = ?", 2).
		UpdateColumn(models.ArticleColumns.Status, 0).Error; err != nil {
		t.Fatal(err)
	}
	s := service.NewArticleService()
	for _, slug := range []string{"draft", "private", "deleted", "missing", "PUBLIC-POST", "123"} {
		if _, err := s.GetPublicBySlug(context.Background(), slug); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("%s: %v", slug, err)
		}
	}
	gin.SetMode(gin.TestMode)
	authorToken, err := utils.GenerateToken(1, "author", "admin")
	if err != nil {
		t.Fatal(err)
	}
	r := SetupUserRouter(fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}})
	for _, tt := range []struct {
		path   string
		status int
	}{
		{"/blog/public-post", 200}, {"/blog/1", 301}, {"/blog/PUBLIC-POST/", 301},
		{"/api/articles/slug/public-post", 200},
		{"/api/articles/slug/public-post?status=2&author_id=999&show_all=true&sort_order=asc", 200},
		{"/blog/draft", 404}, {"/blog/2", 404}, {"/blog/private", 404}, {"/blog/3", 404}, {"/blog/deleted", 404}, {"/blog/4", 404},
		{"/api/articles/slug/draft?status=0&author_id=1&show_all=true", 404},
		{"/api/articles/slug/private?status=2&show_all=true", 404}, {"/api/articles/slug/deleted", 404},
	} {
		t.Run(tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if w.Code != tt.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tt.status == 301 && w.Header().Get("Location") != "https://canonical.example/blog/public-post" {
				t.Fatal(w.Header())
			}
			if tt.status == 200 && strings.HasPrefix(tt.path, "/api/") {
				var response struct {
					Code int                    `json:"code"`
					Data models.ArticleResponse `json:"data"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Code != 0 || response.Data.ID != 1 || response.Data.PublishAt == nil || response.Data.PublishAt.UnixMilli() != now.UnixMilli() {
					t.Fatalf("public API publication contract: %s", w.Body.String())
				}
			}
			if tt.path == "/blog/public-post" {
				values, ld := parseHead(t, w.Body.String())
				if len(values["title"]) != 1 || values["title"][0] != "Public | InkSpace" || ld["url"] != "https://canonical.example/blog/public-post" {
					t.Fatal(values, ld)
				}
			}
		})
	}
	for _, path := range []string{
		"/blog/draft?status=0&show_all=true", "/blog/private?author_id=1&show_all=true", "/blog/deleted",
		"/api/articles/slug/draft?status=0&author_id=1&show_all=true",
		"/api/articles/slug/private?status=2&author_id=1&show_all=true",
	} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+authorToken)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("author/admin bypassed public filtering for %s: %d %s", path, w.Code, w.Body.String())
		}
	}
	// Cross-end acceptance on the same production assembly and real public SQL
	// readers. Host/forwarded headers and list-query bypass flags are untrusted.
	request := func(method, path, etag string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, nil)
		req.Host = "hostile.invalid"
		req.Header.Set("X-Forwarded-Host", "forwarded.invalid")
		req.Header.Set("Authorization", "Bearer "+authorToken)
		req.Header.Set("If-None-Match", etag)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		for _, secret := range []string{"PRIVATE", "hostile.invalid", "forwarded.invalid", "attacker.invalid", "javascript:", "secret"} {
			if strings.Contains(w.Body.String(), secret) {
				t.Fatalf("%s leaked %q", path, secret)
			}
		}
		return w
	}
	settings := request("GET", "/api/settings/public", "")
	var publicSettings struct {
		Code int               `json:"code"`
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(settings.Body.Bytes(), &publicSettings); err != nil {
		t.Fatal(err)
	}
	if settings.Code != 200 || publicSettings.Code != 0 || publicSettings.Data["public_site_url"] != "https://canonical.example" || publicSettings.Data["site_logo"] != "" {
		t.Fatal(settings.Code, publicSettings)
	}
	api := request("GET", "/api/articles/slug/public-post?status=0&show_all=true", "")
	var apiArticle struct {
		Code int `json:"code"`
		Data struct {
			SEO struct {
				Canonical string `json:"canonical"`
			} `json:"seo"`
		} `json:"data"`
	}
	if err := json.Unmarshal(api.Body.Bytes(), &apiArticle); err != nil {
		t.Fatal(err)
	}
	if api.Code != 200 || apiArticle.Code != 0 || apiArticle.Data.SEO.Canonical != "https://canonical.example/blog/public-post" {
		t.Fatal(api.Code, api.Body.String())
	}
	for _, path := range []string{"/feed.xml", "/atom.xml", "/sitemap.xml", "/sitemaps/1.xml", "/robots.txt"} {
		w := request("GET", path+"?status=0&show_all=true&author_id=1", "")
		if w.Code != 200 {
			t.Fatal(path, w.Code)
		}
		if path == "/feed.xml" || path == "/atom.xml" {
			entries := decodedFeedEntries(t, w.Body.Bytes())
			if len(entries) != 1 || entries[0].URL != "https://canonical.example/blog/public-post" {
				t.Fatal(entries)
			}
		} else if path != "/robots.txt" {
			doc := parseSitemap(t, w)
			if len(doc.URLs) != 3 || doc.URLs[2].Location != "https://canonical.example/blog/public-post" {
				t.Fatal(doc)
			}
		}
		head := request("HEAD", path, "")
		if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("ETag") != w.Header().Get("ETag") {
			t.Fatal(path, head.Code, head.Header())
		}
		for _, method := range []string{"GET", "HEAD"} {
			conditional := request(method, path, w.Header().Get("ETag"))
			if conditional.Code != 304 || conditional.Body.Len() != 0 {
				t.Fatal(path, conditional.Code)
			}
		}
	}
	// Visibility is fresh, not inherited from the internal article cache.
	if err := db.Model(&models.Article{}).
		Where(models.ArticleColumns.ID+" = ?", 1).
		UpdateColumn(models.ArticleColumns.Status, 2).Error; err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog/public-post", nil))
	if w.Code != 404 {
		t.Fatal("unpublished article remained public")
	}
	for _, path := range []string{"/feed.xml", "/atom.xml", "/sitemap.xml"} {
		w := request("GET", path, "")
		if w.Code != 200 || strings.Contains(w.Body.String(), "/blog/public-post") {
			t.Fatal("withdrawal not reflected", path, w.Code)
		}
	}
	// Destroy only a table in this disposable random test schema. Conditional
	// requests must not conceal failed DB reads as 304 or successful empty XML.
	if err := db.Migrator().DropTable(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/blog/public-post", "/blog/1", "/api/articles/slug/public-post", "/feed.xml", "/atom.xml", "/sitemap.xml", "/sitemaps/1.xml"} {
		methods := []string{"GET", "HEAD"}
		if strings.HasPrefix(path, "/api/") {
			methods = []string{"GET"} // existing JSON API is GET-only
		}
		for _, method := range methods {
			w := request(method, path, "*")
			if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("ETag") != "" || (method == "HEAD" && w.Body.Len() != 0) {
				t.Fatal(path, method, w.Code, w.Header())
			}
		}
	}
	if err := db.Migrator().DropTable(&models.Setting{}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/settings/public", "/"} {
		w := request("GET", path, "*")
		if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal(path, w.Code, w.Header())
		}
	}
	if w := request("GET", "/robots.txt", ""); w.Code != 200 {
		t.Fatal("robots must not depend on DB", w.Code)
	}
}
