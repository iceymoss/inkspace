package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		service: service.NewCommentService(),
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	var uid uint
	if exists {
		uid = userID.(uint)
	}

	comment, err := h.service.Create(&req, uid)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// 如果评论状态为 0（待审核），返回审核中的消息
	if comment.Status == 0 {
		utils.SuccessWithMessage(c, "评论已提交，等待审核", comment.ToResponse())
		return
	}

	// 如果评论已通过（status=1），正常返回，不显示额外消息
	utils.Success(c, comment.ToResponse())
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	// 确保用户已登录（middleware 已经验证，但这里再次确认）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")
	roleStr := "user"
	if role != nil {
		roleStr = role.(string)
	}

	if err := h.service.Delete(uint(id), userID.(uint), roleStr); err != nil {
		utils.Error(c, 403, err.Error()) // 使用 403 Forbidden 表示权限不足
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *CommentHandler) GetList(c *gin.Context) {
	var query models.CommentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	comments, total, err := h.service.GetList(&query)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	commentResponses := make([]*models.CommentResponse, len(comments))
	for i, comment := range comments {
		resp := comment.ToResponse()

		// 加载该评论的回复（默认只加载10条）
		// 查询条件：root_id = comment.ID 的所有评论（包括直接回复和回复的回复）
		var replies []*models.Comment
		if err := database.DB.Where("root_id = ?", comment.ID).
			Preload("User").
			Preload("Article").
			Order("created_at ASC").
			Limit(10). // 默认只返回10条子评论
			Find(&replies).Error; err == nil {
			// 转换为响应格式
			replyResponses := make([]models.CommentResponse, len(replies))
			for j, reply := range replies {
				replyResponses[j] = *reply.ToResponse()
			}
			resp.Replies = replyResponses
		}

		commentResponses[i] = resp
	}

	utils.PageResponse(c, commentResponses, total, query.Page, query.PageSize)
}

func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// GetReplies 获取子评论列表（分页）
func (h *CommentHandler) GetReplies(c *gin.Context) {
	rootID, err := strconv.ParseUint(c.Param("root_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的根评论ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 查询子评论总数
	var total int64
	if err := database.DB.Model(&models.Comment{}).
		Where("root_id = ?", rootID).
		Count(&total).Error; err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 查询子评论列表
	var replies []*models.Comment
	offset := (page - 1) * pageSize
	if err := database.DB.Where("root_id = ?", rootID).
		Preload("User").
		Preload("Article").
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&replies).Error; err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 转换为响应格式
	replyResponses := make([]*models.CommentResponse, len(replies))
	for i, reply := range replies {
		replyResponses[i] = reply.ToResponse()
	}

	utils.PageResponse(c, replyResponses, total, page, pageSize)
}
