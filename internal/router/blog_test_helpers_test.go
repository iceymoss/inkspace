package router

import (
	"io/fs"

	"github.com/iceymoss/inkspace/internal/service"

	"github.com/gin-gonic/gin"
)

func setupUserRouter(assets []fs.FS, articles recentPublicArticleReader, settings feedSiteReader) *gin.Engine {
	public := service.NewArticleService()
	return setupUserRouterWithPublicSEO(assets, articles, settings, public, public)
}

func setupUserRouterWithSitemap(sitemap publicSitemapReader, assets ...fs.FS) *gin.Engine {
	articles := service.NewArticleService()
	return setupUserRouterWithPublicSEO(assets, articles, service.NewSettingService(), sitemap, articles)
}

func setupUserRouterWithDiscovery(assets []fs.FS, articles recentPublicArticleReader, settings feedSiteReader, sitemap publicSitemapReader) *gin.Engine {
	return setupUserRouterWithPublicSEO(assets, articles, settings, sitemap, service.NewArticleService())
}
