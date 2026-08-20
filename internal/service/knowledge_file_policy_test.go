package service

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/models"
)

func TestClassifyKnowledgeFile(t *testing.T) {
	zipHeader := []byte{'P', 'K', 3, 4}
	docx := testZIP(t, map[string]string{"[Content_Types].xml": "types", "word/document.xml": "document"})
	tests := []struct {
		name     string
		fileName string
		content  []byte
		wantKind string
		wantErr  error
	}{
		{name: "markdown", fileName: "readme.md", content: []byte("# hello"), wantKind: models.DocKindMarkdown},
		{name: "code", fileName: "main.go", content: []byte("package main"), wantKind: models.DocKindCode},
		{name: "pdf", fileName: "paper.pdf", content: []byte("%PDF-1.7\n"), wantKind: models.DocKindFile},
		{name: "binary disguised as text", fileName: "note.txt", content: []byte{'a', 0, 'b'}, wantErr: ErrKnowledgeFileType},
		{name: "signature mismatch", fileName: "paper.pdf", content: []byte("plain text"), wantErr: ErrKnowledgeFileType},
		{name: "zip container", fileName: "report.zip", content: zipHeader, wantKind: models.DocKindFile},
		{name: "docx container", fileName: "report.docx", content: docx, wantKind: models.DocKindFile},
		{name: "zip disguised as docx", fileName: "report.docx", content: testZIP(t, map[string]string{"data.txt": "data"}), wantErr: ErrKnowledgeFileType},
		{name: "unknown", fileName: "data.bin", content: []byte{1, 2, 3}, wantErr: ErrKnowledgeFileType},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy, err := classifyKnowledgeFile(test.fileName, test.content, config.UploadConfig{})
			if err != test.wantErr {
				t.Fatalf("classifyKnowledgeFile() error = %v, want %v", err, test.wantErr)
			}
			if err == nil && policy.kind != test.wantKind {
				t.Fatalf("kind = %q, want %q", policy.kind, test.wantKind)
			}
		})
	}
}

func TestMatchesClaimedContentType(t *testing.T) {
	if !matchesClaimedContentType("application/octet-stream", "application/pdf", ".pdf") {
		t.Fatal("generic content type should be accepted after signature validation")
	}
	if matchesClaimedContentType("image/png", "application/pdf", ".pdf") {
		t.Fatal("mismatched content type should be rejected")
	}
	if matchesClaimedContentType("image/png", "text/plain", ".txt") {
		t.Fatal("binary content type should be rejected for text")
	}
}

func testZIP(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create zip entry: %v", err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buffer.Bytes()
}

func TestKnowledgeFileMaxSize(t *testing.T) {
	if got := knowledgeFileMaxSize(config.UploadConfig{}); got != defaultKnowledgeFileMaxSize {
		t.Fatalf("default max size = %d, want %d", got, defaultKnowledgeFileMaxSize)
	}
	if got := knowledgeFileMaxSize(config.UploadConfig{KnowledgeMaxSize: maxKnowledgeFileMaxSize + 1}); got != maxKnowledgeFileMaxSize {
		t.Fatalf("capped max size = %d, want %d", got, maxKnowledgeFileMaxSize)
	}
}

func TestDirectPreviewMimeType(t *testing.T) {
	for _, mimeType := range []string{"application/pdf", "image/jpeg", "image/png", "image/gif", "image/webp"} {
		if !isDirectPreviewMimeType(mimeType) {
			t.Errorf("%s should support direct preview", mimeType)
		}
	}
	for _, mimeType := range []string{"image/svg+xml", "application/zip", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"} {
		if isDirectPreviewMimeType(mimeType) {
			t.Errorf("%s must not use raw direct preview", mimeType)
		}
	}
}

func TestBrowserParsedPreviewMimeType(t *testing.T) {
	docx := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if !isBrowserParsedPreviewMimeType(docx) || !isPreviewSourceMimeType(docx) {
		t.Fatal("DOCX should be available to the controlled browser parser")
	}
	if isBrowserParsedPreviewMimeType("application/vnd.ms-powerpoint") {
		t.Fatal("legacy Office formats must use server-side conversion")
	}
	if !canBrowserParsePreview(docx, 1024) || canBrowserParsePreview(docx, maxBrowserParsedPreviewSize+1) {
		t.Fatal("DOCX browser parsing size limit is incorrect")
	}
}

func TestConfiguredKnowledgeExtension(t *testing.T) {
	policy, err := classifyKnowledgeFile("design.psd", []byte{1, 2, 3}, config.UploadConfig{KnowledgeAllowedExtensions: []string{"PSD"}})
	if err != nil {
		t.Fatalf("classify configured extension: %v", err)
	}
	if policy.kind != models.DocKindFile || policy.previewStatus != "none" {
		t.Fatalf("policy = %+v, want fallback file", policy)
	}
}

func TestSafeKnowledgeFileName(t *testing.T) {
	if got := safeKnowledgeFileName(`../unsafe\report.pdf`); got != "report.pdf" {
		t.Fatalf("safeKnowledgeFileName() = %q, want report.pdf", got)
	}
}

func TestInferCreatableDocType(t *testing.T) {
	tests := []struct {
		name         string
		wantKind     string
		wantLanguage string
		wantOK       bool
	}{
		{name: "README.md", wantKind: models.DocKindMarkdown, wantLanguage: "markdown", wantOK: true},
		{name: "notes.txt", wantKind: models.DocKindText, wantLanguage: "text", wantOK: true},
		{name: "styles.css", wantKind: models.DocKindCode, wantLanguage: "css", wantOK: true},
		{name: "app.js", wantKind: models.DocKindCode, wantLanguage: "javascript", wantOK: true},
		{name: "worker.rs", wantKind: models.DocKindCode, wantLanguage: "rust", wantOK: true},
		{name: "config.yaml", wantKind: models.DocKindCode, wantLanguage: "yaml", wantOK: true},
		{name: "Dockerfile", wantKind: models.DocKindCode, wantLanguage: "dockerfile", wantOK: true},
		{name: "manual.pdf", wantOK: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kind, language, ok := inferCreatableDocType(test.name)
			if kind != test.wantKind || language != test.wantLanguage || ok != test.wantOK {
				t.Fatalf("inferCreatableDocType() = %q, %q, %v", kind, language, ok)
			}
		})
	}
}

func TestKnowledgeFileStorageStrategy(t *testing.T) {
	for _, kind := range []string{models.DocKindMarkdown, models.DocKindText, models.DocKindCode} {
		if usesObjectStorage(kind) {
			t.Fatalf("kind %q should use Doc.Content instead of object storage", kind)
		}
	}
	if !usesObjectStorage(models.DocKindFile) {
		t.Fatal("binary file should use object storage")
	}
}

func TestTextDownloadFormat(t *testing.T) {
	extension, mimeType := textDownloadFormat(models.DocKindCode, "json")
	if extension != ".json" || mimeType != "application/json; charset=utf-8" {
		t.Fatalf("JSON format = %q %q", extension, mimeType)
	}
}
