package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		service: service.NewNotificationService(),
	}
}

// GetList 获取通知列表
func (h *NotificationHandler) GetList(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var query models.NotificationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}

	notifications, total, err := h.service.GetList(userID.(uint), &query)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	notificationResponses := make([]*models.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		notificationResponses[i] = notification.ToResponse()
	}

	utils.PageResponse(c, notificationResponses, total, query.Page, query.PageSize)
}

// MarkAsRead 标记为已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.service.MarkAsRead(uint(id), userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "标记成功", nil)
}

// MarkAllAsRead 全部标记为已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.service.MarkAllAsRead(userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "全部已读", nil)
}

// GetUnreadCount 获取未读数量
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.service.GetUnreadCount(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"count": count})
}

// Delete 删除通知
func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.service.Delete(uint(id), userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

