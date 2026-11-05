package service

import (
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"
)

func TestPublicSiteIdentityProjection(t *testing.T) {
	site, err := PublicSiteIdentityFromSettings("https://blog.example/", []*models.Setting{
		{Key: "site_name", Value: "private name", IsPublic: false},
		{Key: "site_description", Value: "公开描述", IsPublic: true},
		{Key: "site_logo", Value: "/uploads/knowledge/private.png", IsPublic: true},
		{Key: "secret", Value: "private", IsPublic: true},
	})
	if err != nil || site.Name != "InkSpace" || site.Description != "公开描述" || site.Logo != "" || site.URL != "https://blog.example" {
		t.Fatalf("%+v %v", site, err)
	}
	input := PublicArticleSEOInput(&models.Article{Slug: "post", Author: &models.User{Nickname: "公开作者", Email: "secret@example.com", Role: "admin"}})
	m := seo.BuildArticleMetadata(site, input)
	if m.JSONLD["author"].(map[string]any)["name"] != "公开作者" || input.AuthorURL != "" {
		t.Fatal(m)
	}
}
