package router

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"

	"github.com/gin-gonic/gin"
)

type feedArticleReaderFunc func(context.Context) ([]*models.Article, error)

func (f feedArticleReaderFunc) GetRecentPublicArticles(ctx context.Context) ([]*models.Article, error) {
	return f(ctx)
}

type feedIdentityReaderFunc func(context.Context) (seo.SiteIdentity, error)

func (f feedIdentityReaderFunc) GetPublicSiteIdentity(ctx context.Context) (seo.SiteIdentity, error) {
	return f(ctx)
}

// decodedFeedEntries reads both protocols through namespace-aware XML tokens,
// independently of their production serialization structs.
func decodedFeedEntries(t *testing.T, body []byte) []struct{ Title, Summary, URL, Author, Published, Updated string } {
	t.Helper()
	var document struct {
		Channel struct {
			Items []struct {
				Title     string `xml:"title"`
				Summary   string `xml:"description"`
				URL       string `xml:"guid"`
				Author    string `xml:"http://purl.org/dc/elements/1.1/ creator"`
				Published string `xml:"pubDate"`
				Updated   string `xml:"http://www.w3.org/2005/Atom updated"`
			} `xml:"item"`
		} `xml:"channel"`
		Entries []struct {
			Title   string `xml:"title"`
			Summary string `xml:"summary"`
			URL     string `xml:"id"`
			Author  struct {
				Name string `xml:"name"`
			} `xml:"author"`
			Published string `xml:"published"`
			Updated   string `xml:"updated"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal(body, &document); err != nil {
		t.Fatalf("invalid XML: %v\n%s", err, body)
	}
	var entries []struct{ Title, Summary, URL, Author, Published, Updated string }
	for _, item := range document.Channel.Items {
		entries = append(entries, struct{ Title, Summary, URL, Author, Published, Updated string }{item.Title, item.Summary, item.URL, item.Author, item.Published, item.Updated})
	}
	for _, entry := range document.Entries {
		entries = append(entries, struct{ Title, Summary, URL, Author, Published, Updated string }{entry.Title, entry.Summary, entry.URL, entry.Author.Name, entry.Published, entry.Updated})
	}
	return entries
}

func TestFeedsHTTPPublicProjectionAndCaching(t *testing.T) {
	gin.SetMode(gin.TestMode)
	published := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	rows := []*models.Article{
		{ID: 9, Slug: "new-post", Title: `中文 <标题> & "quote"`, Content: "# 安全正文\n\n<script>secret script</script>\n\n<style>secret style</style>\n\n**纯文本** & 中文", PublishAt: &published, UpdatedAt: published.Add(time.Hour), Author: &models.User{Username: "fallback", Nickname: `作者 <张> & "李"`, Email: "secret@example.org", Password: "secret-password"}},
		{ID: 8, Slug: "old-post", Title: "Older", Summary: "**显式摘要**", Content: "DO NOT EXPOSE FULL BODY", CreatedAt: published.Add(-time.Hour), Author: &models.User{Username: "username"}},
		{ID: 7, Slug: "empty-post", Title: "Empty", Content: "<script>hidden</script>", CreatedAt: published.Add(-2 * time.Hour)},
	}
	var articleErr, siteErr error
	articleReader := feedArticleReaderFunc(func(ctx context.Context) ([]*models.Article, error) { return rows, articleErr })
	identityReader := feedIdentityReaderFunc(func(ctx context.Context) (seo.SiteIdentity, error) {
		return seo.SiteIdentity{URL: "https://canonical.example", Name: "公开站点", Description: "站点描述"}, siteErr
	})
	r := gin.New()
	registerFeeds(r, articleReader, identityReader)
	for _, path := range []string{"/feed.xml", "/atom.xml"} {
		get := httptest.NewRecorder()
		r.ServeHTTP(get, httptest.NewRequest(http.MethodGet, path+"?status=0&author_id=999&show_all=true&page_size=10000&sort_order=asc", nil))
		wantType := "application/rss+xml; charset=utf-8"
		if path == "/atom.xml" {
			wantType = "application/atom+xml; charset=utf-8"
		}
		if get.Code != 200 || get.Header().Get("Content-Type") != wantType || get.Header().Get("Cache-Control") != "public, max-age=300" || get.Header().Get("ETag") == "" {
			t.Fatalf("status=%d headers=%v body=%s", get.Code, get.Header(), get.Body.String())
		}
		entries := decodedFeedEntries(t, get.Body.Bytes())
		if len(entries) != 3 || entries[0].Title != rows[0].Title || entries[0].Summary != "安全正文 纯文本 & 中文" || entries[0].URL != "https://canonical.example/blog/new-post" || entries[0].Author != rows[0].Author.Nickname || entries[1].Summary != "显式摘要" || entries[1].Author != "username" || entries[2].Summary != "站点描述" || entries[2].Author != "公开站点" {
			t.Fatalf("public projection: %+v\n%s", entries, get.Body.String())
		}
		if entries[0].Updated != "2024-01-02T04:04:05Z" || entries[1].Updated != "2024-01-02T02:04:05Z" {
			t.Fatal(entries)
		}
		wantPub := "Tue, 02 Jan 2024 03:04:05 +0000"
		if path == "/atom.xml" {
			wantPub = "2024-01-02T03:04:05Z"
		}
		if entries[0].Published != wantPub {
			t.Fatal(entries)
		}
		for _, secret := range []string{"secret", "DO NOT EXPOSE", "password", "localhost"} {
			if strings.Contains(get.Body.String(), secret) {
				t.Fatalf("leaked %q", secret)
			}
		}
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, "https://attacker.example"+path, nil)
			req.Header.Set("If-None-Match", "W/"+get.Header().Get("ETag"))
			r.ServeHTTP(w, req)
			if w.Code != 304 || w.Body.Len() != 0 {
				t.Fatalf("conditional %s status=%d body=%s", method, w.Code, w.Body.String())
			}
		}
		head := httptest.NewRecorder()
		r.ServeHTTP(head, httptest.NewRequest(http.MethodHead, path, nil))
		if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("Content-Length") != get.Header().Get("Content-Length") || head.Header().Get("ETag") != get.Header().Get("ETag") || head.Header().Get("Content-Type") != wantType {
			t.Fatalf("HEAD status=%d headers=%v body=%s", head.Code, head.Header(), head.Body.String())
		}
		for _, failure := range []string{"articles", "site"} {
			if failure == "articles" {
				articleErr = errors.New("private SQL details")
			} else {
				siteErr = errors.New("private settings details")
			}
			for _, method := range []string{http.MethodGet, http.MethodHead} {
				w := httptest.NewRecorder()
				req := httptest.NewRequest(method, path, nil)
				req.Header.Set("If-None-Match", get.Header().Get("ETag"))
				r.ServeHTTP(w, req)
				if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("ETag") != "" || w.Body.Len() != 0 {
					t.Fatalf("%s failure %s: %d %v %s", failure, method, w.Code, w.Header(), w.Body.String())
				}
			}
			articleErr, siteErr = nil, nil
		}
	}
	// A fresh public snapshot is read even with an old validator; content changes
	// and withdrawals cannot accidentally receive a stale 304.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/atom.xml", nil))
	rows = rows[1:]
	next := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/atom.xml", nil)
	req.Header.Set("If-None-Match", w.Header().Get("ETag"))
	r.ServeHTTP(next, req)
	if next.Code != 200 || len(decodedFeedEntries(t, next.Body.Bytes())) != 2 || next.Header().Get("ETag") == w.Header().Get("ETag") {
		t.Fatal("public snapshot did not refresh")
	}
}
