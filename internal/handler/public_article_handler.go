package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type publicSlugReader interface {
	GetPublicBySlug(context.Context, string) (*models.Article, error)
}

// PublicArticleHandler exposes only the public slug lookup, without list-query
// binding or author/admin visibility exceptions.
type PublicArticleHandler struct {
	reader     publicSlugReader
	recordView func(context.Context, uint) error
	site       func(context.Context) (seo.SiteIdentity, error)
}

// NewPublicArticleHandler constructs the uncached public article endpoint.
func NewPublicArticleHandler() *PublicArticleHandler {
	s := service.NewArticleService()
	return &PublicArticleHandler{reader: s, recordView: s.RecordPublicArticleView, site: service.NewSettingService().GetPublicSiteIdentity}
}

// GetBySlug ignores all query parameters; only the canonical slug is accepted.
func (h *PublicArticleHandler) GetBySlug(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	article, err := h.reader.GetPublicBySlug(c.Request.Context(), c.Param("slug"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "文章不存在")
		return
	}
	if err != nil {
		c.Header("Cache-Control", "no-store")
		utils.HTTPError(c, http.StatusInternalServerError, "文章加载失败")
		return
	}
	site, err := h.site(c.Request.Context())
	if err != nil {
		c.Header("Cache-Control", "no-store")
		utils.HTTPError(c, http.StatusInternalServerError, "文章加载失败")
		return
	}
	if h.recordView != nil {
		if err := h.recordView(c.Request.Context(), article.ID); err != nil {
			log.Printf("record public article view: %v", err)
		}
	}
	c.Header("Cache-Control", "no-cache")
	utils.Success(c, struct {
		*models.ArticleResponse
		Author *models.PublicUserResponse `json:"author,omitempty"`
		SEO    seo.ArticleMetadata        `json:"seo"`
	}{ArticleResponse: articleDetailResponse(article), Author: publicArticleAuthor(article), SEO: seo.BuildArticleMetadata(site, service.PublicArticleSEOInput(article))})
}

func publicArticleAuthor(article *models.Article) *models.PublicUserResponse {
	if article.Author == nil {
		return nil
	}
	return article.Author.ToPublicResponse()
}
