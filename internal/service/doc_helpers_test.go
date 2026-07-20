package service

import (
	"strings"
	"testing"
)

func TestRenderMarkdownEscapesRawHTML(t *testing.T) {
	html, err := renderMarkdown("# 标题\n\n<script>alert('xss')</script>")
	if err != nil {
		t.Fatalf("renderMarkdown() error = %v", err)
	}
	if strings.Contains(html, "<script>") {
		t.Fatalf("renderMarkdown() returned executable raw HTML: %s", html)
	}
	if !strings.Contains(html, "<h1>") {
		t.Fatalf("renderMarkdown() did not render Markdown heading: %s", html)
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "  \n", want: 0},
		{name: "Chinese", content: "知识库文档", want: 5},
		{name: "English", content: "hello knowledge base", want: 3},
		{name: "mixed", content: "InkSpace 知识 base 2", want: 5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := countWords(test.content); got != test.want {
				t.Fatalf("countWords(%q) = %d, want %d", test.content, got, test.want)
			}
		})
	}
}
