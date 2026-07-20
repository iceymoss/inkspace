package service

import (
	"bytes"
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/yuin/goldmark"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocService struct{}

func NewDocService() *DocService { return &DocService{} }

func (s *DocService) Create(req *models.DocCreateRequest, ownerID uint) (*models.Doc, error) {
	if err := validateWorkspaceAndCatalog(database.DB, req.WorkspaceID, ownerID, req.CatalogID); err != nil {
		return nil, err
	}
	doc := &models.Doc{
		WorkspaceID: req.WorkspaceID, CatalogID: req.CatalogID, OwnerID: ownerID,
		Title: req.Title, Content: req.Content, WordCount: countWords(req.Content), Sort: req.Sort,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return err
		}
		return tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", req.WorkspaceID, ownerID).
			UpdateColumn("doc_count", gorm.Expr("doc_count + 1")).Error
	})
	if err == nil {
		deleteWorkspaceCache(req.WorkspaceID)
	}
	if err != nil {
		return nil, err
	}
	return s.get(doc.ID, ownerID, database.DB)
}

func (s *DocService) List(workspaceID, ownerID uint, catalogID *uint) ([]*models.Doc, error) {
	if err := requireWorkspace(database.DB, workspaceID, ownerID); err != nil {
		return nil, err
	}
	query := database.DB.Where("workspace_id = ? AND owner_id = ?", workspaceID, ownerID)
	if catalogID != nil {
		if *catalogID == 0 {
			query = query.Where("catalog_id IS NULL")
		} else {
			if err := validateWorkspaceAndCatalog(database.DB, workspaceID, ownerID, catalogID); err != nil {
				return nil, err
			}
			query = query.Where("catalog_id = ?", *catalogID)
		}
	}
	var docs []*models.Doc
	err := query.Order("sort ASC, updated_at DESC").Find(&docs).Error
	return docs, err
}

func (s *DocService) GetEdit(id, ownerID uint) (*models.Doc, error) {
	return s.get(id, ownerID, database.DB)
}

func (s *DocService) Save(id, ownerID uint, req *models.DocSaveRequest) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, ownerID).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		doc.Title = req.Title
		doc.Content = req.Content
		doc.WordCount = countWords(req.Content)
		if err := tx.Model(&doc).Updates(map[string]interface{}{
			"title": doc.Title, "content": doc.Content, "word_count": doc.WordCount,
		}).Error; err != nil {
			return err
		}
		return createDocVersion(tx, &doc, "手动保存")
	})
	if err != nil {
		return nil, err
	}
	return s.get(doc.ID, ownerID, database.DB)
}

func (s *DocService) Autosave(id, ownerID uint, req *models.DocAutosaveRequest) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, ownerID).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if doc.Content == req.Content {
			return nil
		}
		doc.Content = req.Content
		doc.WordCount = countWords(req.Content)
		if err := tx.Model(&doc).Updates(map[string]interface{}{
			"content": doc.Content, "word_count": doc.WordCount,
		}).Error; err != nil {
			return err
		}
		return createDocVersion(tx, &doc, "自动保存")
	})
	if err != nil {
		return nil, err
	}
	return s.get(doc.ID, ownerID, database.DB)
}

func (s *DocService) Publish(id, ownerID uint, status int) (*models.Doc, error) {
	if status != models.DocStatusDraft && status != models.DocStatusPublished {
		return nil, ErrKnowledgeInvalid
	}
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, ownerID).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		updates := map[string]interface{}{"status": status}
		if status == models.DocStatusPublished {
			html, err := renderMarkdown(doc.Content)
			if err != nil {
				return err
			}
			updates["content_html"] = html
		}
		if status == models.DocStatusPublished && doc.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = &now
		}
		return tx.Model(&doc).Updates(updates).Error
	})
	if err != nil {
		return nil, err
	}
	return s.get(id, ownerID, database.DB)
}

func (s *DocService) PublishToBlog(id, ownerID uint, req *models.DocPublishToBlogRequest) (*models.Article, error) {
	doc, err := s.get(id, ownerID, database.DB)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(doc.Title) == "" || strings.TrimSpace(doc.Content) == "" {
		return nil, ErrKnowledgeInvalid
	}
	var categoryCount int64
	if err := database.DB.Model(&models.Category{}).Where("id = ?", req.CategoryID).Count(&categoryCount).Error; err != nil {
		return nil, err
	}
	if categoryCount == 0 {
		return nil, ErrKnowledgeInvalid
	}

	articleReq := &models.ArticleRequest{
		Title: doc.Title, Content: doc.Content, Summary: req.Summary, Cover: req.Cover,
		CategoryID: req.CategoryID, TagIDs: req.TagIDs, Status: 1,
	}
	articleService := NewArticleService()
	if doc.ArticleID != nil {
		var count int64
		if err := database.DB.Model(&models.Article{}).
			Where("id = ? AND author_id = ?", *doc.ArticleID, ownerID).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return articleService.Update(*doc.ArticleID, articleReq, ownerID, "user")
		}
	}

	article, err := articleService.Create(articleReq, ownerID)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Doc{}).Where("id = ? AND owner_id = ?", id, ownerID).
		Update("article_id", article.ID).Error; err != nil {
		_ = articleService.Delete(article.ID, ownerID, "user")
		return nil, err
	}
	return article, nil
}

func (s *DocService) Delete(id, ownerID uint) error {
	var workspaceID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		doc, err := s.get(id, ownerID, tx)
		if err != nil {
			return err
		}
		workspaceID = doc.WorkspaceID
		if err := tx.Where("doc_id = ? AND owner_id = ?", id, ownerID).Delete(&models.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("doc_id = ? AND owner_id = ?", id, ownerID).Delete(&models.DocVersion{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Doc{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", doc.WorkspaceID, ownerID).
			UpdateColumn("doc_count", gorm.Expr("GREATEST(doc_count - 1, 0)")).Error
	})
	if err == nil {
		deleteWorkspaceCache(workspaceID)
	}
	return err
}

func (s *DocService) Move(id, ownerID uint, req *models.DocMoveRequest) (*models.Doc, error) {
	doc, err := s.get(id, ownerID, database.DB)
	if err != nil {
		return nil, err
	}
	if err := validateWorkspaceAndCatalog(database.DB, doc.WorkspaceID, ownerID, req.CatalogID); err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Doc{}).Where("id = ? AND owner_id = ?", id, ownerID).
		Updates(map[string]interface{}{"catalog_id": req.CatalogID, "sort": req.Sort}).Error; err != nil {
		return nil, err
	}
	return s.get(id, ownerID, database.DB)
}

func (s *DocService) Versions(id, ownerID uint) ([]*models.DocVersion, error) {
	if _, err := s.get(id, ownerID, database.DB); err != nil {
		return nil, err
	}
	var versions []*models.DocVersion
	err := database.DB.Where("doc_id = ? AND owner_id = ?", id, ownerID).Order("version DESC").Find(&versions).Error
	return versions, err
}

func (s *DocService) Version(id, ownerID uint, version int) (*models.DocVersion, error) {
	if _, err := s.get(id, ownerID, database.DB); err != nil {
		return nil, err
	}
	var snapshot models.DocVersion
	if err := database.DB.Where("doc_id = ? AND version = ? AND owner_id = ?", id, version, ownerID).First(&snapshot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	return &snapshot, nil
}

func (s *DocService) Rollback(id, ownerID uint, version int) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, ownerID).First(&doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		var snapshot models.DocVersion
		if err := tx.Where("doc_id = ? AND version = ? AND owner_id = ?", id, version, ownerID).First(&snapshot).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		doc.Title = snapshot.Title
		doc.Content = snapshot.Content
		doc.WordCount = countWords(snapshot.Content)
		if err := tx.Model(&doc).Updates(map[string]interface{}{
			"title": doc.Title, "content": doc.Content, "word_count": doc.WordCount,
		}).Error; err != nil {
			return err
		}
		return createDocVersion(tx, &doc, "回滚自 v"+itoa(version))
	})
	if err != nil {
		return nil, err
	}
	return s.get(doc.ID, ownerID, database.DB)
}

func (s *DocService) Search(workspaceID, ownerID uint, keyword string) ([]*models.DocSearchResponse, error) {
	if err := requireWorkspace(database.DB, workspaceID, ownerID); err != nil {
		return nil, err
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*models.DocSearchResponse{}, nil
	}
	var docs []models.Doc
	like := "%" + keyword + "%"
	if err := database.DB.Where("workspace_id = ? AND owner_id = ? AND (title LIKE ? OR content LIKE ?)", workspaceID, ownerID, like, like).
		Order("updated_at DESC").Find(&docs).Error; err != nil {
		return nil, err
	}
	result := make([]*models.DocSearchResponse, 0, len(docs))
	for i := range docs {
		result = append(result, &models.DocSearchResponse{
			ID: docs[i].ID, CatalogID: docs[i].CatalogID, ArticleID: docs[i].ArticleID, Title: docs[i].Title,
			Summary: contentSummary(docs[i].Content, keyword), Status: docs[i].Status,
			WordCount: docs[i].WordCount, UpdatedAt: docs[i].UpdatedAt,
		})
	}
	return result, nil
}

func (s *DocService) get(id, ownerID uint, db *gorm.DB) (*models.Doc, error) {
	var doc models.Doc
	if err := db.Where("id = ? AND owner_id = ?", id, ownerID).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	return &doc, nil
}

func validateWorkspaceAndCatalog(db *gorm.DB, workspaceID, ownerID uint, catalogID *uint) error {
	if err := requireWorkspace(db, workspaceID, ownerID); err != nil {
		return err
	}
	if catalogID == nil {
		return nil
	}
	var count int64
	if err := db.Model(&models.Catalog{}).Where("id = ? AND workspace_id = ? AND owner_id = ?", *catalogID, workspaceID, ownerID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrKnowledgeNotFound
	}
	return nil
}

func createDocVersion(tx *gorm.DB, doc *models.Doc, remark string) error {
	var maxVersion int
	if err := tx.Model(&models.DocVersion{}).Where("doc_id = ?", doc.ID).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion).Error; err != nil {
		return err
	}
	return tx.Create(&models.DocVersion{
		DocID: doc.ID, Version: maxVersion + 1, Title: doc.Title, Content: doc.Content,
		OwnerID: doc.OwnerID, Remark: remark,
	}).Error
}

func renderMarkdown(content string) (string, error) {
	var output bytes.Buffer
	// Goldmark's default renderer escapes raw HTML rather than passing it through.
	if err := goldmark.Convert([]byte(content), &output); err != nil {
		return "", err
	}
	return output.String(), nil
}

func countWords(content string) int {
	count := 0
	inWord := false
	for _, r := range content {
		if unicode.Is(unicode.Han, r) {
			count++
			inWord = false
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				count++
				inWord = true
			}
			continue
		}
		inWord = false
	}
	return count
}

func contentSummary(content, keyword string) string {
	runes := []rune(strings.TrimSpace(content))
	if len(runes) <= 200 {
		return string(runes)
	}
	start := 0
	if index := strings.Index(strings.ToLower(content), strings.ToLower(keyword)); index > 0 {
		start = utf8.RuneCountInString(content[:index]) - 60
		if start < 0 {
			start = 0
		}
	}
	if start+200 > len(runes) {
		start = len(runes) - 200
	}
	return string(runes[start : start+200])
}

// ContentSummary returns a bounded excerpt for document list responses.
func ContentSummary(content string) string {
	return contentSummary(content, "")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
