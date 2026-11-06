package seo

import (
	"strings"
	"testing"
	"time"
)

func TestArticleMetadataPublicProjection(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	m := BuildArticleMetadata(SiteIdentity{URL: "https://blog.example", Name: `站点 & "名称"`, Description: "站点描述", Logo: "/uploads/logo.png"}, ArticleInput{
		Slug: "public-post", Title: `中文 </script><script>alert("x")</script>`,
		Content: "# 正文\n\n<script>secret()</script><style>.private{}</style>\n\n**粗体** [链接](https://secret.example) ![图](private.png)\x00", Cover: "/uploads/knowledge/private.png",
		AuthorName: `作者 & "名字"`, Published: now, Modified: now.Add(time.Hour),
	})
	if m.Description != "正文 粗体 链接" || m.Canonical != "https://blog.example/blog/public-post" || m.Image != "" || m.Title != `中文 </script><script>alert("x")</script> | 站点 & "名称"` {
		t.Fatalf("metadata: %+v", m)
	}
	if m.JSONLD["datePublished"] != "2026-01-02T03:04:05Z" || m.JSONLD["dateModified"] != "2026-01-02T04:04:05Z" {
		t.Fatal(m.JSONLD)
	}
	if _, ok := m.JSONLD["image"]; ok {
		t.Fatal("private cover leaked")
	}
	if got := Summary("<script>secret</script>", "<style>secret</style>", "站点描述"); got != "站点描述" {
		t.Fatal(got)
	}
	if got := Summary("", strings.Repeat("中", 200), ""); len([]rune(got)) != 160 {
		t.Fatal("summary not rune bounded", got)
	}
}

func TestPublicResourceURL(t *testing.T) {
	for _, tt := range []struct{ raw, want string }{
		{"/uploads/cover.png", "https://blog.example/uploads/cover.png"},
		{"https://cdn.example/image.png", "https://cdn.example/image.png"},
		{"", ""}, {"javascript:alert(1)", ""}, {"data:image/png,x", ""}, {"//cdn.example/x", ""},
		{"https://user:pass@cdn.example/x", ""}, {"https://127.0.0.1/x", ""}, {"http://192.168.1.1/x", ""},
		{"https://localhost/x", ""}, {"https://[::1]/x", ""}, {"/uploads/knowledge/secret", ""},
		{"/uploads/%6bnowledge/secret", ""}, {"/uploads/a/../knowledge/secret", ""}, {"/api/private", ""},
		{"https://cdn.example/x?token=secret", ""}, {"https://cdn.example/\nfoo", ""},
		{"https://127.1/x", ""}, {"https://0177.0.0.1/x", ""}, {"https://metadata.internal/x", ""},
		{"/uploads/knowledge%5csecret.png", ""}, {"/uploads/cover%00.png", ""},
	} {
		if got := PublicResourceURL("https://blog.example", tt.raw); got != tt.want {
			t.Errorf("%q: %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestBrowserNumericResourceHostsAreOmitted(t *testing.T) {
	// WHATWG URL hosts accept hexadecimal, octal, mixed-radix, short and
	// single-number IPv4 forms. None may bypass the public resource boundary.
	for _, host := range []string{
		"0x7f.0.0.1", "0X7F.1", "0xa9.0xfe.0xa9.0xfe", "0x7f000001",
		"2130706433", "0177.0.0.1", "017700000001", "127.1", "127.0.1",
		"0x7f.00.0.01", "0x7f.1.", "0x0.0", "0xffffffff", "0x08080808",
		"example.0x7f", "example.123", "example.09",
	} {
		t.Run(host, func(t *testing.T) {
			raw := "http://" + host + "/image.png"
			if got := PublicResourceURL("https://blog.example", raw); got != "" {
				t.Errorf("unsafe browser host accepted: %q", got)
			}
			m := BuildArticleMetadata(SiteIdentity{URL: "https://blog.example", Name: "Site", Logo: raw}, ArticleInput{Slug: "post", Cover: raw, AuthorName: "Author", AuthorURL: raw})
			if m.Image != "" || m.JSONLD["image"] != nil {
				t.Fatal("unsafe cover emitted", m)
			}
			publisher := m.JSONLD["publisher"].(map[string]any)
			author := m.JSONLD["author"].(map[string]any)
			if publisher["logo"] != nil || author["url"] != nil {
				t.Fatal("unsafe logo/author emitted", m.JSONLD)
			}
		})
	}
	for _, raw := range []string{"https://cdn.example/cover.png", "https://0x7f.example/cover.png", "https://8.8.8.8/cover.png"} {
		if got := PublicResourceURL("https://blog.example", raw); got != raw {
			t.Errorf("public resource %q omitted: %q", raw, got)
		}
	}
}

func TestSummaryRemovesMarkupAndUnsafeContents(t *testing.T) {
	for _, tt := range []struct{ source, want string }{
		{"## 中文标题\n\n> **粗体** 和 *斜体*\n\n- 列表 `code`", "中文标题 粗体 和 斜体 列表 code"},
		{"<div>前文<script>alert('秘密')</script><style>秘密</style><iframe>秘密</iframe><svg><text>秘密</text></svg>后文</div>", "前文后文"},
		{"[可见标签](https://private.example/token) ![私有图片说明](secret.png)", "可见标签"},
		{"A\x00\x01\u200bB\nC\tD", "AB C D"},
		{"<script>secret</script>", ""},
	} {
		if got := Summary("", tt.source, ""); got != tt.want {
			t.Errorf("%q: %q want %q", tt.source, got, tt.want)
		}
	}
}
