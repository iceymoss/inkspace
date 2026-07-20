package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type WorkspaceHandler struct{ service *service.WorkspaceService }

func NewWorkspaceHandler() *WorkspaceHandler {
	return &WorkspaceHandler{service: service.NewWorkspaceService()}
}

func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req models.WorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	workspace, err := h.service.Create(&req, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, workspace.ToResponse())
}

func (h *WorkspaceHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	workspaces, err := h.service.List(userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	responses := make([]*models.WorkspaceResponse, len(workspaces))
	for i := range workspaces {
		responses[i] = workspaces[i].ToResponse()
	}
	utils.Success(c, responses)
}

func (h *WorkspaceHandler) Get(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	workspace, err := h.service.Get(id, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, workspace.ToResponse())
}

func (h *WorkspaceHandler) Update(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var req models.WorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	workspace, err := h.service.Update(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, workspace.ToResponse())
}

func (h *WorkspaceHandler) Delete(c *gin.Context) {
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
