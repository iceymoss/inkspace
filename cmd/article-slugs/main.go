// article-slugs is an explicit offline migration, never a server startup job.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iceymoss/inkspace/internal/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func run() error {
	if os.Getenv("INKSPACE_ARTICLE_SLUG_MIGRATION_OFFLINE") != "1" {
		return fmt.Errorf("stop all article writers and set INKSPACE_ARTICLE_SLUG_MIGRATION_OFFLINE=1")
	}
	dsn := os.Getenv("INKSPACE_ARTICLE_SLUG_MIGRATION_DSN")
	if dsn == "" {
		return fmt.Errorf("explicit INKSPACE_ARTICLE_SLUG_MIGRATION_DSN required; no developer config fallback")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	stats, err := database.MigrateArticleSlugs(context.Background(), db)
	fmt.Printf("article slug migration: scanned=%d filled=%d\n", stats.Scanned, stats.Filled)
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
