//go:build ignore

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/pkg/uploader"
)

const (
	sampleDirectory = "docs/samples/knowledge-files"
	workspaceName   = "文件类型样例"
	username        = "iceymoss"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := database.Init(); err != nil {
		panic(err)
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		panic(err)
	}

	var workspace models.Workspace
	result := database.DB.Where("owner_id = ? AND name = ?", user.ID, workspaceName).First(&workspace)
	if result.Error != nil {
		workspace = models.Workspace{
			OwnerID: user.ID, Name: workspaceName,
			Description: "知识库多文件类型上传与预览联调样例",
		}
		if err := database.DB.Create(&workspace).Error; err != nil {
			panic(err)
		}
	}

	storage, err := uploader.NewStorage()
	if err != nil {
		panic(err)
	}
	uploadConfig := config.AppConfig.Upload
	uploadConfig.KnowledgeAllowedExtensions = append(uploadConfig.KnowledgeAllowedExtensions, ".bin")
	fileService := service.NewKnowledgeFileServiceWithStorage(storage, uploadConfig)

	entries, err := os.ReadDir(sampleDirectory)
	if err != nil {
		panic(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "README.md" || strings.HasPrefix(entry.Name(), "invalid-") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	created, skipped := 0, 0
	for _, name := range names {
		var count int64
		if err := database.DB.Model(&models.Doc{}).
			Where("workspace_id = ? AND owner_id = ? AND title = ?", workspace.ID, user.ID, name).
			Count(&count).Error; err != nil {
			panic(err)
		}
		if count > 0 {
			fmt.Printf("skip existing %s\n", name)
			skipped++
			continue
		}

		header, cleanup := multipartFileHeader(filepath.Join(sampleDirectory, name), name)
		_, err := fileService.Upload(context.Background(), workspace.ID, user.ID, &models.KnowledgeFileUploadRequest{Title: name}, header)
		cleanup()
		if err != nil {
			panic(fmt.Errorf("upload %s: %w", name, err))
		}
		fmt.Printf("uploaded %s\n", name)
		created++
	}

	fmt.Printf("workspace=%q id=%d owner=%s created=%d skipped=%d\n", workspace.Name, workspace.ID, user.Username, created, skipped)
}

func multipartFileHeader(path, name string) (*multipart.FileHeader, func()) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(32 << 20); err != nil {
		panic(err)
	}
	files := request.MultipartForm.File["file"]
	if len(files) != 1 {
		panic("multipart file was not created")
	}
	return files[0], func() { _ = request.MultipartForm.RemoveAll() }
}
