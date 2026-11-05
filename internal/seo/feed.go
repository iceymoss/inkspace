package seo

import (
	"encoding/xml"
	"time"
)

const atomNamespace = "http://www.w3.org/2005/Atom"

// Feed is a protocol-neutral public projection. Callers supply canonical URLs,
// safe plain-text summaries and public author names from shared metadata rules.
// Entries retain the public reader's publication order; rendering never sorts.
type Feed struct {
	Title       string
	Description string
	HomeURL     string
	RSSURL      string
	AtomURL     string
	Entries     []FeedEntry
}

// FeedEntry contains only public summary metadata, never full article content or
// private author fields. Published is first publication, Updated is editorial time.
type FeedEntry struct {
	Title     string
	Summary   string
	Author    string
	URL       string
	Published time.Time
	Updated   time.Time
}

type feedLink struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Self        feedLink  `xml:"http://www.w3.org/2005/Atom link"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	GUID        rssGUID `xml:"guid"`
	// RSS author requires an email address. Dublin Core creator carries the same
	// public name as Atom author/name without inventing or leaking an email.
	Creator   string `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Published string `xml:"pubDate"`
	Updated   string `xml:"http://www.w3.org/2005/Atom updated"`
}

type rssGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type atomDocument struct {
	XMLName  xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	ID       string      `xml:"id"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Updated  string      `xml:"updated"`
	Links    []feedLink  `xml:"link"`
	Entries  []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID        string     `xml:"id"`
	Title     string     `xml:"title"`
	Link      feedLink   `xml:"link"`
	Summary   atomText   `xml:"summary"`
	Author    atomAuthor `xml:"author"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
}

type atomText struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

// RenderRSS encodes RSS 2.0 with an Atom self link and update extension, and a
// Dublin Core creator name. Only encoding/xml constructs markup or namespaces.
func RenderRSS(feed Feed) ([]byte, error) {
	channel := rssChannel{
		Title: feed.Title, Link: feed.HomeURL, Description: feed.Description,
		Self: feedLink{Rel: "self", Type: "application/rss+xml", Href: feed.RSSURL},
	}
	for _, entry := range feed.Entries {
		channel.Items = append(channel.Items, rssItem{
			Title: entry.Title, Link: entry.URL, Description: entry.Summary,
			GUID: rssGUID{IsPermaLink: true, Value: entry.URL}, Creator: entry.Author,
			Published: entry.Published.UTC().Format(time.RFC1123Z), Updated: entry.Updated.UTC().Format(time.RFC3339),
		})
	}
	return marshalFeedXML(rssDocument{Version: "2.0", Channel: channel})
}

// RenderAtom encodes Atom 1.0 using the same entries and timestamps as RSS.
// Empty feeds use a deterministic epoch update time instead of the request time,
// so unchanged empty representations remain conditionally cacheable.
func RenderAtom(feed Feed) ([]byte, error) {
	updated := time.Unix(0, 0).UTC()
	document := atomDocument{
		ID: feed.HomeURL, Title: feed.Title, Subtitle: feed.Description,
		Links: []feedLink{
			{Rel: "self", Type: "application/atom+xml", Href: feed.AtomURL},
			{Rel: "alternate", Type: "text/html", Href: feed.HomeURL},
		},
	}
	for _, entry := range feed.Entries {
		if entry.Updated.After(updated) {
			updated = entry.Updated
		}
		document.Entries = append(document.Entries, atomEntry{
			ID: entry.URL, Title: entry.Title,
			Link:    feedLink{Rel: "alternate", Type: "text/html", Href: entry.URL},
			Summary: atomText{Type: "text", Value: entry.Summary}, Author: atomAuthor{Name: entry.Author},
			Published: entry.Published.UTC().Format(time.RFC3339), Updated: entry.Updated.UTC().Format(time.RFC3339),
		})
	}
	document.Updated = updated.UTC().Format(time.RFC3339)
	return marshalFeedXML(document)
}

func marshalFeedXML(document any) ([]byte, error) {
	body, err := xml.Marshal(document)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), body...), nil
}
