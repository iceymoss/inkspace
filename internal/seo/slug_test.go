package seo

import (
	"strings"
	"testing"
)

func TestSlugBase(t *testing.T) {
	tests := []struct{ title, want string }{
		{"Hello, World!", "hello-world"}, {"中文博客", "zhong-wen-bo-ke"},
		{"Go 123", "go-123"}, {"123", "article-123"}, {"--- / \\ 😀", "article-42"},
		{"", "article-42"}, {"a.b\x00c", "a-b-c"}, {strings.Repeat("a", 200), strings.Repeat("a", 120)},
		{"女绿", "nv-lv"},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := SlugBase(tt.title, 42)
			if got != tt.want || !ValidSlug(got) {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
	for _, invalid := range []string{"123", "UPPER", "a/b", "a\\b", "..", "-a", "a-", "a--b", "", strings.Repeat("a", 121)} {
		if ValidSlug(invalid) {
			t.Errorf("accepted %q", invalid)
		}
	}
}

func TestSlugCandidate(t *testing.T) {
	base := strings.Repeat("x", 119) + "a"
	for _, n := range []uint{1, 2, 99, 100000} {
		got := SlugCandidate(base, n)
		if !ValidSlug(got) {
			t.Fatalf("invalid candidate %q", got)
		}
	}
	if SlugCandidate("same", 1) != "same" || SlugCandidate("same", 2) != "same-2" {
		t.Fatal("unstable suffix")
	}
}
