package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type publicArticleReader interface {
	GetPublicBySlug(context.Context, string) (*models.Article, error)
	GetPublicByID(context.Context, uint) (*models.Article, error)
}

// registerPublicArticlePages owns article HTML URLs before generic SPA fallback.
// The resolved public article is passed to render so metadata can be added without
// reimplementing visibility, legacy-ID parsing or canonical-path normalization.
func registerPublicArticlePages(r *gin.Engine, reader publicArticleReader, render func(*gin.Context, *models.Article), base string) {
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/blog/") && strings.Contains(strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/blog/"), "/"), "/") {
			c.Header("Cache-Control", "no-cache")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	})
	page := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		value := strings.ToLower(c.Param("slug"))
		var article *models.Article
		var err error
		if seo.IsNumericArticlePath(value) {
			id, parseErr := strconv.ParseUint(value, 10, 32)
			if parseErr != nil || id == 0 {
				c.Status(http.StatusNotFound)
				return
			}
			article, err = reader.GetPublicByID(c.Request.Context(), uint(id))
		} else if seo.ValidSlug(value) {
			article, err = reader.GetPublicBySlug(c.Request.Context(), value)
		} else {
			c.Status(http.StatusNotFound)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		if err != nil {
			c.Header("Cache-Control", "no-store")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !seo.ValidSlug(article.Slug) {
			c.Header("Cache-Control", "no-store")
			c.Status(http.StatusInternalServerError)
			return
		}
		canonical := seo.ArticlePath(article.Slug)
		if c.Request.URL.EscapedPath() != canonical {
			origin, err := seo.NormalizePublicSiteURL(base)
			if err != nil {
				c.Header("Cache-Control", "no-store")
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Redirect(http.StatusMovedPermanently, seo.ArticleURL(origin, article.Slug))
			return
		}
		render(c, article)
	}
	for _, path := range []string{"/blog/:slug", "/blog/:slug/"} {
		r.GET(path, page)
		r.HEAD(path, page)
	}
}
