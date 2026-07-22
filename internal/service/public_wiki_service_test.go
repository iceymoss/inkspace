package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
)

func TestBuildPublicCatalogTreePrunesDraftOnlyBranchesAndKeepsAncestors(t *testing.T) {
	catalogs := []models.Catalog{
		{ID: 1, Name: "root", Sort: 2},
		{ID: 2, ParentID: uintPtr(1), Name: "published", Sort: 1},
		{ID: 3, ParentID: uintPtr(1), Name: "draft only", Sort: 0},
		{ID: 4, Name: "other root", Sort: 1},
	}
	docs := []models.Doc{{ID: 10, CatalogID: uintPtr(2), Status: models.DocStatusPublished}}

	got := buildPublicCatalogTree(catalogs, docs)
	if len(got) != 1 || got[0].ID != 1 || len(got[0].Children) != 1 || got[0].Children[0].ID != 2 {
		t.Fatalf("buildPublicCatalogTree() returned unexpected tree: %#v", got)
	}
	if got[0].Children == nil || got[0].Children[0].Children == nil {
		t.Fatal("public catalog children must serialize as arrays")
	}
}

func TestBuildPublicCatalogTreeOrdersBySortThenID(t *testing.T) {
	catalogs := []models.Catalog{
		{ID: 3, Sort: 1},
		{ID: 2, Sort: 0},
		{ID: 1, Sort: 0},
	}
	docs := []models.Doc{
		{CatalogID: uintPtr(3)},
		{CatalogID: uintPtr(2)},
		{CatalogID: uintPtr(1)},
	}

	tree := buildPublicCatalogTree(catalogs, docs)
	got := []uint{tree[0].ID, tree[1].ID, tree[2].ID}
	if want := []uint{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("catalog order = %v, want %v", got, want)
	}
}

func TestSanitizePublicWikiHTML(t *testing.T) {
	input := `<h2 onclick="bad()">Title</h2><script>alert(1)</script>` +
		`<a href="javascript:alert(2)">bad</a><img src="data:text/html,boom" onerror="bad()">` +
		`<table><tr><th colspan="2">safe</th></tr></table><pre class="language-go"><code>ok</code></pre>`
	got := sanitizePublicWikiHTML(input)

	for _, unsafe := range []string{"onclick", "<script", "javascript:", "data:text/html", "onerror"} {
		if strings.Contains(strings.ToLower(got), unsafe) {
			t.Fatalf("sanitized HTML still contains %q: %s", unsafe, got)
		}
	}
	for _, safe := range []string{"<h2>Title</h2>", "<table>", `colspan="2"`, `class="language-go"`} {
		if !strings.Contains(got, safe) {
			t.Fatalf("sanitized HTML removed %q: %s", safe, got)
		}
	}
}

func TestStatusAfterContentMutation(t *testing.T) {
	if got := statusAfterContentMutation(models.DocStatusPublished); got != models.DocStatusDraft {
		t.Fatalf("published status became %d, want draft", got)
	}
	if got := statusAfterContentMutation(models.DocStatusDraft); got != models.DocStatusDraft {
		t.Fatalf("draft status became %d, want draft", got)
	}
}
