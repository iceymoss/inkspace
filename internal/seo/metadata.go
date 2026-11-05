package seo

import (
	"bytes"
	"strings"
	"time"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	xhtml "golang.org/x/net/html"
)

// SiteIdentity contains only public site settings and the validated canonical origin.
type SiteIdentity struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo,omitempty"`
}

// ArticleInput is the public, privacy-minimal input shared by page and feed builders.
// Callers must resolve visibility before constructing it. Zero dates are omitted.
type ArticleInput struct {
	Slug, Title, Summary, Content, Cover, AuthorName, AuthorURL string
	Published, Modified                                         time.Time
}

// ArticleMetadata is also returned by the public slug API so SPA navigation uses
// exactly the same descriptions, safe resources, canonical URLs and JSON-LD as Go.
type ArticleMetadata struct {
	Title       string         `json:"title"`
	Headline    string         `json:"headline"`
	Description string         `json:"description"`
	Canonical   string         `json:"canonical"`
	Image       string         `json:"image,omitempty"`
	Card        string         `json:"card"`
	JSONLD      map[string]any `json:"json_ld"`
}

// BuildArticleMetadata constructs metadata from already-public inputs. No request
// headers, account secrets or author profile objects enter this projection.
func BuildArticleMetadata(site SiteIdentity, article ArticleInput) ArticleMetadata {
	m := ArticleMetadata{
		Title: article.Title + " | " + site.Name, Headline: article.Title,
		Description: Summary(article.Summary, article.Content, site.Description),
		Canonical:   ArticleURL(site.URL, article.Slug), Image: PublicResourceURL(site.URL, article.Cover), Card: "summary",
	}
	if m.Image != "" {
		m.Card = "summary_large_image"
	}
	publisher := map[string]any{"@type": "Organization", "name": site.Name}
	if logo := PublicResourceURL(site.URL, site.Logo); logo != "" {
		publisher["logo"] = map[string]any{"@type": "ImageObject", "url": logo}
	}
	m.JSONLD = map[string]any{
		"@context": "https://schema.org", "@type": "BlogPosting", "headline": article.Title,
		"description": m.Description, "url": m.Canonical, "publisher": publisher,
	}
	if m.Image != "" {
		m.JSONLD["image"] = m.Image
	}
	if !article.Published.IsZero() {
		m.JSONLD["datePublished"] = article.Published.UTC().Format(time.RFC3339)
	}
	if !article.Modified.IsZero() {
		m.JSONLD["dateModified"] = article.Modified.UTC().Format(time.RFC3339)
	}
	if article.AuthorName != "" {
		author := map[string]any{"@type": "Person", "name": article.AuthorName}
		if u := PublicResourceURL(site.URL, article.AuthorURL); u != "" {
			author["url"] = u
		}
		m.JSONLD["author"] = author
	}
	return m
}

// Summary selects the first nonempty plain-text excerpt, limited to 160 Unicode
// code points. Markdown syntax, HTML, script/style contents, images and control
// characters never enter discovery documents; links contribute only their labels.
func Summary(summary, content, fallback string) string {
	for _, source := range []string{summary, content, fallback} {
		var rendered bytes.Buffer
		md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
		if err := md.Convert([]byte(source), &rendered); err != nil {
			continue
		}
		root, err := xhtml.Parse(&rendered)
		if err != nil {
			continue
		}
		var text strings.Builder
		var visit func(*xhtml.Node)
		visit = func(n *xhtml.Node) {
			if n.Type == xhtml.ElementNode {
				switch n.Data {
				case "script", "style", "img", "iframe", "object", "template", "noscript", "svg":
					return
				}
			}
			if n.Type == xhtml.TextNode {
				text.WriteString(n.Data)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				visit(c)
			}
			if n.Type == xhtml.ElementNode {
				switch n.Data {
				case "p", "div", "br", "li", "h1", "h2", "h3", "h4", "h5", "h6", "pre", "blockquote", "tr":
					text.WriteByte(' ')
				}
			}
		}
		visit(root)
		clean := strings.Map(func(r rune) rune {
			if unicode.IsControl(r) && !unicode.IsSpace(r) || unicode.Is(unicode.Cf, r) || r == '\ufffd' {
				return -1
			}
			return r
		}, text.String())
		clean = strings.Join(strings.Fields(clean), " ")
		if clean == "" {
			continue
		}
		runes := []rune(clean)
		if len(runes) > 160 {
			clean = string(runes[:160])
		}
		return clean
	}
	return ""
}
