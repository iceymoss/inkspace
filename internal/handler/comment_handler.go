package handler

import (
	"strconv"

	"mysite/internal/database"
	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

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

	utils.Success(c, comment.ToResponse())
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if err := h.service.Delete(uint(id), userID.(uint), role.(string)); err != nil {
		utils.Error(c, 400, err.Error())
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
		
		// 加载该评论的所有回复（包括回复的回复）
		// 查询条件：root_id = comment.ID 的所有评论（包括直接回复和回复的回复）
		var replies []*models.Comment
		if err := database.DB.Where("root_id = ?", comment.ID).
			Preload("User").
			Preload("Article").
			Order("created_at ASC").
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

