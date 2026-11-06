package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"
)

func TestPublicSlugAPIViewFailureUsesStructuredNonfatalLog(t *testing.T) {
	core, logs := observer.New(zap.WarnLevel)
	restore := zap.ReplaceGlobals(zap.New(core))
	t.Cleanup(restore)
	h := &PublicArticleHandler{reader: slugReader{}, site: func(context.Context) (seo.SiteIdentity, error) {
		return seo.SiteIdentity{URL: "https://canonical.example", Name: "Site"}, nil
	}, recordView: func(context.Context, uint) error { return errors.New("view write failed") }}
	r := gin.New()
	r.GET("/api/articles/slug/:slug", h.GetBySlug)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/articles/slug/public-post", nil))
	if w.Code != 200 || strings.Contains(w.Body.String(), "view write failed") {
		t.Fatal(w.Code, w.Body.String())
	}
	entries := logs.All()
	if len(entries) != 1 || entries[0].Message != "record public article view" || entries[0].ContextMap()["article_id"] != uint64(42) || entries[0].ContextMap()["error"] != "view write failed" {
		t.Fatal(entries)
	}
}

type slugReader struct{}

func (slugReader) GetPublicBySlug(_ context.Context, slug string) (*models.Article, error) {
	if slug == "failure" {
		return nil, errors.New("private SQL error")
	}
	if slug != "public-post" {
		return nil, gorm.ErrRecordNotFound
	}
	return &models.Article{ID: 42, Slug: slug, Status: 1, Title: "Public", Author: &models.User{Nickname: "Public author", Email: "secret@example.com", Role: "admin"}}, nil
}

func TestPublicSlugAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &PublicArticleHandler{reader: slugReader{}, site: func(context.Context) (seo.SiteIdentity, error) {
		return seo.SiteIdentity{URL: "https://canonical.example", Name: "Public site"}, nil
	}}
	r.GET("/api/articles/slug/:slug", h.GetBySlug)
	for _, tt := range []struct {
		slug   string
		status int
	}{
		{"public-post", 200}, {"draft", 404}, {"private", 404}, {"deleted", 404}, {"missing", 404}, {"failure", 500},
	} {
		t.Run(tt.slug, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/articles/slug/"+tt.slug+"?status=0&show_all=true&author_id=42", nil))
			if w.Code != tt.status {
				t.Fatalf("status %d", w.Code)
			}
			var body struct {
				Code int
				Data *struct {
					*models.ArticleResponse
					SEO seo.ArticleMetadata `json:"seo"`
				}
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if tt.status == 200 && (body.Code != 0 || body.Data.ID != 42 || body.Data.Slug != "public-post") {
				t.Fatalf("body %+v", body)
			}
			if tt.status == 200 && (body.Data.SEO.Title != "Public | Public site" || body.Data.SEO.Canonical != "https://canonical.example/blog/public-post") {
				t.Fatalf("missing shared metadata: %+v", body.Data)
			}
			if strings.Contains(w.Body.String(), "secret@example.com") || strings.Contains(w.Body.String(), `"role"`) {
				t.Fatal("private account fields leaked")
			}
		})
	}
}

func TestPublicSlugAPISiteFailureIsNotCached(t *testing.T) {
	r := gin.New()
	h := &PublicArticleHandler{reader: slugReader{}, site: func(context.Context) (seo.SiteIdentity, error) {
		return seo.SiteIdentity{}, errors.New("private database details")
	}}
	r.GET("/api/articles/slug/:slug", h.GetBySlug)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/articles/slug/public-post", nil))
	if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || strings.Contains(w.Body.String(), "private database") {
		t.Fatal(w.Code, w.Header(), w.Body.String())
	}
}
