package service

import (
	"context"
	"strings"
	"testing"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type counterSQLLog struct{ sql string }

func (l *counterSQLLog) Printf(_ string, args ...any) {
	for _, arg := range args {
		if s, ok := arg.(string); ok && strings.HasPrefix(s, "UPDATE ") {
			l.sql = s
		}
	}
}

func TestArticleViewStatementsPreserveLegacyAutomaticTimestamp(t *testing.T) {
	// No connection is made: observe the actual public service SQL boundary.
	log := &counterSQLLog{}
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "unused:unused@tcp(127.0.0.1:1)/unused", SkipInitializeWithVersion: true}), &gorm.Config{
		DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true,
		Logger: logger.New(log, logger.Config{LogLevel: logger.Info}),
	})
	if err != nil {
		t.Fatal(err)
	}
	previous := database.DB
	database.DB = db
	t.Cleanup(func() { database.DB = previous })
	for _, record := range []func() error{
		func() error { return NewArticleService().RecordPublicArticleView(context.Background(), 42) },
		func() error { return NewArticleService().IncrementViewCount(42) },
		func() error { id := uint(42); return database.AfterCreateComment(db, &models.Comment{ArticleID: &id}) },
		func() error { id := uint(42); return database.AfterDeleteComment(db, &models.Comment{ArticleID: &id}) },
	} {
		log.sql = ""
		if err := record(); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(log.sql, "`updated_at`=updated_at") {
			t.Fatalf("counter-only SQL must explicitly suppress MySQL ON UPDATE: %s", log.sql)
		}
	}
}
