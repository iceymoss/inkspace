package models

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestArticleSlugInsertAndImmutability(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "unused:unused@tcp(127.0.0.1:1)/unused", SkipInitializeWithVersion: true}), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	article := &Article{ID: 42, Title: "Title", Content: "Body", AuthorID: 1}
	create := db.Create(article)
	if create.Error != nil {
		t.Fatal(create.Error)
	}
	// NULL (rather than an empty string) lets concurrent transactional inserts
	// coexist under the unique index until their final slugs are allocated.
	if strings.Contains(create.Statement.SQL.String(), "`slug`") {
		t.Fatalf("insert: %s", create.Statement.SQL.String())
	}
	update := db.Model(article).Updates(map[string]interface{}{"slug": "changed", "title": "New"})
	if update.Error != nil {
		t.Fatal(update.Error)
	}
	if strings.Contains(update.Statement.SQL.String(), "`slug`") {
		t.Fatalf("mutable slug: %s", update.Statement.SQL.String())
	}
	if article.ToResponse().ID != 42 || article.ToResponse().Slug != article.Slug {
		t.Fatal("response lost permalink identity")
	}
}
