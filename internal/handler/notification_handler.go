package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

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

// GetNotifications 获取通知列表
// GET /api/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	onlyUnread := c.Query("only_unread") == "true"

	notifications, total, err := h.service.GetNotifications(userID.(uint), page, pageSize, onlyUnread)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	responses := make([]*models.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = notification.ToResponse()
	}

	utils.PageResponse(c, responses, total, page, pageSize)
}

// GetUnreadCount 获取未读通知数量
// GET /api/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Success(c, gin.H{"count": 0})
		return
	}

	count, err := h.service.GetUnreadCount(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"count": count})
}

// MarkAsRead 标记通知为已读
// PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.MarkAsRead(uint(id), userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "标记成功", nil)
}

// MarkAllAsRead 标记所有通知为已读
// PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.MarkAllAsRead(userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "全部标记为已读", nil)
}

// DeleteNotification 删除通知
// DELETE /api/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.DeleteNotification(uint(id), userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// DeleteAllRead 删除所有已读通知
// DELETE /api/notifications/read-all
func (h *NotificationHandler) DeleteAllRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.DeleteAllRead(userID.(uint)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
