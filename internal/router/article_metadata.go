package router

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"
	"github.com/iceymoss/inkspace/internal/service"

	"github.com/gin-gonic/gin"
)

type publicSiteReader interface {
	GetPublicSiteIdentity(context.Context) (seo.SiteIdentity, error)
}

func articleMetadataRenderer(assets fs.FS, settings publicSiteReader) func(*gin.Context, *models.Article) {
	return func(c *gin.Context, article *models.Article) {
		site, err := settings.GetPublicSiteIdentity(c.Request.Context())
		if err != nil {
			htmlFailure(c, http.StatusInternalServerError)
			return
		}
		metadata := seo.BuildArticleMetadata(site, service.PublicArticleSEOInput(article))
		serveBlogIndex(c, assets, site, &metadata)
	}
}

func serveBlogIndex(c *gin.Context, assets fs.FS, site seo.SiteIdentity, metadata *seo.ArticleMetadata) {
	if assets == nil {
		htmlFailure(c, http.StatusServiceUnavailable)
		return
	}
	index, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		htmlFailure(c, http.StatusServiceUnavailable)
		return
	}
	body, err := seo.RenderHead(index, site, metadata)
	if err != nil {
		htmlFailure(c, http.StatusServiceUnavailable)
		return
	}
	c.Header("Cache-Control", "no-cache")
	c.Header("Content-Type", "text/html; charset=utf-8")
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
}

func htmlFailure(c *gin.Context, status int) { c.Header("Cache-Control", "no-store"); c.Status(status) }
