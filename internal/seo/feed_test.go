package seo

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestFeedsRoundTripSharedSemantics(t *testing.T) {
	published := time.Date(2024, 6, 7, 8, 9, 10, 0, time.FixedZone("CST", 8*60*60))
	updated := published.Add(2 * time.Hour)
	feed := Feed{
		Title: `中文 & "Site"`, Description: "站点 <描述>", HomeURL: "https://blog.example/",
		RSSURL: "https://blog.example/feed.xml", AtomURL: "https://blog.example/atom.xml",
		Entries: []FeedEntry{{
			Title: `文章 <标题> & "quote"`, Summary: `摘要中文 & <纯文本>`, Author: `作者 "张" & 李`,
			URL: "https://blog.example/blog/stable-slug", Published: published, Updated: updated,
		}},
	}
	rssBody, err := RenderRSS(feed)
	if err != nil {
		t.Fatal(err)
	}
	var rss struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Links       []struct {
				XMLName xml.Name
				Value   string `xml:",chardata"`
				Href    string `xml:"href,attr"`
				Rel     string `xml:"rel,attr"`
				Type    string `xml:"type,attr"`
			} `xml:"link"`
			Items []struct {
				Title       string `xml:"title"`
				Description string `xml:"description"`
				Link        string `xml:"link"`
				GUID        struct {
					Value       string `xml:",chardata"`
					IsPermaLink string `xml:"isPermaLink,attr"`
				} `xml:"guid"`
				Author    string `xml:"author"`
				Creator   string `xml:"http://purl.org/dc/elements/1.1/ creator"`
				Published string `xml:"pubDate"`
				Updated   string `xml:"http://www.w3.org/2005/Atom updated"`
			} `xml:"item"`
		} `xml:"channel"`
	}
	if err := xml.Unmarshal(rssBody, &rss); err != nil {
		t.Fatal(err)
	}
	if rss.Version != "2.0" || rss.Channel.Title != feed.Title || rss.Channel.Description != feed.Description || len(rss.Channel.Items) != 1 {
		t.Fatalf("RSS envelope: %s", rssBody)
	}
	if len(rss.Channel.Links) != 2 || rss.Channel.Links[0].XMLName.Space != "" || rss.Channel.Links[0].Value != feed.HomeURL || rss.Channel.Links[1].XMLName.Space != atomNamespace || rss.Channel.Links[1].Rel != "self" || rss.Channel.Links[1].Href != feed.RSSURL || rss.Channel.Links[1].Type != "application/rss+xml" {
		t.Fatal("RSS self link", string(rssBody))
	}
	item := rss.Channel.Items[0]
	want := feed.Entries[0]
	if item.Title != want.Title || item.Description != want.Summary || item.Creator != want.Author || item.Author != "" || item.Link != want.URL || item.GUID.Value != want.URL || item.GUID.IsPermaLink != "true" {
		t.Fatalf("RSS entry: %s", rssBody)
	}
	pub, err := time.Parse(time.RFC1123Z, item.Published)
	if err != nil || !pub.Equal(published) || item.Updated != "2024-06-07T02:09:10Z" {
		t.Fatalf("RSS times: %s (%v)", rssBody, err)
	}
	atomBody, err := RenderAtom(feed)
	if err != nil {
		t.Fatal(err)
	}
	var atom struct {
		XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
		ID       string   `xml:"id"`
		Title    string   `xml:"title"`
		Subtitle string   `xml:"subtitle"`
		Updated  string   `xml:"updated"`
		Links    []struct {
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Entries []struct {
			ID      string `xml:"id"`
			Title   string `xml:"title"`
			Summary struct {
				Value string `xml:",chardata"`
				Type  string `xml:"type,attr"`
			} `xml:"summary"`
			Author struct {
				Name  string `xml:"name"`
				Email string `xml:"email"`
			} `xml:"author"`
			Published string `xml:"published"`
			Updated   string `xml:"updated"`
			Link      struct {
				Href string `xml:"href,attr"`
				Rel  string `xml:"rel,attr"`
			} `xml:"link"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal(atomBody, &atom); err != nil {
		t.Fatal(err)
	}
	if atom.Title != feed.Title || atom.Subtitle != feed.Description || atom.ID != feed.HomeURL || atom.Updated != "2024-06-07T02:09:10Z" || len(atom.Entries) != 1 {
		t.Fatalf("Atom envelope: %s", atomBody)
	}
	self := false
	for _, link := range atom.Links {
		if link.Rel == "self" && link.Href == feed.AtomURL && link.Type == "application/atom+xml" {
			self = true
		}
	}
	if !self {
		t.Fatal("Atom self link", string(atomBody))
	}
	entry := atom.Entries[0]
	if entry.Title != item.Title || entry.Summary.Value != item.Description || entry.Summary.Type != "text" || entry.Author.Name != item.Creator || entry.Author.Email != "" || entry.ID != item.GUID.Value || entry.Link.Href != item.Link || entry.Link.Rel != "alternate" || entry.Published != "2024-06-07T00:09:10Z" || entry.Updated != item.Updated {
		t.Fatalf("Atom entry: %s", atomBody)
	}
	for _, body := range [][]byte{rssBody, atomBody} {
		if !strings.HasPrefix(string(body), xml.Header) || strings.Contains(string(body), "<content") {
			t.Fatalf("XML declaration/full text: %s", body)
		}
	}
}
