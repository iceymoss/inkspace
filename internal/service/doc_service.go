package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocService struct{}

func NewDocService() *DocService { return &DocService{} }

func (s *DocService) Create(req *models.DocCreateRequest, ownerID uint) (*models.Doc, error) {
	kind := normalizedDocKind(req.Kind)
	language := req.Language
	if req.FileName != "" {
		inferredKind, inferredLanguage, ok := inferCreatableDocType(req.FileName)
		if !ok {
			return nil, ErrKnowledgeFileType
		}
		kind, language = inferredKind, inferredLanguage
	}
	if !isCreatableDocKind(kind) {
		return nil, ErrKnowledgeInvalid
	}
	if len(req.Content) > maxEditableKnowledgeTextSize {
		return nil, ErrKnowledgeTextTooLarge
	}
	if err := validateStructuredDocContent(kind, language, req.Content); err != nil {
		return nil, err
	}
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		workspace, err := authorizeWorkspaceForUpdate(tx, req.WorkspaceID, ownerID, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if req.CatalogID != nil {
			var catalog models.Catalog
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND workspace_id = ? AND owner_id = ?", *req.CatalogID, req.WorkspaceID, workspace.OwnerID).
				First(&catalog).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrKnowledgeNotFound
				}
				return err
			}
		}
		doc = models.Doc{
			WorkspaceID: req.WorkspaceID, CatalogID: req.CatalogID, OwnerID: workspace.OwnerID,
			Title: req.Title, Content: req.Content, Kind: kind, Language: language, Revision: 1,
			WordCount: countWords(req.Content), Sort: req.Sort,
		}
		createResult := tx.Create(&doc)
		if createResult.Error != nil {
			return createResult.Error
		}
		if createResult.RowsAffected != 1 {
			return ErrKnowledgeNotFound
		}
		result := tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", req.WorkspaceID, workspace.OwnerID).
			UpdateColumn("doc_count", gorm.Expr("doc_count + 1"))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrKnowledgeNotFound
		}
		return nil
	})
	if err == nil {
		deleteWorkspaceCache(req.WorkspaceID)
	}
	if err != nil {
		return nil, err
	}
	result, _, err := s.get(doc.ID, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) List(workspaceID, ownerID uint, catalogID *uint) ([]*models.Doc, error) {
	workspace, err := requireWorkspace(database.DB, workspaceID, ownerID, WorkspacePermissionView)
	if err != nil {
		return nil, err
	}
	query := database.DB.Where("workspace_id = ? AND owner_id = ?", workspaceID, workspace.OwnerID)
	if catalogID != nil {
		if *catalogID == 0 {
			query = query.Where("catalog_id IS NULL")
		} else {
			if _, err := validateWorkspaceAndCatalog(database.DB, workspaceID, ownerID, catalogID, WorkspacePermissionView); err != nil {
				return nil, err
			}
			query = query.Where("catalog_id = ?", *catalogID)
		}
	}
	var docs []*models.Doc
	err = query.Order("sort ASC, updated_at DESC").Find(&docs).Error
	return docs, err
}

func (s *DocService) GetEdit(id, ownerID uint) (*models.Doc, error) {
	doc, _, err := s.get(id, ownerID, database.DB, WorkspacePermissionEdit)
	if err == nil && !isOnlineEditableDoc(doc) {
		return nil, ErrDocNotEditable
	}
	return doc, err
}

func (s *DocService) Detail(id, userID uint) (*models.DocDetailResponse, error) {
	var doc models.Doc
	if err := database.DB.Where("id = ?", id).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	var workspace models.Workspace
	if err := database.DB.Where("id = ?", doc.WorkspaceID).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	if doc.OwnerID != workspace.OwnerID {
		return nil, ErrKnowledgeNotFound
	}

	role, member, err := workspaceRole(database.DB, &workspace, userID)
	if err != nil {
		return nil, err
	}
	public := workspace.IsPublic && doc.Status == models.DocStatusPublished
	if !member && !public {
		return nil, ErrKnowledgeNotFound
	}

	kind := normalizedDocKind(doc.Kind)
	response := &models.DocDetailResponse{
		ID: doc.ID, WorkspaceID: doc.WorkspaceID, CatalogID: doc.CatalogID, ArticleID: doc.ArticleID,
		Title: doc.Title, Kind: kind, Language: doc.Language, Revision: normalizedRevision(doc.Revision),
		Status: doc.Status, WordCount: doc.WordCount, ViewCount: doc.ViewCount, PublishedAt: doc.PublishedAt,
		CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt,
		Capabilities: docCapabilities(role, member, kind, len(doc.Content)),
	}
	if kind == models.DocKindMarkdown {
		if member {
			html, err := renderMarkdown(doc.Content)
			if err != nil {
				return nil, err
			}
			response.ContentHTML = sanitizePublicWikiHTML(html)
		} else {
			response.ContentHTML = sanitizePublicWikiHTML(doc.ContentHTML)
		}
	} else if kind == models.DocKindText || kind == models.DocKindCode {
		response.Content = previewTextContent(doc.Content)
	}
	attachmentID := doc.AttachmentID
	if !member && public {
		attachmentID = doc.PublishedAttachmentID
	}
	if attachmentID != nil {
		var attachment models.Attachment
		if err := database.DB.Where("id = ?", *attachmentID).First(&attachment).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			previewStatus := effectiveAttachmentPreviewStatus(&attachment)
			response.Attachment = &models.DocAttachmentDetail{
				ID: attachment.ID, FileName: attachment.FileName, FileSize: attachment.FileSize,
				MimeType: attachment.MimeType, Extension: attachment.Extension, Checksum: attachment.Checksum,
				PreviewStatus: previewStatus, PreviewMimeType: attachment.PreviewMimeType,
			}
			response.Capabilities.CanDownload = true
		}
	}
	return response, nil
}

func effectiveAttachmentPreviewStatus(attachment *models.Attachment) string {
	if attachment != nil && canBrowserParsePreview(attachment.MimeType, attachment.FileSize) && attachment.PreviewStatus == "pending" {
		return "ready"
	}
	if attachment == nil {
		return "none"
	}
	return attachment.PreviewStatus
}

func (s *DocService) Save(id, ownerID uint, req *models.DocSaveRequest) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		_, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND owner_id = ?", id, workspace.OwnerID).First(&doc).Error; err != nil {
			return err
		}
		if !isOnlineEditableDoc(&doc) {
			return ErrDocNotEditable
		}
		if len(req.Content) > maxEditableKnowledgeTextSize {
			return ErrKnowledgeTextTooLarge
		}
		if err := validateStructuredDocContent(normalizedDocKind(doc.Kind), doc.Language, req.Content); err != nil {
			return err
		}
		if normalizedRevision(doc.Revision) != req.Revision {
			return &DocRevisionConflictError{Revision: normalizedRevision(doc.Revision)}
		}
		if doc.Title == req.Title && doc.Content == req.Content {
			doc.Revision = normalizedRevision(doc.Revision)
			return createDocVersion(tx, &doc, "手动保存")
		}
		doc.Title = req.Title
		doc.Content = req.Content
		doc.WordCount = countWords(req.Content)
		doc.Status = statusAfterContentMutation(doc.Status)
		doc.Revision = req.Revision + 1
		result := tx.Model(&models.Doc{}).
			Where("id = ? AND owner_id = ? AND (revision = ? OR (revision = 0 AND ? = 1))", id, workspace.OwnerID, req.Revision, req.Revision).
			Updates(map[string]interface{}{
				"title": doc.Title, "content": doc.Content, "word_count": doc.WordCount,
				"status": doc.Status, "revision": doc.Revision,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return docRevisionConflict(tx, id)
		}
		return createDocVersion(tx, &doc, "手动保存")
	})
	if err != nil {
		return nil, err
	}
	result, _, err := s.get(doc.ID, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) Autosave(id, ownerID uint, req *models.DocAutosaveRequest) (*models.Doc, error) {
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		_, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND owner_id = ?", id, workspace.OwnerID).First(&doc).Error; err != nil {
			return err
		}
		if !isOnlineEditableDoc(&doc) {
			return ErrDocNotEditable
		}
		if len(req.Content) > maxEditableKnowledgeTextSize {
			return ErrKnowledgeTextTooLarge
		}
		if err := validateStructuredDocContent(normalizedDocKind(doc.Kind), doc.Language, req.Content); err != nil {
			return err
		}
		if normalizedRevision(doc.Revision) != req.Revision {
			return &DocRevisionConflictError{Revision: normalizedRevision(doc.Revision)}
		}
		if doc.Content == req.Content {
			doc.Revision = normalizedRevision(doc.Revision)
			return nil
		}
		doc.Content = req.Content
		doc.WordCount = countWords(req.Content)
		doc.Status = statusAfterContentMutation(doc.Status)
		doc.Revision = req.Revision + 1
		result := tx.Model(&models.Doc{}).
			Where("id = ? AND owner_id = ? AND (revision = ? OR (revision = 0 AND ? = 1))", id, workspace.OwnerID, req.Revision, req.Revision).
			Updates(map[string]interface{}{
				"content": doc.Content, "word_count": doc.WordCount, "status": doc.Status, "revision": doc.Revision,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return docRevisionConflict(tx, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	result, _, err := s.get(doc.ID, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) Publish(id, ownerID uint, status int) (*models.Doc, error) {
	if status != models.DocStatusDraft && status != models.DocStatusPublished {
		return nil, ErrKnowledgeInvalid
	}
	var doc models.Doc
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		_, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, workspace.OwnerID).First(&doc).Error; err != nil {
			return err
		}
		updates := map[string]interface{}{"status": status}
		if status == models.DocStatusPublished {
			if normalizedDocKind(doc.Kind) == models.DocKindMarkdown {
				html, err := renderMarkdown(doc.Content)
				if err != nil {
					return err
				}
				updates["content_html"] = sanitizePublicWikiHTML(html)
			}
			updates["published_attachment_id"] = doc.AttachmentID
			updates["published_revision"] = normalizedRevision(doc.Revision)
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
	result, _, err := s.get(id, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) PublishToBlog(id, ownerID uint, req *models.DocPublishToBlogRequest) (*models.Article, error) {
	var article *models.Article
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		doc, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND owner_id = ?", id, workspace.OwnerID).First(doc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if strings.TrimSpace(doc.Title) == "" || strings.TrimSpace(doc.Content) == "" {
			return ErrKnowledgeInvalid
		}
		var categoryCount int64
		if err := tx.Model(&models.Category{}).Where("id = ?", req.CategoryID).Count(&categoryCount).Error; err != nil {
			return err
		}
		if categoryCount == 0 {
			return ErrKnowledgeInvalid
		}

		articleReq := &models.ArticleRequest{
			Title: doc.Title, Content: doc.Content, Summary: req.Summary, Cover: req.Cover,
			CategoryID: req.CategoryID, TagIDs: req.TagIDs, Status: 1,
		}
		articleService := NewArticleService()
		if doc.ArticleID != nil {
			var count int64
			if err := tx.Model(&models.Article{}).
				Where("id = ? AND author_id = ?", *doc.ArticleID, workspace.OwnerID).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				article, err = articleService.Update(*doc.ArticleID, articleReq, workspace.OwnerID, "user")
				return err
			}
		}

		article, err = articleService.Create(articleReq, workspace.OwnerID)
		if err != nil {
			return err
		}
		result := tx.Model(&models.Doc{}).Where("id = ? AND owner_id = ?", id, workspace.OwnerID).
			Update("article_id", article.ID)
		if result.Error != nil {
			_ = articleService.Delete(article.ID, workspace.OwnerID, "user")
			return result.Error
		}
		if result.RowsAffected != 1 {
			_ = articleService.Delete(article.ID, workspace.OwnerID, "user")
			return ErrKnowledgeNotFound
		}
		return nil
	})
	return article, err
}

func (s *DocService) Delete(id, ownerID uint) error {
	var workspaceID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		doc, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		workspaceID = doc.WorkspaceID
		if err := tx.Where("doc_id = ? AND owner_id = ?", id, workspace.OwnerID).Delete(&models.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("doc_id = ? AND owner_id = ?", id, workspace.OwnerID).Delete(&models.DocVersion{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND owner_id = ?", id, workspace.OwnerID).Delete(&models.Doc{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", doc.WorkspaceID, workspace.OwnerID).
			UpdateColumn("doc_count", gorm.Expr("GREATEST(doc_count - 1, 0)")).Error
	})
	if err == nil {
		deleteWorkspaceCache(workspaceID)
	}
	return err
}

func (s *DocService) Move(id, ownerID uint, req *models.DocMoveRequest) (*models.Doc, error) {
	doc, workspace, err := s.get(id, ownerID, database.DB, WorkspacePermissionEdit)
	if err != nil {
		return nil, err
	}
	if _, err := validateWorkspaceAndCatalog(database.DB, doc.WorkspaceID, ownerID, req.CatalogID, WorkspacePermissionEdit); err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Doc{}).Where("id = ? AND owner_id = ?", id, workspace.OwnerID).
		Updates(map[string]interface{}{"catalog_id": req.CatalogID, "sort": req.Sort}).Error; err != nil {
		return nil, err
	}
	result, _, err := s.get(id, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) Versions(id, ownerID uint) ([]*models.DocVersion, error) {
	doc, workspace, err := s.get(id, ownerID, database.DB, WorkspacePermissionView)
	if err != nil {
		return nil, err
	}
	var versions []*models.DocVersion
	err = database.DB.Where("doc_id = ? AND owner_id = ?", doc.ID, workspace.OwnerID).Order("version DESC").Find(&versions).Error
	return versions, err
}

func (s *DocService) Version(id, ownerID uint, version int) (*models.DocVersion, error) {
	_, workspace, err := s.get(id, ownerID, database.DB, WorkspacePermissionView)
	if err != nil {
		return nil, err
	}
	var snapshot models.DocVersion
	if err := database.DB.Where("doc_id = ? AND version = ? AND owner_id = ?", id, version, workspace.OwnerID).First(&snapshot).Error; err != nil {
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
		_, workspace, err := s.get(id, ownerID, tx, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND owner_id = ?", id, workspace.OwnerID).First(&doc).Error; err != nil {
			return err
		}
		var snapshot models.DocVersion
		if err := tx.Where("doc_id = ? AND version = ? AND owner_id = ?", id, version, workspace.OwnerID).First(&snapshot).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		doc.Title = snapshot.Title
		doc.Content = snapshot.Content
		doc.Kind = normalizedDocKind(snapshot.Kind)
		doc.Language = snapshot.Language
		doc.AttachmentID = snapshot.AttachmentID
		doc.WordCount = countWords(snapshot.Content)
		doc.Status = statusAfterContentMutation(doc.Status)
		doc.Revision = normalizedRevision(doc.Revision) + 1
		if err := tx.Model(&doc).Updates(map[string]interface{}{
			"title": doc.Title, "content": doc.Content, "kind": doc.Kind, "language": doc.Language,
			"attachment_id": doc.AttachmentID, "word_count": doc.WordCount, "status": doc.Status, "revision": doc.Revision,
		}).Error; err != nil {
			return err
		}
		return createDocVersion(tx, &doc, "回滚自 v"+itoa(version))
	})
	if err != nil {
		return nil, err
	}
	result, _, err := s.get(doc.ID, ownerID, database.DB, WorkspacePermissionEdit)
	return result, err
}

func (s *DocService) Search(workspaceID, ownerID uint, keyword string) ([]*models.DocSearchResponse, error) {
	workspace, err := requireWorkspace(database.DB, workspaceID, ownerID, WorkspacePermissionView)
	if err != nil {
		return nil, err
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*models.DocSearchResponse{}, nil
	}
	var docs []models.Doc
	like := "%" + keyword + "%"
	if err := database.DB.Where("workspace_id = ? AND owner_id = ? AND (title LIKE ? OR content LIKE ?)", workspaceID, workspace.OwnerID, like, like).
		Order("updated_at DESC").Find(&docs).Error; err != nil {
		return nil, err
	}
	result := make([]*models.DocSearchResponse, 0, len(docs))
	for i := range docs {
		result = append(result, &models.DocSearchResponse{
			ID: docs[i].ID, CatalogID: docs[i].CatalogID, ArticleID: docs[i].ArticleID, Title: docs[i].Title,
			Kind: normalizedDocKind(docs[i].Kind), Language: docs[i].Language, Editable: isOnlineEditableDoc(&docs[i]),
			Summary: contentSummary(docs[i].Content, keyword), Status: docs[i].Status,
			WordCount: docs[i].WordCount, UpdatedAt: docs[i].UpdatedAt,
		})
	}
	return result, nil
}

func (s *DocService) get(id, ownerID uint, db *gorm.DB, permission WorkspacePermission) (*models.Doc, *models.Workspace, error) {
	return workspaceForDoc(db, id, ownerID, permission)
}

func validateWorkspaceAndCatalog(db *gorm.DB, workspaceID, ownerID uint, catalogID *uint, permission WorkspacePermission) (*models.Workspace, error) {
	workspace, err := requireWorkspace(db, workspaceID, ownerID, permission)
	if err != nil {
		return nil, err
	}
	if catalogID == nil {
		return workspace, nil
	}
	var count int64
	if err := db.Model(&models.Catalog{}).Where("id = ? AND workspace_id = ? AND owner_id = ?", *catalogID, workspaceID, workspace.OwnerID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrKnowledgeNotFound
	}
	return workspace, nil
}

func createDocVersion(tx *gorm.DB, doc *models.Doc, remark string) error {
	var maxVersion int
	if err := tx.Model(&models.DocVersion{}).Where("doc_id = ?", doc.ID).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion).Error; err != nil {
		return err
	}
	return tx.Create(&models.DocVersion{
		DocID: doc.ID, Version: maxVersion + 1, Title: doc.Title, Content: doc.Content,
		Revision: normalizedRevision(doc.Revision), Kind: normalizedDocKind(doc.Kind), Language: doc.Language,
		AttachmentID: doc.AttachmentID, OwnerID: doc.OwnerID, Remark: remark,
	}).Error
}

func workspaceRole(db *gorm.DB, workspace *models.Workspace, userID uint) (string, bool, error) {
	if userID == 0 {
		return "", false, nil
	}
	if workspace.OwnerID == userID {
		return models.WorkspaceRoleOwner, true, nil
	}
	var member models.WorkspaceMember
	if err := db.Where("workspace_id = ? AND user_id = ?", workspace.ID, userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, nil
		}
		return "", false, err
	}
	if !HasWorkspacePermission(member.Role, WorkspacePermissionView) {
		return "", false, nil
	}
	return member.Role, true, nil
}

func docCapabilities(role string, member bool, kind string, contentSize int) models.DocCapabilities {
	canEdit := member && HasWorkspacePermission(role, WorkspacePermissionEdit)
	canManage := member && HasWorkspacePermission(role, WorkspacePermissionManageMembers)
	onlineEditable := isEditableDocKind(kind) && contentSize <= maxEditableKnowledgeTextSize
	return models.DocCapabilities{
		CanView: true, CanEdit: canEdit && onlineEditable, CanReplace: canEdit && kind == models.DocKindFile,
		CanDownload: true, CanPublish: canEdit, CanShare: canEdit,
		CanDelete: canEdit, CanManageMembers: canManage,
	}
}

func normalizedDocKind(kind string) string {
	if kind == "" {
		return models.DocKindMarkdown
	}
	return kind
}

func normalizedRevision(revision uint64) uint64 {
	if revision == 0 {
		return 1
	}
	return revision
}

func isCreatableDocKind(kind string) bool {
	return kind == models.DocKindMarkdown || kind == models.DocKindText || kind == models.DocKindCode
}

func isEditableDocKind(kind string) bool { return isCreatableDocKind(kind) }

func isOnlineEditableDoc(doc *models.Doc) bool {
	return doc != nil && isEditableDocKind(normalizedDocKind(doc.Kind)) && len(doc.Content) <= maxEditableKnowledgeTextSize
}

func IsDocOnlineEditable(doc *models.Doc) bool {
	return isOnlineEditableDoc(doc)
}

func previewTextContent(content string) string {
	if len(content) <= maxEditableKnowledgeTextSize {
		return content
	}
	end := maxEditableKnowledgeTextSize
	for end > 0 && !utf8.RuneStart(content[end]) {
		end--
	}
	return content[:end]
}

func validateStructuredDocContent(kind, language, content string) error {
	if kind != models.DocKindCode || strings.TrimSpace(content) == "" {
		return nil
	}
	switch strings.ToLower(language) {
	case "json":
		if !json.Valid([]byte(content)) {
			return fmt.Errorf("%w: JSON 语法无效", ErrKnowledgeInvalid)
		}
	case "yaml":
		var value interface{}
		if err := yaml.Unmarshal([]byte(content), &value); err != nil {
			return fmt.Errorf("%w: YAML 语法无效: %v", ErrKnowledgeInvalid, err)
		}
	}
	return nil
}

func docRevisionConflict(tx *gorm.DB, id uint) error {
	var current models.Doc
	if err := tx.Select("revision").Where("id = ?", id).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrKnowledgeNotFound
		}
		return err
	}
	return &DocRevisionConflictError{Revision: normalizedRevision(current.Revision)}
}

func renderMarkdown(content string) (string, error) {
	var output bytes.Buffer
	// Goldmark's default renderer escapes raw HTML rather than passing it through.
	if err := goldmark.Convert([]byte(content), &output); err != nil {
		return "", err
	}
	return output.String(), nil
}

func statusAfterContentMutation(status int) int {
	if status == models.DocStatusPublished {
		return models.DocStatusDraft
	}
	return status
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
