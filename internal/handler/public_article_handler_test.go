package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type slugReader struct{}

func (slugReader) GetPublicBySlug(_ context.Context, slug string) (*models.Article, error) {
	if slug == "failure" {
		return nil, errors.New("private SQL error")
	}
	if slug != "public-post" {
		return nil, gorm.ErrRecordNotFound
	}
	return &models.Article{ID: 42, Slug: slug, Status: 1, Title: "Public"}, nil
}

func TestPublicSlugAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &PublicArticleHandler{reader: slugReader{}}
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
				Data *models.ArticleResponse
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if tt.status == 200 && (body.Code != 0 || body.Data.ID != 42 || body.Data.Slug != "public-post") {
				t.Fatalf("body %+v", body)
			}
		})
	}
}
