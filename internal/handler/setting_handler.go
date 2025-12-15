package handler

import (
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	service *service.SettingService
}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{
		service: service.NewSettingService(),
	}
}

// GetPublicSettings 获取公开配置
func (h *SettingHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.service.GetPublic()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 转换为map格式
	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	utils.Success(c, result)
}

// GetAllSettings 获取所有配置（管理员）
func (h *SettingHandler) GetAllSettings(c *gin.Context) {
	settings, err := h.service.GetAll()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	settingResponses := make([]*models.SettingResponse, len(settings))
	for i, setting := range settings {
		settingResponses[i] = setting.ToResponse()
	}

	utils.Success(c, settingResponses)
}

// UpdateSetting 更新配置
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	var req models.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	setting, err := h.service.Set(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, setting.ToResponse())
}

// BatchUpdateSettings 批量更新配置
func (h *SettingHandler) BatchUpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.BatchSet(req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "配置更新成功", nil)
}

// DeleteSetting 删除配置
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.BadRequest(c, "配置键不能为空")
		return
	}

	if err := h.service.Delete(key); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
