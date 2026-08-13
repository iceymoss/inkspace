package service

import (
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
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

func TestDocDefaultsAndCapabilities(t *testing.T) {
	if got := normalizedDocKind(""); got != models.DocKindMarkdown {
		t.Fatalf("normalizedDocKind(\"\") = %q, want markdown", got)
	}
	if got := normalizedRevision(0); got != 1 {
		t.Fatalf("normalizedRevision(0) = %d, want 1", got)
	}
	if !isCreatableDocKind(models.DocKindMarkdown) || !isCreatableDocKind(models.DocKindText) ||
		!isCreatableDocKind(models.DocKindCode) || isCreatableDocKind(models.DocKindFile) {
		t.Fatal("creatable document kind allowlist is incorrect")
	}

	viewer := docCapabilities(models.WorkspaceRoleViewer, true, models.DocKindMarkdown)
	if !viewer.CanView || viewer.CanEdit || viewer.CanPublish || viewer.CanDelete {
		t.Fatalf("viewer capabilities = %+v", viewer)
	}
	visitor := docCapabilities("", false, models.DocKindMarkdown)
	if !visitor.CanView || visitor.CanEdit || visitor.CanManageMembers {
		t.Fatalf("visitor capabilities = %+v", visitor)
	}
	owner := docCapabilities(models.WorkspaceRoleOwner, true, models.DocKindMarkdown)
	if !owner.CanEdit || !owner.CanPublish || !owner.CanShare || !owner.CanDelete || !owner.CanManageMembers {
		t.Fatalf("owner capabilities = %+v", owner)
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
