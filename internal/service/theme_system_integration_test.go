package service

import (
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestThemeSystemMySQLIntegration(t *testing.T) {
	dsn := os.Getenv("INKSPACE_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set INKSPACE_TEST_MYSQL_DSN to run MySQL integration tests")
	}
	if os.Getenv("INKSPACE_TEST_ALLOW_DROP") != "1" {
		t.Fatal("set INKSPACE_TEST_ALLOW_DROP=1 to acknowledge destructive test schema reset")
	}
	dsnConfig, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("parse test database DSN: %v", err)
	}
	if !strings.HasSuffix(dsnConfig.DBName, "_test") {
		t.Fatalf("refusing to reset database %q: test database name must end with _test", dsnConfig.DBName)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrator().DropTable(&models.UserAppearance{}, &models.Doc{}, &models.Catalog{}, &models.Workspace{}); err != nil {
		t.Fatalf("reset test schema: %v", err)
	}
	if err := db.AutoMigrate(&models.UserAppearance{}, &models.Workspace{}, &models.Catalog{}, &models.Doc{}); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}

	previousDB := database.DB
	database.DB = db
	t.Cleanup(func() { database.DB = previousDB })

	t.Run("appearance defaults and concurrent upsert", func(t *testing.T) {
		appearanceService := NewUserAppearanceService()
		got, err := appearanceService.Get(42)
		if err != nil {
			t.Fatalf("get default appearance: %v", err)
		}
		if got.UITheme != models.DefaultUITheme || got.ColorScheme != models.DefaultColorScheme {
			t.Fatalf("unexpected defaults: %+v", got)
		}

		var wg sync.WaitGroup
		errorsCh := make(chan error, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				scheme := "light"
				if index%2 == 1 {
					scheme = "dark"
				}
				_, saveErr := appearanceService.Save(42, &models.UserAppearanceRequest{
					UITheme: models.DefaultUITheme, ColorScheme: scheme,
				})
				errorsCh <- saveErr
			}(i)
		}
		wg.Wait()
		close(errorsCh)
		for saveErr := range errorsCh {
			if saveErr != nil {
				t.Fatalf("concurrent save: %v", saveErr)
			}
		}

		var count int64
		if err := db.Model(&models.UserAppearance{}).Where("user_id = ?", 42).Count(&count).Error; err != nil {
			t.Fatalf("count appearance rows: %v", err)
		}
		if count != 1 {
			t.Fatalf("appearance row count = %d, want 1", count)
		}
	})

	t.Run("public wiki filters drafts and private workspaces", func(t *testing.T) {
		privateWorkspace := models.Workspace{OwnerID: 7, Name: "Private"}
		publicWorkspace := models.Workspace{OwnerID: 7, Name: "Public", IsPublic: true}
		if err := db.Create(&privateWorkspace).Error; err != nil {
			t.Fatalf("create private workspace: %v", err)
		}
		if privateWorkspace.IsPublic {
			t.Fatal("workspace without visibility must default to private")
		}
		if err := db.Create(&publicWorkspace).Error; err != nil {
			t.Fatalf("create public workspace: %v", err)
		}

		parent := models.Catalog{WorkspaceID: publicWorkspace.ID, OwnerID: 7, Name: "Parent"}
		if err := db.Create(&parent).Error; err != nil {
			t.Fatalf("create parent catalog: %v", err)
		}
		child := models.Catalog{WorkspaceID: publicWorkspace.ID, OwnerID: 7, ParentID: &parent.ID, Name: "Child"}
		draftOnly := models.Catalog{WorkspaceID: publicWorkspace.ID, OwnerID: 7, Name: "Draft only"}
		if err := db.Create(&child).Error; err != nil {
			t.Fatalf("create child catalog: %v", err)
		}
		if err := db.Create(&draftOnly).Error; err != nil {
			t.Fatalf("create draft catalog: %v", err)
		}

		published := models.Doc{
			WorkspaceID: publicWorkspace.ID, CatalogID: &child.ID, OwnerID: 7,
			Title: "Published", Content: "secret markdown", ContentHTML: "<p>Public</p><script>alert(1)</script>",
			Status: models.DocStatusPublished,
		}
		draft := models.Doc{
			WorkspaceID: publicWorkspace.ID, CatalogID: &draftOnly.ID, OwnerID: 7,
			Title: "Draft", Status: models.DocStatusDraft,
		}
		if err := db.Create(&published).Error; err != nil {
			t.Fatalf("create published doc: %v", err)
		}
		if err := db.Create(&draft).Error; err != nil {
			t.Fatalf("create draft doc: %v", err)
		}

		wikiService := NewPublicWikiService()
		workspaces, total, err := wikiService.Workspaces(1, 20)
		if err != nil {
			t.Fatalf("list public workspaces: %v", err)
		}
		if total != 1 || len(workspaces) != 1 || workspaces[0].DocCount != 1 {
			t.Fatalf("unexpected public workspace list: total=%d list=%+v", total, workspaces)
		}

		tree, err := wikiService.Tree(publicWorkspace.ID)
		if err != nil {
			t.Fatalf("get public tree: %v", err)
		}
		if len(tree.Docs) != 1 || len(tree.Catalogs) != 1 || len(tree.Catalogs[0].Children) != 1 {
			t.Fatalf("unexpected pruned tree: %+v", tree)
		}
		if tree.Catalogs[0].Name != "Parent" || tree.Catalogs[0].Children[0].Name != "Child" {
			t.Fatalf("unexpected catalog ancestry: %+v", tree.Catalogs)
		}

		publicDoc, err := wikiService.Doc(published.ID)
		if err != nil {
			t.Fatalf("get public doc: %v", err)
		}
		if publicDoc.ContentHTML != "<p>Public</p>" {
			t.Fatalf("content was not sanitized: %q", publicDoc.ContentHTML)
		}
		if _, err := wikiService.Doc(draft.ID); !errors.Is(err, ErrKnowledgeNotFound) {
			t.Fatalf("draft lookup error = %v, want not found", err)
		}

		if err := db.Model(&models.Workspace{}).Where("id = ?", publicWorkspace.ID).Update("is_public", false).Error; err != nil {
			t.Fatalf("close public workspace: %v", err)
		}
		if _, err := wikiService.Doc(published.ID); !errors.Is(err, ErrKnowledgeNotFound) {
			t.Fatalf("closed workspace lookup error = %v, want not found", err)
		}
	})
}
