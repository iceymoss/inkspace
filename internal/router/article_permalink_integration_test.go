package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/testutil"

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
	config.AppConfig = &config.Config{PublicSiteURL: "https://canonical.example"}
	database.DB = db
	t.Cleanup(func() { database.DB = old; config.AppConfig = oldConfig })
	for _, setting := range []models.Setting{
		{Key: "site_name", Value: "PRIVATE SITE NAME", IsPublic: false},
		{Key: "site_description", Value: "公开描述", IsPublic: true},
	} {
		if err := db.Create(&setting).Error; err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	for _, a := range []models.Article{
		{ID: 1, Slug: "public-post", Title: "Public", Content: "Body", Status: 1, AuthorID: 1},
		{ID: 2, Slug: "draft", Title: "Draft", Content: "secret", Status: 1, AuthorID: 1},
		{ID: 3, Slug: "private", Title: "Private", Content: "secret", Status: 2, AuthorID: 1},
		{ID: 4, Slug: "deleted", Title: "Deleted", Content: "secret", Status: 1, AuthorID: 1, DeletedAt: gorm.DeletedAt{Time: now, Valid: true}},
	} {
		if err := db.Create(&a).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Model(&models.Article{}).Where("id=2").UpdateColumn("status", 0).Error; err != nil {
		t.Fatal(err)
	}
	s := service.NewArticleService()
	for _, slug := range []string{"draft", "private", "deleted", "missing", "PUBLIC-POST", "123"} {
		if _, err := s.GetPublicBySlug(context.Background(), slug); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("%s: %v", slug, err)
		}
	}
	gin.SetMode(gin.TestMode)
	r := SetupUserRouter(fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}})
	for _, tt := range []struct {
		path   string
		status int
	}{
		{"/blog/public-post", 200}, {"/blog/1", 301}, {"/blog/PUBLIC-POST/", 301},
		{"/api/articles/slug/public-post", 200},
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
			if tt.path == "/blog/public-post" {
				values, ld := parseHead(t, w.Body.String())
				if len(values["title"]) != 1 || values["title"][0] != "Public | InkSpace" || ld["url"] != "https://canonical.example/blog/public-post" {
					t.Fatal(values, ld)
				}
			}
		})
	}
	// Visibility is fresh, not inherited from the internal article cache.
	if err := db.Model(&models.Article{}).Where("id=1").UpdateColumn("status", 2).Error; err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog/public-post", nil))
	if w.Code != 404 {
		t.Fatal("unpublished article remained public")
	}
}
