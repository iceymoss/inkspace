package router

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/service"
)

type sitemapReaderFunc func(context.Context) ([]service.PublicSitemapArticle, error)

func (f sitemapReaderFunc) GetPublicSitemapArticles(ctx context.Context) ([]service.PublicSitemapArticle, error) {
	return f(ctx)
}

func TestSitemapHTTPShardBoundaries(t *testing.T) {
	sitemapConfig(t)
	for _, tt := range []struct {
		articles int
		shards   []int
	}{
		{0, nil}, {49997, nil}, {49998, nil},
		{49999, []int{50000, 1}}, {50000, []int{50000, 2}},
		{99998, []int{50000, 50000}}, {99999, []int{50000, 50000, 1}},
	} {
		t.Run(fmt.Sprint(tt.articles), func(t *testing.T) {
			articles := make([]service.PublicSitemapArticle, tt.articles)
			for i := range articles {
				// Sparse IDs catch accidentally treating IDs as contiguous offsets.
				articles[i] = service.PublicSitemapArticle{ID: uint(3*i + 7), Slug: fmt.Sprintf("post-%d", 3*i+7)}
			}
			reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) { return articles, nil })
			router := setupUserRouterWithSitemap(reader)
			root := sitemapRequest(router, http.MethodGet, "/sitemap.xml", "")
			if root.Code != http.StatusOK {
				t.Fatalf("root status %d: %s", root.Code, root.Body.String())
			}
			doc := parseSitemap(t, root)
			if tt.shards == nil {
				if doc.XMLName.Local != "urlset" || len(doc.URLs) != tt.articles+2 || len(doc.Shards) != 0 {
					t.Fatalf("expected one sitemap, got %s urls=%d shards=%d", doc.XMLName.Local, len(doc.URLs), len(doc.Shards))
				}
				return
			}
			if doc.XMLName.Local != "sitemapindex" || len(doc.Shards) != len(tt.shards) || len(doc.URLs) != 0 {
				t.Fatalf("expected index, got %s urls=%d shards=%d", doc.XMLName.Local, len(doc.URLs), len(doc.Shards))
			}
			seen := make(map[string]bool)
			position := 0
			for i, shard := range doc.Shards {
				path := fmt.Sprintf("/sitemaps/%d.xml", i+1)
				if shard.Location != "https://canonical.example"+path {
					t.Fatalf("unstable/noncanonical shard: %q", shard.Location)
				}
				response := sitemapRequest(router, http.MethodGet, path+"?show_all=true&status=0&page=999", "")
				if response.Code != http.StatusOK {
					t.Fatalf("shard %d status=%d", i+1, response.Code)
				}
				part := parseSitemap(t, response)
				if part.XMLName.Local != "urlset" || len(part.URLs) != tt.shards[i] {
					t.Fatalf("shard %d size=%d want=%d", i+1, len(part.URLs), tt.shards[i])
				}
				for _, entry := range part.URLs {
					if seen[entry.Location] {
						t.Fatalf("duplicate URL: %q", entry.Location)
					}
					seen[entry.Location] = true
					want := "https://canonical.example/"
					if position == 1 {
						want = "https://canonical.example/blog"
					} else if position >= 2 {
						want = fmt.Sprintf("https://canonical.example/blog/post-%d", 3*(position-2)+7)
					}
					if entry.Location != want {
						t.Fatalf("missing or reordered URL at %d: %s want %s", position, entry.Location, want)
					}
					position++
				}
			}
			if position != tt.articles+2 {
				t.Fatalf("union=%d, want %d", position, tt.articles+2)
			}
			for _, invalid := range []string{"/sitemaps/0.xml", "/sitemaps/01.xml", "/sitemaps/-1.xml", "/sitemaps/no.xml", "/sitemaps/999999999999999999999.xml", fmt.Sprintf("/sitemaps/%d.xml", len(tt.shards)+1)} {
				for _, method := range []string{http.MethodGet, http.MethodHead} {
					response := sitemapRequest(router, method, invalid, "*")
					if response.Code != http.StatusNotFound || response.Header().Get("Cache-Control") != "no-store" || response.Body.Len() != 0 {
						t.Fatalf("invalid shard %s %s: %d %v %s", method, invalid, response.Code, response.Header(), response.Body.String())
					}
				}
			}
		})
	}
}

type parsedSitemap struct {
	XMLName xml.Name
	URLs    []struct {
		Location string `xml:"loc"`
		Modified string `xml:"lastmod"`
	} `xml:"url"`
	Shards []struct {
		Location string `xml:"loc"`
	} `xml:"sitemap"`
}

func parseSitemap(t *testing.T, response *httptest.ResponseRecorder) parsedSitemap {
	t.Helper()
	var doc parsedSitemap
	if err := xml.Unmarshal(response.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid XML: %v; %s", err, response.Body.String())
	}
	if doc.XMLName.Space != "http://www.sitemaps.org/schemas/sitemap/0.9" {
		t.Fatalf("wrong XML namespace: %s", response.Body.String())
	}
	return doc
}

func sitemapConfig(t *testing.T) {
	t.Helper()
	old := config.AppConfig
	config.AppConfig = &config.Config{PublicSiteURL: "https://canonical.example"}
	t.Cleanup(func() { config.AppConfig = old })
}

func sitemapRequest(router http.Handler, method, path, etag string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	request.Host = "attacker.example"
	request.Header.Set("X-Forwarded-Host", "attacker-forwarded.example")
	if etag != "" {
		request.Header.Set("If-None-Match", etag)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestSitemapHTTPPublicURLsAndLastmod(t *testing.T) {
	sitemapConfig(t)
	created := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	published, modified := created.Add(time.Hour), created.Add(2*time.Hour)
	reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) {
		return []service.PublicSitemapArticle{
			{ID: 1, Slug: "content-time", CreatedAt: created, PublishAt: &published, UpdatedAt: modified},
			{ID: 2, Slug: "publication-time", CreatedAt: created, PublishAt: &published},
			{ID: 3, Slug: "creation-time", CreatedAt: created, PublishAt: new(time.Time)},
			{ID: 4, Slug: "no-time"},
		}, nil
	})
	router := setupUserRouterWithSitemap(reader)
	response := sitemapRequest(router, http.MethodGet, "/sitemap.xml?status=0&show_all=true&author_id=123&page_size=1", "")
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Fatalf("response=%d %v %s", response.Code, response.Header(), response.Body.String())
	}
	doc := parseSitemap(t, response)
	if doc.XMLName.Local != "urlset" || len(doc.URLs) != 6 {
		t.Fatalf("unexpected sitemap: %s", response.Body.String())
	}
	for i, want := range []struct{ location, modified string }{
		{"https://canonical.example/", ""},
		{"https://canonical.example/blog", ""},
		{"https://canonical.example/blog/content-time", "2020-01-02T05:04:05Z"},
		{"https://canonical.example/blog/publication-time", "2020-01-02T04:04:05Z"},
		{"https://canonical.example/blog/creation-time", "2020-01-02T03:04:05Z"},
		{"https://canonical.example/blog/no-time", ""},
	} {
		if doc.URLs[i].Location != want.location || doc.URLs[i].Modified != want.modified {
			t.Fatalf("URL %d=%+v, want %+v", i, doc.URLs[i], want)
		}
	}
}

func TestSitemapAndRobotsHTTPRepresentation(t *testing.T) {
	sitemapConfig(t)
	articles := []service.PublicSitemapArticle{{ID: 1, Slug: "visible"}}
	var readError error
	reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) { return articles, readError })
	router := setupUserRouterWithSitemap(reader, fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}})
	for _, tt := range []struct{ path, contentType string }{
		{"/sitemap.xml", "application/xml; charset=utf-8"},
		{"/sitemaps/1.xml", "application/xml; charset=utf-8"},
		{"/robots.txt", "text/plain; charset=utf-8"},
	} {
		t.Run(tt.path, func(t *testing.T) {
			get := sitemapRequest(router, http.MethodGet, tt.path, "")
			etag := get.Header().Get("ETag")
			if get.Code != 200 || get.Header().Get("Content-Type") != tt.contentType || get.Header().Get("Cache-Control") != "public, max-age=300" || etag == "" || get.Body.Len() == 0 {
				t.Fatalf("GET: %d %v %s", get.Code, get.Header(), get.Body.String())
			}
			if strings.Contains(get.Body.String(), "attacker") || strings.Contains(get.Body.String(), "<!doctype") {
				t.Fatal("discovery fell through to SPA or trusted request host")
			}
			if tt.path == "/robots.txt" && !strings.Contains(get.Body.String(), "Sitemap: https://canonical.example/sitemap.xml\n") {
				t.Fatal(get.Body.String())
			}
			head := sitemapRequest(router, http.MethodHead, tt.path, "")
			if head.Code != get.Code || head.Body.Len() != 0 {
				t.Fatalf("HEAD: %d %s", head.Code, head.Body.String())
			}
			for _, header := range []string{"Content-Type", "Cache-Control", "ETag", "Content-Length"} {
				if head.Header().Get(header) != get.Header().Get(header) {
					t.Fatalf("HEAD %s=%q GET=%q", header, head.Header().Get(header), get.Header().Get(header))
				}
			}
			for _, condition := range []string{etag, "W/" + etag, `"other", W/` + etag, "*"} {
				for _, method := range []string{http.MethodGet, http.MethodHead} {
					conditional := sitemapRequest(router, method, tt.path, condition)
					if conditional.Code != 304 || conditional.Body.Len() != 0 || conditional.Header().Get("ETag") != etag {
						t.Fatalf("conditional %s %s: %d %v %s", method, condition, conditional.Code, conditional.Header(), conditional.Body.String())
					}
				}
			}
			mismatch := sitemapRequest(router, http.MethodGet, tt.path, `"different"`)
			if mismatch.Code != 200 || mismatch.Body.String() != get.Body.String() {
				t.Fatal("nonmatching ETag did not return representation")
			}
		})
	}
	before := sitemapRequest(router, http.MethodGet, "/sitemap.xml", "")
	articles = nil // A fresh query must remove unpublished articles, including for conditional requests.
	after := sitemapRequest(router, http.MethodGet, "/sitemap.xml", before.Header().Get("ETag"))
	if after.Code != 200 || after.Header().Get("ETag") == before.Header().Get("ETag") || strings.Contains(after.Body.String(), "/blog/visible") {
		t.Fatalf("stale public visibility: %d %s", after.Code, after.Body.String())
	}
	readError = errors.New("private SQL detail")
	for _, path := range []string{"/sitemap.xml", "/sitemaps/1.xml"} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			failure := sitemapRequest(router, method, path, "*")
			if failure.Code != 500 || failure.Header().Get("Cache-Control") != "no-store" || failure.Header().Get("ETag") != "" || strings.Contains(failure.Body.String(), "SQL") || (method == http.MethodHead && failure.Body.Len() != 0) {
				t.Fatalf("failure cached/masked: %d %v %s", failure.Code, failure.Header(), failure.Body.String())
			}
		}
	}
	// robots only needs validated configuration, not any DB/private settings.
	if response := sitemapRequest(router, http.MethodGet, "/robots.txt", ""); response.Code != 200 {
		t.Fatalf("robots unnecessarily depends on database: %d", response.Code)
	}
}

func TestSitemapAndRobotsRejectInvalidOrigin(t *testing.T) {
	sitemapConfig(t)
	config.AppConfig.PublicSiteURL = ""
	reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) { return nil, nil })
	router := setupUserRouterWithSitemap(reader)
	for _, path := range []string{"/sitemap.xml", "/sitemaps/1.xml", "/robots.txt"} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			response := sitemapRequest(router, method, path, "*")
			if response.Code != 500 || response.Header().Get("Cache-Control") != "no-store" || strings.Contains(response.Body.String(), "attacker") || (method == http.MethodHead && response.Body.Len() != 0) {
				t.Fatalf("invalid config: %d %v %s", response.Code, response.Header(), response.Body.String())
			}
		}
	}
}

func TestSitemapIndexHTTPHeadAndVisibilityChanges(t *testing.T) {
	sitemapConfig(t)
	articles := make([]service.PublicSitemapArticle, 49999)
	for i := range articles {
		articles[i] = service.PublicSitemapArticle{ID: uint(i + 1), Slug: fmt.Sprintf("post-%d", i+1)}
	}
	reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) { return articles, nil })
	router := setupUserRouterWithSitemap(reader)
	get := sitemapRequest(router, http.MethodGet, "/sitemap.xml", "")
	if get.Code != 200 || len(parseSitemap(t, get).Shards) != 2 {
		t.Fatalf("index not returned: %d %s", get.Code, get.Body.String())
	}
	head := sitemapRequest(router, http.MethodHead, "/sitemap.xml", "")
	if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("ETag") != get.Header().Get("ETag") || head.Header().Get("Content-Length") != get.Header().Get("Content-Length") {
		t.Fatalf("index HEAD differs: %d %v", head.Code, head.Header())
	}
	if conditional := sitemapRequest(router, http.MethodGet, "/sitemap.xml", get.Header().Get("ETag")); conditional.Code != 304 || conditional.Body.Len() != 0 {
		t.Fatal("index conditional GET failed")
	}
	partBefore := sitemapRequest(router, http.MethodGet, "/sitemaps/1.xml", "")
	// An edit changes lastmod/ETag but cannot move an ID into a different shard.
	articles[49997].UpdatedAt = time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	partAfter := sitemapRequest(router, http.MethodGet, "/sitemaps/1.xml", partBefore.Header().Get("ETag"))
	if partAfter.Code != 200 || partAfter.Header().Get("ETag") == partBefore.Header().Get("ETag") {
		t.Fatal("content edit failed to invalidate shard")
	}
	before, after := parseSitemap(t, partBefore), parseSitemap(t, partAfter)
	for i := range before.URLs {
		if before.URLs[i].Location != after.URLs[i].Location {
			t.Fatal("content edit moved shard membership")
		}
	}
	// Independent requests are fresh, not snapshot-frozen. Removing one article
	// can shrink the index to a URL set and makes the old second shard disappear.
	articles = articles[1:]
	rootAfter := sitemapRequest(router, http.MethodGet, "/sitemap.xml", get.Header().Get("ETag"))
	if rootAfter.Code != 200 || len(parseSitemap(t, rootAfter).URLs) != 50000 || strings.Contains(rootAfter.Body.String(), "/blog/post-1<") {
		t.Fatal("root retained unpublished membership")
	}
	missing := sitemapRequest(router, http.MethodGet, "/sitemaps/2.xml", "*")
	if missing.Code != 404 || missing.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("removed shard incorrectly cached: %d %v", missing.Code, missing.Header())
	}
}

func TestSitemapRejectsInvalidPersistedSlug(t *testing.T) {
	sitemapConfig(t)
	for _, slug := range []string{"", "../private", "123", "UPPERCASE", "bad/segment"} {
		t.Run(slug, func(t *testing.T) {
			reader := sitemapReaderFunc(func(context.Context) ([]service.PublicSitemapArticle, error) {
				return []service.PublicSitemapArticle{{ID: 1, Slug: slug}}, nil
			})
			router := setupUserRouterWithSitemap(reader)
			for _, path := range []string{"/sitemap.xml", "/sitemaps/1.xml"} {
				response := sitemapRequest(router, http.MethodGet, path, "*")
				if response.Code != 500 || response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("ETag") != "" {
					t.Fatalf("unsafe canonical emitted: %d %v %s", response.Code, response.Header(), response.Body.String())
				}
			}
		})
	}
}
