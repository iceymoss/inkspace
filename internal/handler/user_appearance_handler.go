package handler

import (
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserAppearanceHandler struct {
	service *service.UserAppearanceService
}

func NewUserAppearanceHandler() *UserAppearanceHandler {
	return &UserAppearanceHandler{service: service.NewUserAppearanceService()}
}

func (h *UserAppearanceHandler) Get(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	appearance, err := h.service.Get(userID.(uint))
	if err != nil {
		zap.L().Error("get user appearance failed", zap.Error(err))
		utils.InternalServerError(c, "服务暂时不可用")
		return
	}

	utils.Success(c, appearance)
}

func (h *UserAppearanceHandler) Update(c *gin.Context) {
	var req models.UserAppearanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	appearance, err := h.service.Save(userID.(uint), &req)
	if err != nil {
		zap.L().Error("save user appearance failed", zap.Error(err))
		utils.InternalServerError(c, "服务暂时不可用")
		return
	}

	utils.Success(c, appearance)
}
