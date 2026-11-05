package seo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

// HeadStart and HeadEnd survive Vite HTML minification. Only this owned region
// is replaced, preserving built asset/preload tags elsewhere in the head.
const HeadStart = `<template id="inkspace-head-start"></template>`
const HeadEnd = `<template id="inkspace-head-end"></template>`

var headTemplate = template.Must(template.New("head").Parse(`
<title data-inkspace-meta>{{.Title}}</title>
<meta data-inkspace-meta name="description" content="{{.Description}}">
<link data-inkspace-feed rel="alternate" type="application/rss+xml" title="RSS" href="{{.RSS}}">
<link data-inkspace-feed rel="alternate" type="application/atom+xml" title="Atom" href="{{.Atom}}">
{{if .Article}}
<link data-inkspace-meta rel="canonical" href="{{.Article.Canonical}}">
<meta data-inkspace-meta property="og:type" content="article">
<meta data-inkspace-meta property="og:title" content="{{.Article.Headline}}">
<meta data-inkspace-meta property="og:description" content="{{.Article.Description}}">
<meta data-inkspace-meta property="og:url" content="{{.Article.Canonical}}">
{{if .Article.Image}}<meta data-inkspace-meta property="og:image" content="{{.Article.Image}}">{{end}}
<meta data-inkspace-meta name="twitter:card" content="{{.Article.Card}}">
<meta data-inkspace-meta name="twitter:title" content="{{.Article.Headline}}">
<meta data-inkspace-meta name="twitter:description" content="{{.Article.Description}}">
{{if .Article.Image}}<meta data-inkspace-meta name="twitter:image" content="{{.Article.Image}}">{{end}}
<script data-inkspace-meta type="application/ld+json">{{.JSONLD}}</script>
{{end}}
<script id="inkspace-site" type="application/json">{{.Site}}</script>
`))

// RenderHead replaces the explicit owned HTML region using contextual escaping.
// JSON is encoded with encoding/json's HTML escaping before entering script
// contexts, so user strings can never terminate either script element.
func RenderHead(index []byte, site SiteIdentity, article *ArticleMetadata) ([]byte, error) {
	start, end := bytes.Index(index, []byte(HeadStart)), bytes.Index(index, []byte(HeadEnd))
	if start < 0 || end < start || bytes.Count(index, []byte(HeadStart)) != 1 || bytes.Count(index, []byte(HeadEnd)) != 1 {
		return nil, fmt.Errorf("missing or duplicate HTML head placeholders")
	}
	siteJSON, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}
	var articleJSON []byte
	if article != nil {
		articleJSON, err = json.Marshal(article.JSONLD)
		if err != nil {
			return nil, err
		}
	}
	data := struct {
		Title, Description, RSS, Atom string
		Article                       *ArticleMetadata
		Site, JSONLD                  template.JS
	}{Title: site.Name, Description: site.Description, RSS: AbsoluteURL(site.URL, "/feed.xml"), Atom: AbsoluteURL(site.URL, "/atom.xml"), Article: article,
		// The only trusted script values are bytes produced by json.Marshal above.
		Site: template.JS(siteJSON), JSONLD: template.JS(articleJSON)}
	if article != nil {
		data.Title, data.Description = article.Title, article.Description
	}
	var head bytes.Buffer
	if err := headTemplate.Execute(&head, data); err != nil {
		return nil, err
	}
	result := append([]byte(nil), index[:start+len(HeadStart)]...)
	result = append(result, head.Bytes()...)
	result = append(result, index[end:]...)
	return result, nil
}
