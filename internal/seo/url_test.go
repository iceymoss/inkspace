package seo

import "testing"

func TestNormalizePublicSiteURL(t *testing.T) {
	for _, tt := range []struct{ input, want string }{
		{"https://EXAMPLE.com///", "https://example.com"},
		{"http://localhost:8081/", "http://localhost:8081"},
		{"", ""}, {"//example.com", ""}, {"ftp://example.com", ""},
		{"https://u:p@example.com", ""}, {"https://example.com/blog", ""},
		{"https://example.com?x=1", ""}, {"https://example.com#x", ""},
		{"https://example.com?", ""}, {"https://example.com/#", ""},
		{"https://bad host", ""}, {"https://example.com:99999", ""},
		{"https://.", ""}, {"https://-invalid.example", ""}, {"https://example..com", ""}, {"https://example.com:", ""},
		{"https://[example.com]", ""},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got, err := NormalizePublicSiteURL(tt.input)
			if tt.want == "" {
				if err == nil {
					t.Fatalf("accepted %q", tt.input)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
