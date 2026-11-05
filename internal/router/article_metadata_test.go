package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

const metadataHTML = `<!doctype html><html><head><meta charset="utf-8"><template id="inkspace-head-start"></template><title>old</title><template id="inkspace-head-end"></template><script type="module" src="/assets/app.js"></script></head><body><div id="app"></div></body></html>`

type metadataReader struct {
	pageReader
	article *models.Article
}

func (r metadataReader) GetPublicBySlug(_ context.Context, slug string) (*models.Article, error) {
	if r.article != nil {
		return r.article, nil
	}
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	return &models.Article{Slug: slug, Title: `中文 " & </script><script>alert(1)</script>`, Content: "# 正文\n\n<script>secret</script>\n\n**纯文本**", Cover: "/uploads/knowledge/private.png", PublishAt: &now, UpdatedAt: now.Add(time.Hour), Author: &models.User{Nickname: `作者 </script> & "`, Email: "secret@example.com", Role: "admin"}}, nil
}

func TestArticleMetadataOptionalFieldsAndEscaping(t *testing.T) {
	assets := fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}}
	r := gin.New()
	registerPublicArticlePages(r, metadataReader{article: &models.Article{
		Slug: "public-post", Title: "标题\\中文", Summary: `摘要 "引号" & <em>格式</em> \\ 中文`,
		Cover: "/uploads/cover.png",
	}}, articleMetadataRenderer(assets, fixedSite{}), "https://canonical.example")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog/public-post", nil))
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	values, ld := parseHead(t, w.Body.String())
	if values["description"][0] != `摘要 "引号" & 格式 \ 中文` {
		t.Fatal(values["description"])
	}
	if values["og:image"][0] != "https://canonical.example/uploads/cover.png" || values["twitter:card"][0] != "summary_large_image" {
		t.Fatal(values)
	}
	if ld["image"] != "https://canonical.example/uploads/cover.png" || ld["publisher"].(map[string]any)["logo"].(map[string]any)["url"] != "https://canonical.example/uploads/logo.png" {
		t.Fatal(ld)
	}
	for _, key := range []string{"author", "datePublished", "dateModified"} {
		if _, ok := ld[key]; ok {
			t.Fatal("missing field was invented", key)
		}
	}
}

type fixedSite struct{ fail bool }

func (s fixedSite) GetPublicSiteIdentity(context.Context) (seo.SiteIdentity, error) {
	if s.fail {
		return seo.SiteIdentity{}, context.DeadlineExceeded
	}
	return seo.SiteIdentity{URL: "https://canonical.example", Name: `站点 < & "`, Description: "默认描述", Logo: "/uploads/logo.png"}, nil
}

func parseHead(t *testing.T, body string) (map[string][]string, map[string]any) {
	t.Helper()
	root, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	values := map[string][]string{}
	var ld map[string]any
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		attrs := map[string]string{}
		for _, a := range n.Attr {
			attrs[a.Key] = a.Val
		}
		key, value := "", ""
		switch n.Data {
		case "title":
			key = "title"
			if n.FirstChild != nil {
				value = n.FirstChild.Data
			}
		case "meta":
			key = attrs["name"]
			if key == "" {
				key = attrs["property"]
			}
			value = attrs["content"]
		case "link":
			key = attrs["rel"]
			value = attrs["href"]
		case "script":
			if attrs["type"] == "application/ld+json" {
				if n.FirstChild == nil {
					t.Fatal("missing JSON-LD")
				}
				if err := json.Unmarshal([]byte(n.FirstChild.Data), &ld); err != nil {
					t.Fatal(err)
				}
				values["jsonld"] = append(values["jsonld"], "one")
			} else if attrs["src"] == "" && attrs["type"] != "application/json" {
				t.Fatal("injected executable script")
			}
		}
		if key != "" {
			values[key] = append(values[key], value)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(root)
	return values, ld
}

func TestArticleInitialHTMLMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	assets := fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}}
	r := gin.New()
	registerPublicArticlePages(r, metadataReader{}, articleMetadataRenderer(assets, fixedSite{}), "https://canonical.example")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://untrusted.example/blog/public-post", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 || w.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("%d %v", w.Code, w.Header())
	}
	values, ld := parseHead(t, w.Body.String())
	for key, want := range map[string]string{
		"title": `中文 " & </script><script>alert(1)</script> | 站点 < & "`, "description": "正文 纯文本",
		"canonical": "https://canonical.example/blog/public-post", "og:url": "https://canonical.example/blog/public-post",
		"og:title": `中文 " & </script><script>alert(1)</script>`, "og:description": "正文 纯文本",
		"twitter:card": "summary", "twitter:title": `中文 " & </script><script>alert(1)</script>`, "twitter:description": "正文 纯文本",
	} {
		if len(values[key]) != 1 || values[key][0] != want {
			t.Errorf("%s: %q want %q", key, values[key], want)
		}
	}
	if len(values["og:image"]) != 0 || len(values["twitter:image"]) != 0 || len(values["jsonld"]) != 1 {
		t.Fatal(values)
	}
	if ld["@type"] != "BlogPosting" || ld["datePublished"] != "2026-01-02T03:04:05Z" || ld["author"].(map[string]any)["name"] != `作者 </script> & "` {
		t.Fatal(ld)
	}
	if strings.Contains(w.Body.String(), "secret@example.com") || strings.Contains(w.Body.String(), "knowledge/private") || strings.Contains(w.Body.String(), "untrusted.example") {
		t.Fatal("private or request-host value leaked")
	}
	if len(values["alternate"]) != 2 || values["alternate"][0] != "https://canonical.example/feed.xml" || values["alternate"][1] != "https://canonical.example/atom.xml" {
		t.Fatal(values)
	}
	for _, site := range []fixedSite{{}, {fail: true}} {
		r = gin.New()
		registerPublicArticlePages(r, metadataReader{}, articleMetadataRenderer(assets, site), "https://canonical.example")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodHead, "/blog/public-post", nil))
		want := 200
		if site.fail {
			want = 500
		}
		if w.Code != want || w.Body.Len() != 0 {
			t.Fatalf("HEAD %d %s", w.Code, w.Body.String())
		}
		if site.fail && w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal(w.Header())
		}
	}
}

func TestGenericBlogHeadAndMissingPlaceholders(t *testing.T) {
	assets := fstest.MapFS{"index.html": {Data: []byte(metadataHTML)}, "assets/app.js": {Data: []byte("app")}}
	r := gin.New()
	serveSPA(r, assets, func(c *gin.Context) {
		site, err := (fixedSite{}).GetPublicSiteIdentity(c.Request.Context())
		if err != nil {
			t.Fatal(err)
		}
		serveBlogIndex(c, assets, site, nil)
	})
	for _, path := range []string{"/", "/blog", "/index.html"} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(method, path, nil))
			if w.Code != 200 || w.Header().Get("Cache-Control") != "no-cache" {
				t.Fatalf("%s: %d", path, w.Code)
			}
			if method == http.MethodHead {
				if w.Body.Len() != 0 {
					t.Fatal("HEAD body")
				}
				continue
			}
			values, ld := parseHead(t, w.Body.String())
			if len(values["alternate"]) != 2 || len(values["canonical"]) != 0 || ld != nil || !strings.Contains(w.Body.String(), "/assets/app.js") {
				t.Fatal(values)
			}
		}
	}
	for _, path := range []string{"/feed.xml", "/atom.xml", "/sitemap.xml", "/robots.txt"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != 404 {
			t.Fatal("discovery fell into SPA", path)
		}
	}
	r = gin.New()
	registerPublicArticlePages(r, metadataReader{}, articleMetadataRenderer(fstest.MapFS{"index.html": {Data: []byte("<title>no placeholders</title>")}}, fixedSite{}), "https://canonical.example")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog/public-post", nil))
	if w.Code != 503 || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("missing placeholders: %d", w.Code)
	}
}
