package handler

import (
	"strconv"

	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	service *service.LikeService
}

func NewLikeHandler() *LikeHandler {
	return &LikeHandler{
		service: service.NewLikeService(),
	}
}

// LikeArticle 点赞文章
func (h *LikeHandler) LikeArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	ip := c.ClientIP()

	if err := h.service.LikeArticle(userID, uint(articleID), ip); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeArticle 取消点赞文章
func (h *LikeHandler) UnlikeArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	ip := c.ClientIP()

	if err := h.service.UnlikeArticle(userID, uint(articleID), ip); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "取消点赞成功", nil)
}

// CheckArticleLiked 检查是否已点赞文章
func (h *LikeHandler) CheckArticleLiked(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	ip := c.ClientIP()

	isLiked, err := h.service.IsArticleLiked(userID, uint(articleID), ip)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"is_liked": isLiked})
}

// LikeComment 点赞评论
func (h *LikeHandler) LikeComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评论ID")
		return
	}

	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	ip := c.ClientIP()

	if err := h.service.LikeComment(userID, uint(commentID), ip); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeComment 取消点赞评论
func (h *LikeHandler) UnlikeComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评论ID")
		return
	}

	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	ip := c.ClientIP()

	if err := h.service.UnlikeComment(userID, uint(commentID), ip); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "取消点赞成功", nil)
}

