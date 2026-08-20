package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/pkg/uploader"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KnowledgeFileService struct {
	storage uploader.Storage
	config  config.UploadConfig
}

type KnowledgeFileDownload struct {
	FileName     string
	MimeType     string
	Size         int64
	AttachmentID uint
	Reader       io.ReadCloser
}

type KnowledgeFilePreview struct {
	MimeType     string
	Size         int64
	AttachmentID uint
	Reader       io.ReadCloser
}

func NewKnowledgeFileService() (*KnowledgeFileService, error) {
	storage, err := uploader.NewStorage()
	if err != nil {
		return nil, err
	}
	return NewKnowledgeFileServiceWithStorage(storage, config.AppConfig.Upload), nil
}

func NewKnowledgeFileServiceWithStorage(storage uploader.Storage, cfg config.UploadConfig) *KnowledgeFileService {
	return &KnowledgeFileService{storage: storage, config: cfg}
}

func (s *KnowledgeFileService) MaxSize() int64 {
	return knowledgeFileMaxSize(s.config)
}

func (s *KnowledgeFileService) AuthorizeUpload(workspaceID, userID uint) error {
	_, err := authorizeWorkspace(database.DB, workspaceID, userID, WorkspacePermissionEdit)
	return err
}

func (s *KnowledgeFileService) OpenDownload(ctx context.Context, docID, userID uint) (*KnowledgeFileDownload, error) {
	detail, err := NewDocService().Detail(docID, userID)
	if err != nil {
		return nil, err
	}
	if isEditableDocKind(detail.Kind) {
		var doc models.Doc
		if err := database.DB.Select("content").Where("id = ?", docID).First(&doc).Error; err != nil {
			return nil, err
		}
		content := doc.Content
		extension, mimeType := textDownloadFormat(detail.Kind, detail.Language)
		fileName := detail.Title
		if !strings.EqualFold(filepath.Ext(fileName), extension) {
			fileName += extension
		}
		return &KnowledgeFileDownload{
			FileName: fileName, MimeType: mimeType, Size: int64(len(content)),
			Reader: io.NopCloser(strings.NewReader(content)),
		}, nil
	}
	if detail.Attachment == nil {
		return nil, ErrKnowledgeNotFound
	}
	var attachment models.Attachment
	if err := database.DB.Where("id = ?", detail.Attachment.ID).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	if attachment.ObjectKey == "" {
		return nil, ErrKnowledgeNotFound
	}
	reader, err := s.storage.Open(ctx, attachment.ObjectKey)
	if err != nil {
		return nil, err
	}
	return &KnowledgeFileDownload{
		FileName: attachment.FileName, MimeType: attachment.MimeType, Size: attachment.FileSize,
		AttachmentID: attachment.ID, Reader: reader,
	}, nil
}

func (s *KnowledgeFileService) OpenPreview(ctx context.Context, docID, userID uint) (*KnowledgeFilePreview, error) {
	detail, err := NewDocService().Detail(docID, userID)
	if err != nil {
		return nil, err
	}
	if detail.Attachment == nil || detail.Attachment.PreviewStatus != "ready" || !isPreviewSourceMimeType(detail.Attachment.MimeType) {
		return nil, ErrDocPreviewUnavailable
	}
	if isBrowserParsedPreviewMimeType(detail.Attachment.MimeType) && !canBrowserParsePreview(detail.Attachment.MimeType, detail.Attachment.FileSize) {
		return nil, ErrDocPreviewUnavailable
	}
	var attachment models.Attachment
	if err := database.DB.Where("id = ?", detail.Attachment.ID).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	if attachment.ObjectKey == "" {
		return nil, ErrKnowledgeNotFound
	}
	reader, err := s.storage.Open(ctx, attachment.ObjectKey)
	if err != nil {
		return nil, err
	}
	return &KnowledgeFilePreview{
		MimeType: attachment.MimeType, Size: attachment.FileSize, AttachmentID: attachment.ID, Reader: reader,
	}, nil
}

func (s *KnowledgeFileService) Upload(ctx context.Context, workspaceID, userID uint, req *models.KnowledgeFileUploadRequest, header *multipart.FileHeader) (*models.DocDetailResponse, error) {
	if err := s.AuthorizeUpload(workspaceID, userID); err != nil {
		return nil, err
	}
	if header == nil || header.Size <= 0 || header.Size > knowledgeFileMaxSize(s.config) {
		if header != nil && header.Size > knowledgeFileMaxSize(s.config) {
			return nil, ErrKnowledgeFileTooLarge
		}
		return nil, ErrKnowledgeInvalid
	}
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(io.LimitReader(file, knowledgeFileMaxSize(s.config)+1))
	closeErr := file.Close()
	if err != nil {
		return nil, err
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if int64(len(content)) > knowledgeFileMaxSize(s.config) {
		return nil, ErrKnowledgeFileTooLarge
	}
	if len(content) == 0 {
		return nil, ErrKnowledgeInvalid
	}

	fileName := safeKnowledgeFileName(header.Filename)
	policy, err := classifyKnowledgeFile(fileName, content, s.config)
	if err != nil {
		return nil, err
	}
	if !matchesClaimedContentType(header.Header.Get("Content-Type"), policy.mimeType, filepath.Ext(fileName)) {
		return nil, ErrKnowledgeFileType
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	}
	if title == "" || len([]rune(title)) > 200 {
		return nil, ErrKnowledgeInvalid
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	digest := sha256.Sum256(content)
	doc := models.Doc{
		WorkspaceID: workspaceID, CatalogID: req.CatalogID, Title: title, Kind: policy.kind,
		Language: policy.language, Revision: 1,
	}
	if policy.kind == models.DocKindMarkdown || policy.kind == models.DocKindText || policy.kind == models.DocKindCode {
		doc.Content = strings.TrimPrefix(string(content), "\ufeff")
		doc.WordCount = countWords(doc.Content)
		if err := validateStructuredDocContent(doc.Kind, doc.Language, doc.Content); err != nil {
			return nil, err
		}
	}

	var attachment *models.Attachment
	objectKey := ""
	if usesObjectStorage(policy.kind) {
		objectKey = fmt.Sprintf("knowledge/%d/%s/%s%s", workspaceID, time.Now().Format("2006/01/02"), uuid.NewString(), ext)
		if err := s.storage.Put(ctx, objectKey, bytes.NewReader(content), int64(len(content)), policy.mimeType); err != nil {
			return nil, err
		}
		storageType := strings.ToLower(strings.TrimSpace(s.config.StorageType))
		if storageType == "" {
			storageType = "local"
		}
		attachment = &models.Attachment{
			UserID: userID, FileName: fileName, FilePath: objectKey, FileSize: int64(len(content)), FileType: policy.fileType,
			MimeType: policy.mimeType, Extension: ext, StorageType: storageType, ObjectKey: objectKey,
			Checksum: hex.EncodeToString(digest[:]), PreviewStatus: policy.previewStatus,
		}
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		workspace, err := authorizeWorkspaceForUpdate(tx, workspaceID, userID, WorkspacePermissionEdit)
		if err != nil {
			return err
		}
		if req.CatalogID != nil {
			var catalog models.Catalog
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND workspace_id = ? AND owner_id = ?", *req.CatalogID, workspaceID, workspace.OwnerID).
				First(&catalog).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrKnowledgeNotFound
				}
				return err
			}
		}
		if attachment != nil {
			if err := tx.Create(attachment).Error; err != nil {
				return err
			}
			doc.AttachmentID = &attachment.ID
		}
		doc.OwnerID = workspace.OwnerID
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}
		result := tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", workspaceID, workspace.OwnerID).
			UpdateColumn("doc_count", gorm.Expr("doc_count + 1"))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrKnowledgeNotFound
		}
		return nil
	})
	if err != nil {
		if objectKey != "" {
			if deleteErr := s.storage.Delete(context.Background(), objectKey); deleteErr != nil {
				zap.L().Error("failed to remove orphaned knowledge file", zap.String("object_key", objectKey), zap.Error(deleteErr))
			}
		}
		return nil, err
	}
	deleteWorkspaceCache(workspaceID)
	return NewDocService().Detail(doc.ID, userID)
}

func usesObjectStorage(kind string) bool {
	return kind == models.DocKindFile
}

func textDownloadFormat(kind, language string) (string, string) {
	if kind == models.DocKindMarkdown {
		return ".md", "text/markdown; charset=utf-8"
	}
	formats := map[string][2]string{
		"json": {".json", "application/json; charset=utf-8"}, "yaml": {".yaml", "application/yaml; charset=utf-8"},
		"jsonl": {".jsonl", "application/x-ndjson; charset=utf-8"},
		"xml":   {".xml", "application/xml; charset=utf-8"}, "toml": {".toml", "application/toml; charset=utf-8"},
		"csv": {".csv", "text/csv; charset=utf-8"}, "go": {".go", "text/plain; charset=utf-8"},
		"javascript": {".js", "text/plain; charset=utf-8"}, "typescript": {".ts", "text/plain; charset=utf-8"},
		"python": {".py", "text/plain; charset=utf-8"}, "vue": {".vue", "text/plain; charset=utf-8"},
		"ini": {".ini", "text/plain; charset=utf-8"}, "dotenv": {".env", "text/plain; charset=utf-8"},
	}
	if format, ok := formats[language]; ok {
		return format[0], format[1]
	}
	return ".txt", "text/plain; charset=utf-8"
}

func matchesClaimedContentType(claimed, canonical, ext string) bool {
	claimed = strings.ToLower(strings.TrimSpace(strings.Split(claimed, ";")[0]))
	if claimed == "" || claimed == "application/octet-stream" {
		return true
	}
	if claimed == canonical || (canonical == "image/jpeg" && claimed == "image/jpg") {
		return true
	}
	if strings.HasPrefix(canonical, "text/") || canonical == "application/json" || canonical == "application/xml" ||
		canonical == "application/yaml" || canonical == "application/toml" || canonical == "application/x-ndjson" {
		return strings.HasPrefix(claimed, "text/") || claimed == "application/json" || claimed == "application/xml" ||
			claimed == "application/yaml" || claimed == "application/x-yaml" || claimed == "application/toml" || claimed == "application/x-ndjson"
	}
	switch strings.ToLower(ext) {
	case ".docx", ".xlsx", ".pptx", ".odt", ".ods", ".odp":
		return claimed == "application/zip" || claimed == "application/x-zip-compressed"
	case ".rar":
		return claimed == "application/x-rar-compressed"
	case ".gz", ".tgz":
		return claimed == "application/x-gzip"
	}
	return false
}

func safeKnowledgeFileName(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == '/' || r == '\\' {
			return -1
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		return "file"
	}
	runes := []rune(name)
	if len(runes) > 255 {
		name = string(runes[:255])
	}
	return name
}
