package handler

import (
	"errors"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

func knowledgeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrKnowledgeNotFound):
		utils.NotFound(c, err.Error())
	case errors.Is(err, service.ErrShareDisabled), errors.Is(err, service.ErrShareExpired):
		utils.Forbidden(c, err.Error())
	case errors.Is(err, service.ErrCatalogCycle), errors.Is(err, service.ErrKnowledgeInvalid):
		utils.BadRequest(c, err.Error())
	default:
		zap.L().Error("knowledge base request failed", zap.Error(err))
		utils.InternalServerError(c, "服务暂时不可用")
	}
}

func pathUint(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil || value == 0 {
		utils.BadRequest(c, "无效的ID")
		return 0, false
	}
	return uint(value), true
}

func currentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return 0, false
	}
	userID, ok := value.(uint)
	if !ok {
		utils.Unauthorized(c, "登录信息无效")
		return 0, false
	}
	return userID, true
}
