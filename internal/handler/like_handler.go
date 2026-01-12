package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

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

// LikeWork 点赞作品
// POST /api/works/:id/like
func (h *LikeHandler) LikeWork(c *gin.Context) {
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

	if err := h.service.LikeWork(userID.(uint), uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "操作成功", nil)
}

// LikeArticle 点赞文章
// POST /api/articles/:id/like
func (h *LikeHandler) LikeArticle(c *gin.Context) {
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

	if err := h.service.LikeArticle(userID.(uint), uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "操作成功", nil)
}

// CheckWorkLiked 检查是否已点赞作品
// GET /api/works/:id/liked
func (h *LikeHandler) CheckWorkLiked(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		// 兼容两种字段名
		utils.Success(c, gin.H{"liked": false, "is_liked": false})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	liked, err := h.service.CheckWorkLiked(userID.(uint), uint(id))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 返回两种字段名以兼容前后端
	utils.Success(c, gin.H{"liked": liked, "is_liked": liked})
}

// CheckArticleLiked 检查是否已点赞文章
// GET /api/articles/:id/is-liked (旧版)
// GET /api/articles/:id/liked (新版)
func (h *LikeHandler) CheckArticleLiked(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		// 兼容两种字段名
		utils.Success(c, gin.H{"liked": false, "is_liked": false})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	liked, err := h.service.CheckArticleLiked(userID.(uint), uint(id))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 返回两种字段名以兼容前后端
	utils.Success(c, gin.H{"liked": liked, "is_liked": liked})
}

// UnlikeArticle 取消点赞文章（兼容旧API）
// DELETE /api/articles/:id/like
func (h *LikeHandler) UnlikeArticle(c *gin.Context) {
	// 与 LikeArticle 相同，点赞是toggle操作
	h.LikeArticle(c)
}

// LikeComment 点赞评论
// POST /api/comments/:id/like
func (h *LikeHandler) LikeComment(c *gin.Context) {
	utils.SuccessWithMessage(c, "评论点赞功能待实现", nil)
}

// UnlikeComment 取消点赞评论
// DELETE /api/comments/:id/like
func (h *LikeHandler) UnlikeComment(c *gin.Context) {
	utils.SuccessWithMessage(c, "评论取消点赞功能待实现", nil)
}

// CheckCommentLiked 检查是否已点赞评论
// GET /api/comments/:id/is-liked
func (h *LikeHandler) CheckCommentLiked(c *gin.Context) {
	// 兼容两种字段名
	utils.Success(c, gin.H{"liked": false, "is_liked": false})
}
