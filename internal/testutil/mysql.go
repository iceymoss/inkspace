// Package testutil provides opt-in isolated integration infrastructure.
package testutil

import (
	"fmt"
	"os"
	"testing"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SEOTestMySQL creates a fresh randomly named database using an explicitly
// supplied test-server DSN. A selected database is refused to prevent accidental
// use of application/developer schemas. Only this new database is ever dropped.
func SEOTestMySQL(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INKSPACE_SEO_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set INKSPACE_SEO_TEST_MYSQL_DSN to an isolated MySQL server DSN with no database selected")
	}
	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBName != "" {
		t.Fatal("SEO test DSN must not select a database; tests create their own")
	}
	cfg.ParseTime = true
	admin, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	adminSQL, err := admin.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := adminSQL.Close(); err != nil {
			t.Error(err)
		}
	})
	name := fmt.Sprintf("inkspace_seo_test_%d_%d", os.Getpid(), time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + name + "` CHARACTER SET utf8mb4").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + name + "`").Error; err != nil {
			t.Error(err)
		}
	})
	cfg.DBName = name
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent), DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}
