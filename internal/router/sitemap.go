package router

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/seo"
	"github.com/iceymoss/inkspace/internal/service"

	"github.com/gin-gonic/gin"
)

type publicSitemapReader interface {
	GetPublicSitemapArticles(context.Context) ([]service.PublicSitemapArticle, error)
}

const sitemapURLLimit = 50000

func registerSitemapRoutes(routes gin.IRoutes, reader publicSitemapReader, base string) {
	sitemap := func(c *gin.Context) {
		page := 0 // The root chooses URL set or index; shard numbering starts at 1.
		if name := c.Param("shard"); name != "" {
			number, ok := strings.CutSuffix(name, ".xml")
			parsed, err := strconv.Atoi(number)
			if !ok || err != nil || parsed < 1 || strconv.Itoa(parsed) != number {
				sitemapNotFound(c)
				return
			}
			page = parsed
		}
		origin, err := seo.NormalizePublicSiteURL(base)
		if err != nil {
			discoveryError(c)
			return
		}
		articles, err := reader.GetPublicSitemapArticles(c.Request.Context())
		if err != nil {
			discoveryError(c)
			return
		}
		for _, article := range articles {
			if !seo.ValidSlug(article.Slug) {
				discoveryError(c)
				return
			}
		}
		// One ordered compact SELECT supplies both count and membership. Edits do
		// not move IDs, but publication/deletion can shift membership between HTTP
		// requests: this is deliberately not a persistent cross-request snapshot.
		total := len(articles) + 2
		shards := (total-1)/sitemapURLLimit + 1
		if page > shards {
			sitemapNotFound(c)
			return
		}
		var body []byte
		if page == 0 && total > sitemapURLLimit {
			locations := make([]string, 0, shards)
			for shard := 1; shard <= shards; shard++ {
				locations = append(locations, seo.AbsoluteURL(origin, "/sitemaps/"+strconv.Itoa(shard)+".xml"))
			}
			body, err = seo.RenderSitemapIndex(locations)
		} else {
			if page == 0 {
				page = 1
			}
			start := (page - 1) * sitemapURLLimit
			end := min(start+sitemapURLLimit, total)
			urls := make([]seo.SitemapURL, 0, end-start)
			// The two static URLs participate in the SAME limit as articles. Only
			// shard 1 contains them; article offsets therefore subtract exactly 2.
			for position := start; position < end; position++ {
				switch position {
				case 0:
					urls = append(urls, seo.SitemapURL{Location: seo.AbsoluteURL(origin, "/")})
				case 1:
					urls = append(urls, seo.SitemapURL{Location: seo.AbsoluteURL(origin, "/blog")})
				default:
					article := articles[position-2]
					urls = append(urls, seo.SitemapURL{Location: seo.ArticleURL(origin, article.Slug), LastModified: article.LastModified()})
				}
			}
			body, err = seo.RenderSitemap(urls)
		}
		if err != nil {
			discoveryError(c)
			return
		}
		writeDiscovery(c, "application/xml; charset=utf-8", body)
	}
	routes.GET("/sitemap.xml", sitemap)
	routes.HEAD("/sitemap.xml", sitemap)
	routes.GET("/sitemaps/:shard", sitemap)
	routes.HEAD("/sitemaps/:shard", sitemap)
	robots := func(c *gin.Context) {
		origin, err := seo.NormalizePublicSiteURL(base)
		if err != nil {
			discoveryError(c)
			return
		}
		writeDiscovery(c, "text/plain; charset=utf-8", seo.RenderRobots(seo.AbsoluteURL(origin, "/sitemap.xml")))
	}
	routes.GET("/robots.txt", robots)
	routes.HEAD("/robots.txt", robots)
}

func sitemapNotFound(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Status(http.StatusNotFound)
}
