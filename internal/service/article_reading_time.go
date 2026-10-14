package service

import (
	"html"
	"math"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// EstimateArticleReadingTime returns the estimated minutes needed to read Markdown content.
func EstimateArticleReadingTime(content string) int {
	source := []byte(content)
	root := goldmark.New(goldmark.WithExtensions(extension.Linkify)).Parser().Parse(text.NewReader(source))
	var readable []rune

	_ = ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := node.(type) {
		case *ast.AutoLink, *ast.RawHTML:
			return ast.WalkSkipChildren, nil
		case *ast.HTMLBlock:
			if n.HTMLBlockType == ast.HTMLBlockType1 {
				first := n.Lines().At(0)
				if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(string(first.Value(source)))), "<pre") {
					return ast.WalkSkipChildren, nil
				}
			}
			if n.HTMLBlockType == ast.HTMLBlockType2 ||
				n.HTMLBlockType == ast.HTMLBlockType3 || n.HTMLBlockType == ast.HTMLBlockType4 ||
				n.HTMLBlockType == ast.HTMLBlockType5 {
				return ast.WalkSkipChildren, nil
			}
			var block strings.Builder
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				block.Write(segment.Value(source))
			}
			if n.HasClosure() {
				block.Write(n.ClosureLine.Value(source))
			}
			readable = append(readable, ' ')
			readable = append(readable, []rune(html.UnescapeString(stripHTMLTags(block.String())))...)
			readable = append(readable, ' ')
			return ast.WalkSkipChildren, nil
		case *ast.Image:
			readable = append(readable, ' ')
			return ast.WalkContinue, nil
		case *ast.CodeBlock:
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				readable = append(readable, []rune(string(segment.Value(source)))...)
				readable = append(readable, ' ')
			}
			return ast.WalkSkipChildren, nil
		case *ast.FencedCodeBlock:
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				readable = append(readable, []rune(string(segment.Value(source)))...)
				readable = append(readable, ' ')
			}
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			segment := n.Segment
			readable = append(readable, []rune(string(segment.Value(source)))...)
			if n.SoftLineBreak() || n.HardLineBreak() {
				readable = append(readable, ' ')
			}
		case *ast.String:
			readable = append(readable, []rune(string(n.Value))...)
		case *ast.Paragraph, *ast.Heading, *ast.TextBlock:
			if len(readable) > 0 {
				readable = append(readable, ' ')
			}
		}
		return ast.WalkContinue, nil
	})

	han, words, inWord := 0, 0, false
	for _, r := range readable {
		if unicode.Is(unicode.Han, r) {
			han++
			inWord = false
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	minutes := int(math.Ceil(float64(han)/300 + float64(words)/200))
	if minutes < 1 {
		return 1
	}
	return minutes
}

func stripHTMLTags(content string) string {
	var readable strings.Builder
	insideTag := false
	var quote rune
	for _, r := range content {
		switch {
		case insideTag && quote != 0:
			if r == quote {
				quote = 0
			}
		case insideTag && (r == '"' || r == '\''):
			quote = r
		case insideTag && r == '>':
			insideTag = false
			readable.WriteRune(' ')
		case r == '<':
			insideTag = true
			readable.WriteRune(' ')
		case !insideTag:
			readable.WriteRune(r)
		}
	}
	return readable.String()
}
