package seo

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestRenderSitemapProtocol(t *testing.T) {
	modified := time.Date(2026, 1, 2, 3, 4, 5, 0, time.FixedZone("east", 8*3600))
	body, err := RenderSitemap([]SitemapURL{
		{Location: "https://example.com/"},
		{Location: "https://example.com/blog"},
		{Location: "https://example.com/blog/chinese?a=1&b=中文", LastModified: modified},
	})
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		XMLName xml.Name
		URLs    []struct {
			Location string `xml:"loc"`
			Modified string `xml:"lastmod"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.XMLName.Local != "urlset" || doc.XMLName.Space != "http://www.sitemaps.org/schemas/sitemap/0.9" || len(doc.URLs) != 3 {
		t.Fatalf("invalid sitemap: %s", body)
	}
	if doc.URLs[0].Modified != "" || doc.URLs[1].Modified != "" || doc.URLs[2].Modified != "2026-01-01T19:04:05Z" {
		t.Fatalf("unreliable or incorrect lastmod: %+v", doc.URLs)
	}
	if doc.URLs[2].Location != "https://example.com/blog/chinese?a=1&b=中文" || !strings.Contains(string(body), "&amp;") {
		t.Fatalf("XML did not roundtrip: %s", body)
	}
}

func TestRenderSitemapIndexProtocol(t *testing.T) {
	body, err := RenderSitemapIndex([]string{"https://example.com/sitemaps/1.xml", "https://example.com/sitemaps/2.xml"})
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		XMLName   xml.Name
		Locations []string `xml:"sitemap>loc"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.XMLName.Local != "sitemapindex" || doc.XMLName.Space != "http://www.sitemaps.org/schemas/sitemap/0.9" || len(doc.Locations) != 2 || doc.Locations[1] != "https://example.com/sitemaps/2.xml" {
		t.Fatalf("invalid index: %s", body)
	}
	if strings.Contains(string(body), "lastmod") {
		t.Fatalf("index invents modification times: %s", body)
	}
}
