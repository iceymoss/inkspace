package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/iceymoss/inkspace/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pageReader struct{ fail bool }

func (p pageReader) GetPublicBySlug(_ context.Context, slug string) (*models.Article, error) {
	if p.fail {
		return nil, errors.New("database unavailable")
	}
	if slug != "public-post" {
		return nil, gorm.ErrRecordNotFound
	}
	return &models.Article{ID: 7, Slug: slug, Status: 1}, nil
}
func (p pageReader) GetPublicByID(ctx context.Context, id uint) (*models.Article, error) {
	if id != 7 {
		return nil, gorm.ErrRecordNotFound
	}
	return p.GetPublicBySlug(ctx, "public-post")
}

func TestPublicArticlePages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	assets := fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}}
	registerPublicArticlePages(r, pageReader{}, articleMetadataRenderer(assets, fixedSite{}), "https://canonical.example")
	serveSPA(r, assets)
	tests := []struct {
		path     string
		status   int
		location string
	}{
		{"/blog/public-post", 200, ""}, {"/blog/7", 301, "https://canonical.example/blog/public-post"},
		{"/blog/0007", 301, "https://canonical.example/blog/public-post"}, {"/blog/PUBLIC-POST/", 301, "https://canonical.example/blog/public-post"},
		{"/blog/public-post/", 301, "https://canonical.example/blog/public-post"}, {"/blog/private", 404, ""},
		{"/blog/%70ublic-post", 301, "https://canonical.example/blog/public-post"},
		{"/blog/draft?status=0&show_all=true&author_id=1", 404, ""},
		{"/blog/deleted", 404, ""}, {"/blog/8", 404, ""}, {"/blog/99999999999999999", 404, ""},
		{"/blog/0", 404, ""}, {"/blog/a.b", 404, ""}, {"/blog/a%2Fb", 404, ""},
	}
	for _, tt := range tests {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			t.Run(method+tt.path, func(t *testing.T) {
				w := httptest.NewRecorder()
				r.ServeHTTP(w, httptest.NewRequest(method, tt.path, nil))
				if w.Code != tt.status || w.Header().Get("Location") != tt.location {
					t.Fatalf("status %d location %q", w.Code, w.Header().Get("Location"))
				}
				if method == http.MethodHead && w.Body.Len() != 0 {
					t.Fatal("HEAD has body")
				}
				if tt.status == 404 && w.Body.Len() != 0 {
					t.Fatal("404 leaked article HTML")
				}
			})
		}
	}
	r = gin.New()
	registerPublicArticlePages(r, pageReader{fail: true}, articleMetadataRenderer(assets, fixedSite{}), "https://canonical.example")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog/public-post", nil))
	if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("database error was hidden/cached")
	}
}

func TestPublicArticleRoutesAreWired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupUserRouter(fstest.MapFS{"index.html": {Data: []byte("<html>app</html>")}})
	want := map[string]bool{
		"GET /api/articles/slug/:slug": false,
		"GET /api/articles/:id":        false,
		"GET /blog/:slug":              false,
		"HEAD /blog/:slug":             false,
		"GET /blog/:slug/":             false,
	}
	for _, route := range r.Routes() {
		key := route.Method + " " + route.Path
		if _, exists := want[key]; exists {
			want[key] = true
		}
	}
	for route, found := range want {
		if !found {
			t.Errorf("missing %s", route)
		}
	}
	for _, path := range []string{"/blog/a.b", "/blog/a/b", "/blog/a%2Fb"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusNotFound {
			t.Fatalf("%s status %d", path, w.Code)
		}
	}
}
