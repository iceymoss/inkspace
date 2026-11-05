package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/seo"
)

// GetPublicSiteIdentity reads only public, non-deleted identity settings. Errors
// are preserved; missing settings use explicit defaults, not private rows.
func (s *SettingService) GetPublicSiteIdentity(ctx context.Context) (seo.SiteIdentity, error) {
	if config.AppConfig == nil {
		return seo.SiteIdentity{}, fmt.Errorf("PUBLIC_SITE_URL is not configured")
	}
	if _, err := seo.NormalizePublicSiteURL(config.AppConfig.PublicSiteURL); err != nil {
		return seo.SiteIdentity{}, err
	}
	var settings []*models.Setting
	if err := database.DB.WithContext(ctx).
		Model(&models.Setting{}).
		Where(models.SettingColumns.IsPublic+" = ?", true).
		Where("`"+models.SettingColumns.Key+"` IN ?", []string{models.SettingSiteName, models.SettingSiteDescription, models.SettingSiteLogo}).
		Find(&settings).Error; err != nil {
		return seo.SiteIdentity{}, err
	}
	return PublicSiteIdentityFromSettings(config.AppConfig.PublicSiteURL, settings)
}

// PublicSiteIdentityFromSettings projects a safe public identity without exposing
// unrelated settings. It validates the origin even in development and tests.
func PublicSiteIdentityFromSettings(base string, settings []*models.Setting) (seo.SiteIdentity, error) {
	normalized, err := seo.NormalizePublicSiteURL(base)
	if err != nil {
		return seo.SiteIdentity{}, err
	}
	site := seo.SiteIdentity{URL: normalized, Name: "InkSpace"}
	for _, setting := range settings {
		if setting == nil || !setting.IsPublic || setting.DeletedAt.Valid {
			continue
		}
		value := strings.TrimSpace(setting.Value)
		switch setting.Key {
		case models.SettingSiteName:
			if value != "" {
				site.Name = value
			}
		case models.SettingSiteDescription:
			site.Description = seo.Summary(value, "", "")
		case models.SettingSiteLogo:
			site.Logo = seo.PublicResourceURL(normalized, value)
		}
	}
	return site, nil
}

// PublicArticleSEOInput explicitly selects public author names, never the full
// account DTO. Author URLs are omitted until a public profile is confirmed.
func PublicArticleSEOInput(article *models.Article) seo.ArticleInput {
	input := seo.ArticleInput{Slug: article.Slug, Title: article.Title, Summary: article.Summary,
		Content: article.Content, Cover: article.Cover, Published: article.CreatedAt, Modified: article.UpdatedAt}
	if article.PublishAt != nil && !article.PublishAt.IsZero() {
		input.Published = *article.PublishAt
	}
	if input.Modified.IsZero() {
		input.Modified = input.Published
	}
	if article.Author != nil {
		input.AuthorName = article.Author.Nickname
		if strings.TrimSpace(input.AuthorName) == "" {
			input.AuthorName = article.Author.Username
		}
	}
	return input
}
