package service

import (
	"archive/zip"
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/models"
)

const (
	defaultKnowledgeFileMaxSize  = int64(100 * 1024 * 1024)
	maxKnowledgeFileMaxSize      = int64(100 * 1024 * 1024)
	maxEditableKnowledgeTextSize = 2 * 1024 * 1024
	maxBrowserParsedPreviewSize  = int64(10 * 1024 * 1024)
)

type knowledgeFilePolicy struct {
	kind          string
	language      string
	fileType      string
	mimeType      string
	previewStatus string
}

var knowledgeFilePolicies = map[string]knowledgeFilePolicy{
	".md": {models.DocKindMarkdown, "markdown", "document", "text/markdown", "ready"}, ".markdown": {models.DocKindMarkdown, "markdown", "document", "text/markdown", "ready"},
	".txt": {models.DocKindText, "text", "document", "text/plain", "ready"}, ".log": {models.DocKindText, "text", "document", "text/plain", "ready"},
	".ini": {models.DocKindText, "ini", "document", "text/plain", "ready"}, ".conf": {models.DocKindText, "text", "document", "text/plain", "ready"}, ".env": {models.DocKindText, "dotenv", "document", "text/plain", "ready"},
	".json": {models.DocKindCode, "json", "document", "application/json", "ready"}, ".jsonl": {models.DocKindCode, "jsonl", "document", "application/x-ndjson", "ready"},
	".yaml": {models.DocKindCode, "yaml", "document", "application/yaml", "ready"}, ".yml": {models.DocKindCode, "yaml", "document", "application/yaml", "ready"},
	".xml": {models.DocKindCode, "xml", "document", "application/xml", "ready"}, ".toml": {models.DocKindCode, "toml", "document", "application/toml", "ready"}, ".csv": {models.DocKindCode, "csv", "document", "text/csv", "ready"},
	".pdf": {models.DocKindFile, "", "document", "application/pdf", "ready"},
	".jpg": {models.DocKindFile, "", "image", "image/jpeg", "ready"}, ".jpeg": {models.DocKindFile, "", "image", "image/jpeg", "ready"},
	".png": {models.DocKindFile, "", "image", "image/png", "ready"}, ".gif": {models.DocKindFile, "", "image", "image/gif", "ready"}, ".webp": {models.DocKindFile, "", "image", "image/webp", "ready"},
	".svg": {models.DocKindFile, "", "image", "image/svg+xml", "pending"},
	".doc": {models.DocKindFile, "", "document", "application/msword", "pending"}, ".docx": {models.DocKindFile, "", "document", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "ready"},
	".odt": {models.DocKindFile, "", "document", "application/vnd.oasis.opendocument.text", "pending"}, ".rtf": {models.DocKindFile, "", "document", "application/rtf", "pending"},
	".xls": {models.DocKindFile, "", "document", "application/vnd.ms-excel", "pending"}, ".xlsx": {models.DocKindFile, "", "document", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "pending"},
	".ods": {models.DocKindFile, "", "document", "application/vnd.oasis.opendocument.spreadsheet", "pending"},
	".ppt": {models.DocKindFile, "", "document", "application/vnd.ms-powerpoint", "pending"}, ".pptx": {models.DocKindFile, "", "document", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "pending"},
	".odp": {models.DocKindFile, "", "document", "application/vnd.oasis.opendocument.presentation", "pending"},
	".zip": {models.DocKindFile, "", "archive", "application/zip", "pending"}, ".rar": {models.DocKindFile, "", "archive", "application/vnd.rar", "pending"},
	".7z": {models.DocKindFile, "", "archive", "application/x-7z-compressed", "pending"}, ".tar": {models.DocKindFile, "", "archive", "application/x-tar", "pending"},
	".gz": {models.DocKindFile, "", "archive", "application/gzip", "pending"}, ".tgz": {models.DocKindFile, "", "archive", "application/gzip", "pending"},
}

var codeLanguages = map[string]string{
	".go": "go", ".js": "javascript", ".jsx": "jsx", ".ts": "typescript", ".tsx": "tsx", ".vue": "vue", ".py": "python", ".java": "java", ".kt": "kotlin", ".kts": "kotlin",
	".c": "c", ".h": "c", ".cpp": "cpp", ".hpp": "cpp", ".cs": "csharp", ".rs": "rust", ".php": "php", ".rb": "ruby", ".swift": "swift", ".scala": "scala",
	".sh": "shell", ".bash": "shell", ".zsh": "shell", ".fish": "shell", ".sql": "sql", ".html": "html", ".css": "css", ".scss": "scss", ".less": "less", ".dockerfile": "dockerfile",
}

func knowledgeFileMaxSize(cfg config.UploadConfig) int64 {
	if cfg.KnowledgeMaxSize <= 0 {
		return defaultKnowledgeFileMaxSize
	}
	if cfg.KnowledgeMaxSize > maxKnowledgeFileMaxSize {
		return maxKnowledgeFileMaxSize
	}
	return cfg.KnowledgeMaxSize
}

func isDirectPreviewMimeType(mimeType string) bool {
	switch mimeType {
	case "application/pdf", "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func isBrowserParsedPreviewMimeType(mimeType string) bool {
	return mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
}

func isPreviewSourceMimeType(mimeType string) bool {
	return isDirectPreviewMimeType(mimeType) || isBrowserParsedPreviewMimeType(mimeType)
}

func canBrowserParsePreview(mimeType string, size int64) bool {
	return isBrowserParsedPreviewMimeType(mimeType) && size > 0 && size <= maxBrowserParsedPreviewSize
}

func classifyKnowledgeFile(name string, content []byte, cfg config.UploadConfig) (knowledgeFilePolicy, error) {
	ext := knowledgeFileExtension(name)
	policy, builtIn := knowledgeFilePolicies[ext]
	if language, ok := codeLanguages[ext]; ok {
		policy = knowledgeFilePolicy{models.DocKindCode, language, "document", "text/plain", "ready"}
		builtIn = true
	}
	if !builtIn {
		if !cfg.KnowledgeAllowUnknown && !extensionConfigured(ext, cfg.KnowledgeAllowedExtensions) {
			return knowledgeFilePolicy{}, ErrKnowledgeFileType
		}
		if ext == "" {
			return knowledgeFilePolicy{}, ErrKnowledgeFileType
		}
		return knowledgeFilePolicy{models.DocKindFile, "", "other", "application/octet-stream", "none"}, nil
	}

	if policy.kind == models.DocKindMarkdown || policy.kind == models.DocKindText || policy.kind == models.DocKindCode || ext == ".svg" {
		if bytes.IndexByte(content, 0) >= 0 || !utf8.Valid(content) {
			return knowledgeFilePolicy{}, ErrKnowledgeFileType
		}
		return policy, nil
	}
	if !matchesFileSignature(ext, content) {
		return knowledgeFilePolicy{}, ErrKnowledgeFileType
	}
	return policy, nil
}

func inferCreatableDocType(name string) (string, string, bool) {
	ext := knowledgeFileExtension(name)
	policy, ok := knowledgeFilePolicies[ext]
	if language, code := codeLanguages[ext]; code {
		return models.DocKindCode, language, true
	}
	if !ok || !isCreatableDocKind(policy.kind) {
		return "", "", false
	}
	return policy.kind, policy.language, true
}

func knowledgeFileExtension(name string) string {
	if strings.EqualFold(filepath.Base(name), "dockerfile") {
		return ".dockerfile"
	}
	return strings.ToLower(filepath.Ext(name))
}

func extensionConfigured(ext string, configured []string) bool {
	for _, allowed := range configured {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if allowed != "" && !strings.HasPrefix(allowed, ".") {
			allowed = "." + allowed
		}
		if ext == allowed {
			return true
		}
	}
	return false
}

func matchesFileSignature(ext string, content []byte) bool {
	detected := http.DetectContentType(content)
	switch ext {
	case ".pdf":
		return bytes.HasPrefix(content, []byte("%PDF-"))
	case ".jpg", ".jpeg":
		return detected == "image/jpeg"
	case ".png":
		return detected == "image/png"
	case ".gif":
		return detected == "image/gif"
	case ".webp":
		return len(content) >= 12 && string(content[:4]) == "RIFF" && string(content[8:12]) == "WEBP"
	case ".zip":
		return isZIP(content)
	case ".docx":
		return zipContains(content, "[Content_Types].xml", "word/")
	case ".xlsx":
		return zipContains(content, "[Content_Types].xml", "xl/")
	case ".pptx":
		return zipContains(content, "[Content_Types].xml", "ppt/")
	case ".odt":
		return zipMimetype(content, "application/vnd.oasis.opendocument.text")
	case ".ods":
		return zipMimetype(content, "application/vnd.oasis.opendocument.spreadsheet")
	case ".odp":
		return zipMimetype(content, "application/vnd.oasis.opendocument.presentation")
	case ".rar":
		return bytes.HasPrefix(content, []byte("Rar!\x1a\x07\x00")) || bytes.HasPrefix(content, []byte("Rar!\x1a\x07\x01\x00"))
	case ".7z":
		return bytes.HasPrefix(content, []byte("7z\xbc\xaf'\x1c"))
	case ".gz", ".tgz":
		return len(content) >= 2 && content[0] == 0x1f && content[1] == 0x8b
	case ".tar":
		return len(content) >= 262 && string(content[257:262]) == "ustar"
	case ".doc", ".xls", ".ppt":
		return bytes.HasPrefix(content, []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1})
	case ".rtf":
		return bytes.HasPrefix(bytes.TrimSpace(content), []byte(`{\rtf`))
	default:
		return detected != "application/octet-stream" || mime.TypeByExtension(ext) != ""
	}
}

func isZIP(content []byte) bool {
	return bytes.HasPrefix(content, []byte{'P', 'K', 3, 4}) ||
		bytes.HasPrefix(content, []byte{'P', 'K', 5, 6}) ||
		bytes.HasPrefix(content, []byte{'P', 'K', 7, 8})
}

func zipContains(content []byte, exactName, prefix string) bool {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return false
	}
	foundExact, foundPrefix := false, false
	for _, file := range reader.File {
		name := filepath.ToSlash(file.Name)
		foundExact = foundExact || name == exactName
		foundPrefix = foundPrefix || strings.HasPrefix(name, prefix)
	}
	return foundExact && foundPrefix
}

func zipMimetype(content []byte, expected string) bool {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return false
	}
	for _, file := range reader.File {
		if file.Name != "mimetype" || file.UncompressedSize64 > 200 {
			continue
		}
		entry, err := file.Open()
		if err != nil {
			return false
		}
		value, readErr := io.ReadAll(io.LimitReader(entry, 201))
		closeErr := entry.Close()
		return readErr == nil && closeErr == nil && string(value) == expected
	}
	return false
}
