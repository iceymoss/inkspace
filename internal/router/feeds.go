package router

import (
	"context"
	"strings"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"
	"github.com/iceymoss/inkspace/internal/service"

	"github.com/gin-gonic/gin"
)

// recentPublicArticleReader is the shared, fixed-limit publication query, not
// the general article list that accepts author/admin visibility filters.
type recentPublicArticleReader interface {
	GetRecentPublicArticles(context.Context) ([]*models.Article, error)
}

type feedSiteReader interface {
	GetPublicSiteIdentity(context.Context) (seo.SiteIdentity, error)
}

// registerFeeds installs root discovery endpoints independently of SPA assets.
// Both formats share the same public query and projection, with no request
// filters or request Host in the representation. Readers are fresh for HEAD and
// conditional GET as well, so database failures cannot become a cached empty feed.
func registerFeeds(r *gin.Engine, articles recentPublicArticleReader, settings feedSiteReader) {
	for _, format := range []struct {
		path        string
		contentType string
		render      func(seo.Feed) ([]byte, error)
	}{
		{"/feed.xml", "application/rss+xml; charset=utf-8", seo.RenderRSS},
		{"/atom.xml", "application/atom+xml; charset=utf-8", seo.RenderAtom},
	} {
		serve := func(c *gin.Context) {
			rows, err := articles.GetRecentPublicArticles(c.Request.Context())
			if err != nil {
				discoveryError(c)
				return
			}
			site, err := settings.GetPublicSiteIdentity(c.Request.Context())
			if err != nil {
				discoveryError(c)
				return
			}
			feed := seo.Feed{
				Title: site.Name, Description: site.Description,
				HomeURL: seo.AbsoluteURL(site.URL, "/"), RSSURL: seo.AbsoluteURL(site.URL, "/feed.xml"), AtomURL: seo.AbsoluteURL(site.URL, "/atom.xml"),
			}
			for _, article := range rows {
				if article == nil || !seo.ValidSlug(article.Slug) {
					discoveryError(c)
					return
				}
				input := service.PublicArticleSEOInput(article)
				author := input.AuthorName
				if strings.TrimSpace(author) == "" {
					// Atom requires an author for every entry or at feed level.
					// Use the public publisher name in both protocols when no
					// public author name is available; never invent an email.
					author = site.Name
				}
				feed.Entries = append(feed.Entries, seo.FeedEntry{
					Title: input.Title, Summary: seo.Summary(input.Summary, input.Content, site.Description),
					URL: seo.ArticleURL(site.URL, input.Slug), Author: author,
					Published: input.Published, Updated: input.Modified,
				})
			}
			body, err := format.render(feed)
			if err != nil {
				discoveryError(c)
				return
			}
			writeDiscovery(c, format.contentType, body)
		}
		r.GET(format.path, serve)
		r.HEAD(format.path, serve)
	}
}
