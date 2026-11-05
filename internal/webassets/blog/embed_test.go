package blog

import (
	"errors"
	"io/fs"
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/seo"
)

func TestEmbeddedBlogHeadPlaceholders(t *testing.T) {
	for _, name := range []string{"fallback/index.html", "dist/index.html"} {
		index, err := fs.ReadFile(embedded, name)
		if err != nil {
			if name == "dist/index.html" && errors.Is(err, fs.ErrNotExist) {
				t.Log("built HTML absent; fallback checked (run frontend build to check dist)")
				continue
			}
			t.Fatal(err)
		}
		result, err := seo.RenderHead(index, seo.SiteIdentity{URL: "https://canonical.example", Name: "站点"}, nil)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if strings.Count(string(result), `rel="alternate"`) != 2 {
			t.Fatal("missing feed alternates", name)
		}
	}
}
