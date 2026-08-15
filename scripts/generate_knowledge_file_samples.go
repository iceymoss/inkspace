//go:build ignore

package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const outputDir = "docs/samples/knowledge-files"

func main() {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	files := map[string][]byte{
		"README.md": []byte(`# 知识库文件上传样例

本目录由 go run scripts/generate_knowledge_file_samples.go 生成，用于知识库文件上传、类型识别和预览联调。

| 文件 | 预期 kind | 预期状态 |
|------|-----------|----------|
| sample.md | markdown | ready，可在线编辑 |
| sample.txt | text | ready，可在线编辑 |
| sample.go | code | ready，可在线编辑 |
| sample.json | code/json | ready，可在线编辑 |
| sample.yaml | code/yaml | ready，可在线编辑 |
| sample.csv | code/csv | ready，可在线编辑 |
| sample.pdf | file | ready |
| sample.png | file/image | ready |
| sample.svg | file/image | pending，等待安全清洗 |
| sample.zip | file/archive | pending |
| sample.gz | file/archive | pending |
| sample.docx | file/office | pending |
| sample.xlsx | file/office | pending |
| sample.pptx | file/office | pending |
| sample.odt | file/office | pending |
| sample.bin | fallback file | 默认拒绝；配置允许未知扩展名后接收 |
| invalid-disguised.pdf | - | 应拒绝，内容不是 PDF |
| invalid-binary.txt | - | 应拒绝，包含 NUL 字节 |

Office 文件是满足当前上传策略检测的最小容器，不用于验证 LibreOffice 的视觉转换质量。
`),
		"sample.md":             []byte("# InkSpace 知识库样例\n\n这是 **Markdown** 文档。\n\n- 只读预览\n- 显式编辑\n- revision 冲突检测\n\n```go\nfmt.Println(\"hello InkSpace\")\n```\n"),
		"sample.txt":            []byte("InkSpace knowledge file sample\n\nThis is a UTF-8 plain text document.\n第二行包含中文，可用于验证编码和等宽预览。\n"),
		"sample.go":             []byte("package sample\n\nimport \"fmt\"\n\nfunc Hello(name string) string {\n\treturn fmt.Sprintf(\"Hello, %s\", name)\n}\n"),
		"sample.json":           []byte("{\n  \"name\": \"InkSpace\",\n  \"enabled\": true,\n  \"features\": [\"preview\", \"revision\", \"download\"]\n}\n"),
		"sample.yaml":           []byte("name: InkSpace\nenabled: true\nfeatures:\n  - preview\n  - revision\n  - download\n"),
		"sample.csv":            []byte("id,title,status\n1,Knowledge Base,ready\n2,Document Preview,pending\n"),
		"sample.pdf":            minimalPDF(),
		"sample.png":            mustBase64("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="),
		"sample.svg":            []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="640" height="240" viewBox="0 0 640 240"><rect width="640" height="240" fill="#17202a"/><circle cx="100" cy="120" r="54" fill="#57c7ff"/><text x="180" y="110" fill="#fff" font-family="sans-serif" font-size="36">InkSpace</text><text x="180" y="150" fill="#b8c2cc" font-family="sans-serif" font-size="20">Knowledge file preview</text></svg>`),
		"sample.bin":            []byte{0x49, 0x4e, 0x4b, 0x53, 0x50, 0x41, 0x43, 0x45},
		"invalid-disguised.pdf": []byte("This plain text file only has a PDF extension.\n"),
		"invalid-binary.txt":    []byte{'t', 'e', 'x', 't', 0, 'b', 'i', 'n', 'a', 'r', 'y'},
	}

	files["sample.zip"] = zipBytes(map[string][]byte{
		"docs/readme.txt": []byte("Archive manifest sample\n"),
		"src/main.go":     files["sample.go"],
	})
	files["sample.gz"] = gzipBytes([]byte("InkSpace gzip archive sample\n"))
	files["sample.docx"] = zipBytes(map[string][]byte{
		"[Content_Types].xml": []byte(`<?xml version="1.0"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`),
		"word/document.xml":   []byte(`<?xml version="1.0"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>InkSpace DOCX sample</w:t></w:r></w:p></w:body></w:document>`),
	})
	files["sample.xlsx"] = zipBytes(map[string][]byte{
		"[Content_Types].xml": []byte(`<?xml version="1.0"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/></Types>`),
		"xl/workbook.xml":     []byte(`<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets/></workbook>`),
	})
	files["sample.pptx"] = zipBytes(map[string][]byte{
		"[Content_Types].xml":  []byte(`<?xml version="1.0"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/></Types>`),
		"ppt/presentation.xml": []byte(`<?xml version="1.0"?><p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst/></p:presentation>`),
	})
	files["sample.odt"] = zipBytes(map[string][]byte{
		"mimetype":    []byte("application/vnd.oasis.opendocument.text"),
		"content.xml": []byte(`<?xml version="1.0"?><office:document-content xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0"/>`),
	})

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(outputDir, name), files[name], 0644); err != nil {
			panic(err)
		}
	}
	fmt.Printf("generated %d knowledge-file samples in %s\n", len(files), outputDir)
}

func minimalPDF() []byte {
	objects := []string{
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>\nendobj\n",
		"4 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n",
		"5 0 obj\n<< /Length 61 >>\nstream\nBT /F1 24 Tf 72 700 Td (InkSpace PDF preview sample) Tj ET\nendstream\nendobj\n",
	}
	var body strings.Builder
	body.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects))
	for i, object := range objects {
		offsets[i] = body.Len()
		body.WriteString(object)
	}
	xref := body.Len()
	fmt.Fprintf(&body, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&body, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&body, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xref)
	return []byte(body.String())
}

func zipBytes(files map[string][]byte) []byte {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry, err := writer.Create(name)
		if err != nil {
			panic(err)
		}
		if _, err := entry.Write(files[name]); err != nil {
			panic(err)
		}
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func gzipBytes(content []byte) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(content); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func mustBase64(value string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return decoded
}
