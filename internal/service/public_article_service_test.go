package service

import (
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPublicArticleScope(t *testing.T) {
	// MySQL DryRun verifies the shared query contract without connecting to a DB.
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "unused:unused@tcp(127.0.0.1:1)/unused", SkipInitializeWithVersion: true}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	var article models.Article
	query := PublicArticles(db).Where(models.ArticleColumns.Slug+" = ?", "public-post").First(&article)
	sql := query.Statement.SQL.String()
	if !strings.Contains(sql, "status = ?") || !strings.Contains(sql, "deleted_at IS NULL") || !strings.Contains(sql, "slug = ?") {
		t.Fatalf("unsafe public query: %s", sql)
	}
	if len(query.Statement.Vars) < 2 || query.Statement.Vars[0] != 1 || query.Statement.Vars[1] != "public-post" {
		t.Fatalf("vars: %v", query.Statement.Vars)
	}
}
