package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type DocHandler struct{ service *service.DocService }

func NewDocHandler() *DocHandler { return &DocHandler{service: service.NewDocService()} }

func (h *DocHandler) Create(c *gin.Context) {
	var req models.DocCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	doc, err := h.service.Create(&req, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) List(c *gin.Context) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var query models.DocListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	docs, err := h.service.List(workspaceID, userID, query.CatalogID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	responses := make([]*models.DocResponse, len(docs))
	for i := range docs {
		responses[i] = docs[i].ToResponse()
		responses[i].Summary = service.ContentSummary(docs[i].Content)
		responses[i].Content = ""
		responses[i].ContentHTML = ""
	}
	utils.Success(c, responses)
}

func (h *DocHandler) GetEdit(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	doc, err := h.service.GetEdit(id, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) Save(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	var req models.DocSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	doc, err := h.service.Save(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) Autosave(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	var req models.DocAutosaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	doc, err := h.service.Autosave(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) Publish(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	status := models.DocStatusPublished
	if c.Request.ContentLength > 0 {
		var req models.DocPublishRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
		if req.Status != nil {
			status = *req.Status
		}
	}
	doc, err := h.service.Publish(id, userID, status)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) PublishToBlog(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	var req models.DocPublishToBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	article, err := h.service.PublishToBlog(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, article.ToResponse())
}

func (h *DocHandler) Delete(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		knowledgeError(c, err)
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *DocHandler) Move(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	var req models.DocMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	doc, err := h.service.Move(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) Versions(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	versions, err := h.service.Versions(id, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	responses := make([]*models.DocVersionResponse, len(versions))
	for i := range versions {
		responses[i] = versions[i].ToResponse()
		responses[i].Content = ""
	}
	utils.Success(c, responses)
}

func (h *DocHandler) Version(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	version, ok := pathVersion(c)
	if !ok {
		return
	}
	snapshot, err := h.service.Version(id, userID, version)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, snapshot.ToResponse())
}

func (h *DocHandler) Rollback(c *gin.Context) {
	id, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	version, ok := pathVersion(c)
	if !ok {
		return
	}
	doc, err := h.service.Rollback(id, userID, version)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}

func (h *DocHandler) Search(c *gin.Context) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.service.Search(workspaceID, userID, c.Query("q"))
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, result)
}

func docAndUser(c *gin.Context) (uint, uint, bool) {
	id, ok := pathUint(c, "id")
	if !ok {
		return 0, 0, false
	}
	userID, ok := currentUserID(c)
	return id, userID, ok
}

func pathVersion(c *gin.Context) (int, bool) {
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil || version <= 0 {
		utils.BadRequest(c, "无效的版本号")
		return 0, false
	}
	return version, true
}
