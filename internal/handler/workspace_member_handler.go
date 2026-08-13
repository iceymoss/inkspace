package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type WorkspaceMemberHandler struct {
	service *service.WorkspaceMemberService
}

func NewWorkspaceMemberHandler() *WorkspaceMemberHandler {
	return &WorkspaceMemberHandler{service: service.NewWorkspaceMemberService()}
}

func (h *WorkspaceMemberHandler) List(c *gin.Context) {
	workspaceID, userID, ok := workspaceAndUser(c)
	if !ok {
		return
	}
	members, err := h.service.List(workspaceID, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, members)
}

func (h *WorkspaceMemberHandler) Add(c *gin.Context) {
	workspaceID, userID, ok := workspaceAndUser(c)
	if !ok {
		return
	}
	var req models.WorkspaceMemberCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	member, err := h.service.Add(workspaceID, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, member)
}

func (h *WorkspaceMemberHandler) Update(c *gin.Context) {
	workspaceID, userID, ok := workspaceAndUser(c)
	if !ok {
		return
	}
	targetUserID, ok := pathUint(c, "userId")
	if !ok {
		return
	}
	var req models.WorkspaceMemberUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	member, err := h.service.Update(workspaceID, targetUserID, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, member)
}

func (h *WorkspaceMemberHandler) Delete(c *gin.Context) {
	workspaceID, userID, ok := workspaceAndUser(c)
	if !ok {
		return
	}
	targetUserID, ok := pathUint(c, "userId")
	if !ok {
		return
	}
	if err := h.service.Delete(workspaceID, targetUserID, userID); err != nil {
		knowledgeError(c, err)
		return
	}
	utils.SuccessWithMessage(c, "移除成功", nil)
}

func workspaceAndUser(c *gin.Context) (uint, uint, bool) {
	workspaceID, ok := pathUint(c, "id")
	if !ok {
		return 0, 0, false
	}
	userID, ok := currentUserID(c)
	return workspaceID, userID, ok
}
