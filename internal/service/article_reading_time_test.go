package service

import (
	"strings"
	"testing"
)

func TestEstimateArticleReadingTime(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "", want: 1},
		{name: "short mixed text", content: "你好, InkSpace 2026!", want: 1},
		{name: "Chinese boundary", content: strings.Repeat("文", 300), want: 1},
		{name: "Chinese above boundary", content: strings.Repeat("文", 301), want: 2},
		{name: "English boundary", content: strings.Repeat("word ", 200), want: 1},
		{name: "English above boundary", content: strings.Repeat("word ", 201), want: 2},
		{name: "mixed rates add exactly", content: strings.Repeat("文", 150) + "\n\n" + strings.Repeat("word ", 100), want: 1},
		{name: "mixed rates round up", content: strings.Repeat("文", 150) + "\n\n" + strings.Repeat("word ", 101), want: 2},
		{name: "mixed rates add minutes", content: strings.Repeat("文", 300) + "\n\n" + strings.Repeat("word ", 200), want: 2},
		{name: "numeric segments", content: strings.Repeat("123 ", 200), want: 1},
		{name: "numeric segments above boundary", content: strings.Repeat("123 ", 201), want: 2},
		{name: "Markdown heading", content: "# " + strings.Repeat("文", 300), want: 1},
		{name: "Markdown quotation and list", content: "> " + strings.Repeat("文", 150) + "\n\n- " + strings.Repeat("文", 151), want: 2},
		{name: "Markdown formatting", content: "**" + strings.Repeat("文", 301) + "**", want: 2},
		{name: "link display text counted", content: "[" + strings.Repeat("word ", 201) + "](https://example.com)", want: 2},
		{name: "link destination excluded", content: strings.Repeat("文", 300) + "\n\n[读者](https://example.com/" + strings.Repeat("word-", 300) + ")", want: 2},
		{name: "image alt counted but URL excluded", content: "![" + strings.Repeat("word ", 200) + "](https://example.com/" + strings.Repeat("word-", 200) + ")", want: 1},
		{name: "formatted image alt text", content: "![**" + strings.Repeat("word ", 201) + "**](image.png)", want: 2},
		{name: "image only", content: "![](https://example.com/image.png)", want: 1},
		{name: "code text included", content: "`" + strings.Repeat("word ", 101) + "`\n\n```go\n" + strings.Repeat("word ", 100) + "\n```", want: 2},
		{name: "HTML preformatted code text included", content: "<pre class=\"code\">\n" + strings.Repeat("word ", 201) + "\n</pre>", want: 2},
		{name: "HTML script text excluded", content: "<script>\n" + strings.Repeat("word ", 201) + "\n</script>", want: 1},
		{name: "fence language excluded", content: "```" + strings.Repeat("language", 201) + "\n```", want: 1},
		{name: "raw HTML markup excluded", content: strings.Repeat("文", 300) + "\n\n<span class=\"" + strings.Repeat("word-", 300) + "\">正文</span>", want: 2},
		{name: "visible text inside HTML block", content: "<div class=\"" + strings.Repeat("word-", 300) + "\">\n" + strings.Repeat("文", 301) + "\n</div>", want: 2},
		{name: "HTML comment content excluded", content: "<!-- " + strings.Repeat("word ", 201) + " -->\n\n正文", want: 1},
		{name: "HTML attribute with greater than sign excluded", content: "<div data-note=\"" + strings.Repeat("word>", 201) + "\">\n" + strings.Repeat("文", 300) + "\n</div>", want: 1},
		{name: "HTML block inline tag text counted", content: "<div><span>" + strings.Repeat("文", 301) + "</span></div>", want: 2},
		{name: "HTML comment inside block excluded", content: "<div>\n<!-- " + strings.Repeat("word ", 201) + " -->\n正文\n</div>", want: 1},
		{name: "autolink URL excluded", content: strings.Repeat("文", 300) + "\n\n<https://example.com/" + strings.Repeat("word-", 300) + ">", want: 1},
		{name: "bare autolink URL excluded", content: strings.Repeat("文", 300) + "\n\nhttps://example.com/" + strings.Repeat("word-", 300), want: 1},
		{name: "autolink email excluded", content: strings.Repeat("文", 300) + "\n\n<person@example.com>", want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EstimateArticleReadingTime(tt.content); got != tt.want {
				t.Errorf("EstimateArticleReadingTime() = %d, want %d", got, tt.want)
			}
		})
	}
}
