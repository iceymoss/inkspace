package seo

import (
	"strings"
	"testing"
)

func TestRobotsPublicAndPrivatePaths(t *testing.T) {
	body := string(RenderRobots("https://canonical.example/sitemap.xml"))
	var denied []string
	for _, line := range strings.Split(body, "\n") {
		if value, ok := strings.CutPrefix(line, "Disallow: "); ok {
			denied = append(denied, value)
		}
	}
	if !strings.Contains(body, "User-agent: *\n") || !strings.Contains(body, "Allow: /\n") || !strings.Contains(body, "Sitemap: https://canonical.example/sitemap.xml\n") {
		t.Fatalf("missing discovery rules: %s", body)
	}
	for _, tt := range []struct {
		path    string
		private bool
	}{
		{"/api/articles", true}, {"/api/admin/settings", true},
		{"/dashboard", true}, {"/dashboard/docs/12/edit", true},
		{"/dashboard/workspaces/1", true}, {"/dashboard/articles/create", true},
		{"/favorites", true}, {"/profile/edit", true}, {"/login", true},
		{"/uploads/knowledge/secret.pdf", true},
		{"/", false}, {"/blog/public-post", false}, {"/feed.xml", false},
		{"/atom.xml", false}, {"/sitemap.xml", false}, {"/sitemaps/1.xml", false},
		{"/assets/app.js", false}, {"/uploads/cover.png", false},
		{"/users/123", false}, {"/wiki/docs/123", false},
	} {
		t.Run(tt.path, func(t *testing.T) {
			blocked := false
			for _, prefix := range denied {
				blocked = blocked || strings.HasPrefix(tt.path, prefix)
			}
			if blocked != tt.private {
				t.Fatalf("blocked=%t, want %t: %s", blocked, tt.private, body)
			}
		})
	}
}
