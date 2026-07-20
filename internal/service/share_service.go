package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

type ShareService struct{}

func NewShareService() *ShareService { return &ShareService{} }

func (s *ShareService) Create(docID, ownerID uint, req *models.ShareLinkCreateRequest) (*models.ShareLink, error) {
	var doc models.Doc
	if err := database.DB.Where("id = ? AND owner_id = ?", docID, ownerID).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	expiresAt, err := shareExpiration(req.Permanent, req.ExpiresAt, time.Now())
	if err != nil {
		return nil, err
	}
	token, err := generateShareToken()
	if err != nil {
		return nil, err
	}
	link := &models.ShareLink{Token: token, DocID: doc.ID, OwnerID: ownerID, ExpiresAt: expiresAt, Enabled: true}
	if err := database.DB.Create(link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

func (s *ShareService) List(docID, ownerID uint) ([]*models.ShareLink, error) {
	var count int64
	if err := database.DB.Model(&models.Doc{}).Where("id = ? AND owner_id = ?", docID, ownerID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrKnowledgeNotFound
	}
	var links []*models.ShareLink
	err := database.DB.Where("doc_id = ? AND owner_id = ?", docID, ownerID).Order("id DESC").Find(&links).Error
	return links, err
}

func (s *ShareService) Update(id, ownerID uint, req *models.ShareLinkUpdateRequest) (*models.ShareLink, error) {
	link, err := s.get(id, ownerID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]interface{})
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Permanent || req.ExpiresAt != nil {
		expiresAt, err := shareExpiration(req.Permanent, req.ExpiresAt, time.Now())
		if err != nil {
			return nil, err
		}
		updates["expires_at"] = expiresAt
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&models.ShareLink{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.get(link.ID, ownerID)
}

func (s *ShareService) Delete(id, ownerID uint) error {
	result := database.DB.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.ShareLink{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrKnowledgeNotFound
	}
	return nil
}

func (s *ShareService) Public(token string, now time.Time) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var link models.ShareLink
		if err := tx.Where("token = ?", token).First(&link).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if err := validateShareLink(&link, now); err != nil {
			return err
		}
		if err := tx.Where("id = ?", link.DocID).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		html, err := renderMarkdown(doc.Content)
		if err != nil {
			return err
		}
		doc.ContentHTML = html
		if err := tx.Model(&models.ShareLink{}).Where("id = ?", link.ID).
			UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
			return err
		}
		return tx.Model(&models.Doc{}).Where("id = ?", doc.ID).
			UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
	})
	if err != nil {
		return nil, err
	}
	doc.ViewCount++
	return &doc, nil
}

func (s *ShareService) get(id, ownerID uint) (*models.ShareLink, error) {
	var link models.ShareLink
	if err := database.DB.Where("id = ? AND owner_id = ?", id, ownerID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	return &link, nil
}

func validateShareLink(link *models.ShareLink, now time.Time) error {
	if !link.Enabled {
		return ErrShareDisabled
	}
	if link.ExpiresAt != nil && !now.Before(*link.ExpiresAt) {
		return ErrShareExpired
	}
	return nil
}

func shareExpiration(permanent bool, expiresAt *time.Time, now time.Time) (*time.Time, error) {
	if permanent || expiresAt == nil {
		return nil, nil
	}
	if !expiresAt.After(now) {
		return nil, ErrKnowledgeInvalid
	}
	return expiresAt, nil
}

func generateShareToken() (string, error) {
	data := make([]byte, 18)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
