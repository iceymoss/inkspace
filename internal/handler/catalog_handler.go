package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type CatalogHandler struct{ service *service.CatalogService }

func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{service: service.NewCatalogService()}
}

func (h *CatalogHandler) Create(c *gin.Context) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var req models.CatalogCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	catalog, err := h.service.Create(workspaceID, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, catalog.ToResponse())
}

func (h *CatalogHandler) Tree(c *gin.Context) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	tree, err := h.service.Tree(workspaceID, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, tree)
}

func (h *CatalogHandler) Update(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var req models.CatalogUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	catalog, err := h.service.Update(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, catalog.ToResponse())
}

func (h *CatalogHandler) Move(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var req models.CatalogMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	catalog, err := h.service.Move(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, catalog.ToResponse())
}

func (h *CatalogHandler) Delete(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		knowledgeError(c, err)
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
