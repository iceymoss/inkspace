package router

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"

	"gorm.io/gorm"
)

func TestGenericSPASettingsFailureKeepsShellButSEORemainsStrict(t *testing.T) {
	sitemapConfig(t)
	assets := []fs.FS{fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}}}
	feeds := feedArticleReaderFunc(func(context.Context) ([]*models.Article, error) { return nil, nil })
	sitemap := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) { return nil, context.DeadlineExceeded })
	router := setupUserRouterWithPublicSEO(assets, feeds, fixedSite{fail: true}, sitemap, metadataReader{})
	for _, path := range []string{"/", "/blog", "/login", "/dashboard", "/share/token", "/index.html"} {
		for _, method := range []string{"GET", "HEAD"} {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(method, "http://hostile.invalid"+path, nil))
			if w.Code != 200 || w.Header().Get("Cache-Control") != "no-cache" {
				t.Fatalf("%s %s: %d %v", method, path, w.Code, w.Header())
			}
			if method == "HEAD" {
				if w.Body.Len() != 0 {
					t.Fatal("HEAD body")
				}
				continue
			}
			values, ld := parseHead(t, w.Body.String())
			if len(values["alternate"]) != 2 || values["alternate"][0] != "https://canonical.example/feed.xml" || values["alternate"][1] != "https://canonical.example/atom.xml" || values["title"][0] != "InkSpace" || ld != nil || !strings.Contains(w.Body.String(), `"url":"https://canonical.example"`) {
				t.Fatal(path, values, w.Body.String())
			}
		}
	}
	for _, path := range []string{"/blog/public-post", "/feed.xml", "/atom.xml", "/sitemap.xml"} {
		for _, method := range []string{"GET", "HEAD"} {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(method, path, nil))
			if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal(path, w.Code, w.Header())
			}
		}
	}
	// A settings fallback must not silently remove the required absolute links.
	config.AppConfig.PublicSiteURL = ""
	w := httptest.NewRecorder()
	router = setupUserRouterWithPublicSEO(assets, feeds, fixedSite{fail: true}, sitemap, metadataReader{})
	router.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal(w.Code, w.Header())
	}
}

// acceptanceArticles supplies already-public rows at the reader boundary. Real
// SQL visibility filtering is separately exercised by the gated MySQL fixture.
type acceptanceArticles struct {
	article *models.Article
	err     error
}

func (r *acceptanceArticles) GetPublicBySlug(_ context.Context, slug string) (*models.Article, error) {
	if r.err != nil {
		return nil, r.err
	}
	if slug != r.article.Slug {
		return nil, gorm.ErrRecordNotFound
	}
	return r.article, nil
}

func (r *acceptanceArticles) GetPublicByID(ctx context.Context, id uint) (*models.Article, error) {
	if id == r.article.ID {
		return r.GetPublicBySlug(ctx, r.article.Slug)
	}
	if r.err != nil {
		return nil, r.err
	}
	return nil, gorm.ErrRecordNotFound
}

func TestPublicSEOAcceptanceHTTP(t *testing.T) {
	sitemapConfig(t)
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	reader := &acceptanceArticles{article: &models.Article{
		ID: 1, Slug: "public-post", Title: `中文 & " </script><script>alert(1)</script>`,
		Summary: "**公开摘要**", Content: "PRIVATE FULL BODY", Cover: "javascript:alert(1)",
		PublishAt: &now, UpdatedAt: now, Author: &models.User{Nickname: "作者 & 中文", Email: "PRIVATE EMAIL"},
	}}
	feeds := feedArticleReaderFunc(func(context.Context) ([]*models.Article, error) { return []*models.Article{reader.article}, reader.err })
	sitemap := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) {
		return []service.PublicSitemapArticle{{ID: 1, Slug: "public-post", UpdatedAt: now}}, reader.err
	})
	assets := []fs.FS{fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}}}
	router := setupUserRouterWithPublicSEO(assets, feeds, fixedSite{}, sitemap, reader)
	request := func(method, path, etag string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, "http://hostile.invalid"+path, nil)
		req.Header.Set("X-Forwarded-Host", "forwarded.invalid")
		req.Header.Set("Forwarded", "host=forwarded.invalid;proto=http")
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		for _, secret := range []string{"hostile.invalid", "forwarded.invalid", "PRIVATE", "javascript:"} {
			if strings.Contains(w.Body.String(), secret) || strings.Contains(w.Header().Get("Location"), secret) {
				t.Fatalf("%s leaked %q", path, secret)
			}
		}
		return w
	}
	page := request("GET", "/blog/public-post", "")
	if page.Code != 200 {
		t.Fatal(page.Code)
	}
	values, ld := parseHead(t, page.Body.String())
	canonical := "https://canonical.example/blog/public-post"
	if values["canonical"][0] != canonical || values["og:url"][0] != canonical || ld["url"] != canonical || values["description"][0] != "公开摘要" || len(values["og:image"]) != 0 || len(values["jsonld"]) != 1 {
		t.Fatal(values, ld)
	}
	for _, method := range []string{"GET", "HEAD"} {
		redirect := request(method, "/blog/1", "")
		if redirect.Code != 301 || redirect.Header().Get("Location") != canonical {
			t.Fatal(redirect.Code, redirect.Header())
		}
		for _, path := range []string{"/blog/draft", "/blog/private", "/blog/deleted", "/blog/2", "/blog/missing"} {
			if w := request(method, path+"?status=0&show_all=true", ""); w.Code != 404 {
				t.Fatal(path, w.Code)
			}
		}
	}
	if head := request("HEAD", "/blog/public-post", ""); head.Code != 200 || head.Body.Len() != 0 {
		t.Fatal(head.Code, head.Body.String())
	}
	for _, path := range []string{"/", "/blog", "/index.html"} {
		w := request("GET", path, "")
		if w.Code != 200 {
			t.Fatal(path, w.Code)
		}
		generic, genericLD := parseHead(t, w.Body.String())
		if len(generic["alternate"]) != 2 || generic["alternate"][0] != "https://canonical.example/feed.xml" || genericLD != nil {
			t.Fatal(generic, genericLD)
		}
	}
	for path, mime := range map[string]string{"/feed.xml": "application/rss+xml", "/atom.xml": "application/atom+xml", "/sitemap.xml": "application/xml", "/sitemaps/1.xml": "application/xml", "/robots.txt": "text/plain"} {
		w := request("GET", path+"?status=0&show_all=true&author_id=99", "")
		if w.Code != 200 || w.Header().Get("Content-Type") != mime+"; charset=utf-8" || w.Header().Get("ETag") == "" {
			t.Fatal(path, w.Code, w.Header())
		}
		if strings.Contains(path, "feed") || strings.Contains(path, "atom") {
			entries := decodedFeedEntries(t, w.Body.Bytes())
			if len(entries) != 1 || entries[0].URL != canonical || entries[0].Summary != values["description"][0] {
				t.Fatal(entries)
			}
		} else if strings.Contains(path, "sitemap") {
			doc := parseSitemap(t, w)
			if len(doc.URLs) != 3 || doc.URLs[2].Location != canonical {
				t.Fatal(doc)
			}
		} else if !strings.Contains(w.Body.String(), "Sitemap: https://canonical.example/sitemap.xml") {
			t.Fatal(w.Body.String())
		}
		for _, method := range []string{"GET", "HEAD"} {
			conditional := request(method, path, "W/"+w.Header().Get("ETag"))
			if conditional.Code != 304 || conditional.Body.Len() != 0 {
				t.Fatal(path, conditional.Code)
			}
		}
		head := request("HEAD", path, "")
		if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("ETag") != w.Header().Get("ETag") || head.Header().Get("Content-Length") != w.Header().Get("Content-Length") {
			t.Fatal(path, head.Code, head.Header())
		}
	}
	reader.err = errors.New("PRIVATE SQL details")
	for _, path := range []string{"/blog/public-post", "/blog/1", "/feed.xml", "/atom.xml", "/sitemap.xml", "/sitemaps/1.xml"} {
		for _, method := range []string{"GET", "HEAD"} {
			w := request(method, path, "*")
			if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("ETag") != "" || w.Body.Len() != 0 {
				t.Fatal(path, w.Code, w.Header(), w.Body.String())
			}
		}
	}
	// robots is intentionally database-independent, including during outages.
	if w := request("GET", "/robots.txt", ""); w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}
}
