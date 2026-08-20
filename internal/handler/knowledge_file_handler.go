package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
	"go.uber.org/zap"
)

func (h *DocHandler) UploadFile(c *gin.Context) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	fileService, err := service.NewKnowledgeFileService()
	if err != nil {
		knowledgeError(c, err)
		return
	}
	if err := fileService.AuthorizeUpload(workspaceID, userID); err != nil {
		knowledgeError(c, err)
		return
	}
	const multipartOverhead = int64(1 * 1024 * 1024)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, fileService.MaxSize()+multipartOverhead)
	var req models.KnowledgeFileUploadRequest
	bindErr := c.ShouldBind(&req)
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}
	if bindErr != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(bindErr, &maxBytesError) {
			knowledgeError(c, service.ErrKnowledgeFileTooLarge)
			return
		}
		utils.BadRequest(c, bindErr.Error())
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}
	doc, err := fileService.Upload(c.Request.Context(), workspaceID, userID, &req, file)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc)
}

func (h *DocHandler) DownloadFile(c *gin.Context) {
	docID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, _ := optionalUserID(c)
	fileService, err := service.NewKnowledgeFileService()
	if err != nil {
		knowledgeError(c, err)
		return
	}
	download, err := fileService.OpenDownload(c.Request.Context(), docID, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	defer download.Reader.Close()

	fallbackName := fmt.Sprintf("document-%d%s", docID, filepath.Ext(download.FileName))
	c.Header("Content-Type", download.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q; filename*=UTF-8''%s", fallbackName, url.PathEscape(download.FileName)))
	c.Header("Content-Length", fmt.Sprintf("%d", download.Size))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'none'; sandbox")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, download.Reader); err != nil {
		zap.L().Error("stream knowledge file failed", zap.Uint("doc_id", docID), zap.Uint("attachment_id", download.AttachmentID), zap.Error(err))
	}
}

func (h *DocHandler) PreviewFile(c *gin.Context) {
	docID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, _ := optionalUserID(c)
	fileService, err := service.NewKnowledgeFileService()
	if err != nil {
		knowledgeError(c, err)
		return
	}
	preview, err := fileService.OpenPreview(c.Request.Context(), docID, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	defer preview.Reader.Close()

	c.Header("Content-Type", preview.MimeType)
	c.Header("Content-Length", fmt.Sprintf("%d", preview.Size))
	c.Header("Content-Disposition", "inline")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, preview.Reader); err != nil {
		zap.L().Error("stream knowledge preview failed", zap.Uint("doc_id", docID), zap.Uint("attachment_id", preview.AttachmentID), zap.Error(err))
	}
}
