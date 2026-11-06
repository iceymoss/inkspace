package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/testutil"
	"github.com/iceymoss/inkspace/internal/utils"

	"gorm.io/gorm"
)

// TestSitemapHTTPMySQL exercises production router/service visibility against an
// opt-in random schema, never a developer database or a live Redis instance.
func TestSitemapHTTPMySQL(t *testing.T) {
	db := testutil.SEOTestMySQL(t)
	if err := db.AutoMigrate(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	old := database.DB
	database.DB = db
	t.Cleanup(func() { database.DB = old })
	sitemapConfig(t)
	config.AppConfig.JWT = config.JWTConfig{Secret: "sitemap-author-test", ExpireHours: 1}
	authorToken, err := utils.GenerateToken(1, "author", "admin")
	if err != nil {
		t.Fatal(err)
	}
	created := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	published, modified := created.Add(time.Hour), created.Add(2*time.Hour)
	for _, article := range []models.Article{
		{ID: 1, Slug: "public", Title: "public", Content: "body", AuthorID: 1, Status: 1, CreatedAt: created, UpdatedAt: modified, PublishAt: &published},
		{ID: 2, Slug: "draft-secret", Title: "draft", Content: "secret", AuthorID: 1, Status: 1},
		{ID: 3, Slug: "private-secret", Title: "private", Content: "secret", AuthorID: 1, Status: 2},
		{ID: 4, Slug: "deleted-secret", Title: "deleted", Content: "secret", AuthorID: 1, Status: 1, DeletedAt: gorm.DeletedAt{Time: modified, Valid: true}},
	} {
		if err := db.Create(&article).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Model(&models.Article{}).
		Where(models.ArticleColumns.ID+" = ?", 2).
		UpdateColumn(models.ArticleColumns.Status, 0).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupUserRouter()
	for _, path := range []string{
		"/sitemap.xml", "/sitemap.xml?status=0&show_all=true&author_id=1&page_size=1",
		"/sitemaps/1.xml?status=2&show_all=true&author_id=1",
	} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.Header.Set("Authorization", "Bearer "+authorToken)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		doc := parseSitemap(t, response)
		if response.Code != 200 || len(doc.URLs) != 3 || doc.URLs[2].Location != "https://canonical.example/blog/public" || doc.URLs[2].Modified != "2020-01-02T05:04:05Z" || strings.Contains(response.Body.String(), "secret") {
			t.Fatalf("private visibility bypass: %d %s", response.Code, response.Body.String())
		}
	}
	before := sitemapRequest(router, http.MethodGet, "/sitemap.xml", "")
	// Interaction counters must not create a new lastmod/ETag.
	if err := db.Model(&models.Article{}).
		Where(models.ArticleColumns.ID+" = ?", 1).
		UpdateColumn("view_count", 15).Error; err != nil {
		t.Fatal(err)
	}
	unchanged := sitemapRequest(router, http.MethodGet, "/sitemap.xml", before.Header().Get("ETag"))
	if unchanged.Code != http.StatusNotModified {
		t.Fatalf("interaction changed sitemap: %d %s", unchanged.Code, unchanged.Body.String())
	}
	if err := db.Model(&models.Article{}).
		Where(models.ArticleColumns.ID+" = ?", 1).
		UpdateColumn(models.ArticleColumns.Status, 2).Error; err != nil {
		t.Fatal(err)
	}
	after := sitemapRequest(router, http.MethodGet, "/sitemap.xml", before.Header().Get("ETag"))
	if after.Code != 200 || len(parseSitemap(t, after).URLs) != 2 || strings.Contains(after.Body.String(), "/blog/public") {
		t.Fatalf("unpublished article remained cached: %d %s", after.Code, after.Body.String())
	}
	if err := db.Migrator().DropTable(&models.Article{}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/sitemap.xml", "/sitemaps/1.xml"} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			response := sitemapRequest(router, method, path, "*")
			if response.Code != 500 || response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("ETag") != "" || (method == http.MethodHead && response.Body.Len() != 0) {
				t.Fatalf("DB failure masked/cached: %d %v %s", response.Code, response.Header(), response.Body.String())
			}
		}
	}
}
