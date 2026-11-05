package seo

import (
	"encoding/xml"
	"time"
)

const sitemapNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

// SitemapURL is a public canonical URL and its optional reliable modification time.
type SitemapURL struct {
	Location     string
	LastModified time.Time
}

type sitemapEntry struct {
	Location string `xml:"loc"`
	Modified string `xml:"lastmod,omitempty"`
}

// RenderSitemap encodes a URL set using standard XML escaping. Zero timestamps
// are omitted rather than inventing modification times for static pages.
func RenderSitemap(urls []SitemapURL) ([]byte, error) {
	doc := struct {
		XMLName xml.Name       `xml:"urlset"`
		XMLNS   string         `xml:"xmlns,attr"`
		URLs    []sitemapEntry `xml:"url"`
	}{XMLNS: sitemapNamespace, URLs: make([]sitemapEntry, 0, len(urls))}
	for _, entry := range urls {
		item := sitemapEntry{Location: entry.Location}
		if !entry.LastModified.IsZero() {
			item.Modified = entry.LastModified.UTC().Format(time.RFC3339)
		}
		doc.URLs = append(doc.URLs, item)
	}
	body, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), body...), nil
}

// RenderSitemapIndex encodes stable shard URLs without synthetic lastmod values.
func RenderSitemapIndex(locations []string) ([]byte, error) {
	doc := struct {
		XMLName xml.Name       `xml:"sitemapindex"`
		XMLNS   string         `xml:"xmlns,attr"`
		Entries []sitemapEntry `xml:"sitemap"`
	}{XMLNS: sitemapNamespace, Entries: make([]sitemapEntry, 0, len(locations))}
	for _, location := range locations {
		doc.Entries = append(doc.Entries, sitemapEntry{Location: location})
	}
	body, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), body...), nil
}
