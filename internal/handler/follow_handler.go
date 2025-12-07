package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	service *service.FollowService
}

func NewFollowHandler() *FollowHandler {
	return &FollowHandler{
		service: service.NewFollowService(),
	}
}

// Follow 关注用户
// POST /api/users/:id/follow
func (h *FollowHandler) Follow(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	followerID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.Follow(followerID.(uint), uint(followingID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "关注成功", nil)
}

// Unfollow 取消关注
// DELETE /api/users/:id/follow
func (h *FollowHandler) Unfollow(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	followerID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.Unfollow(followerID.(uint), uint(followingID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "取消关注成功", nil)
}

// GetFollowStats 获取关注统计
// GET /api/users/:id/follow-stats
func (h *FollowHandler) GetFollowStats(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	// 获取当前登录用户ID（可能为0，表示未登录）
	currentUserID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		currentUserID = uid.(uint)
	}

	stats, err := h.service.GetFollowStats(uint(userID), currentUserID)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, stats)
}

// GetFollowingList 获取关注列表
// GET /api/users/:id/following
func (h *FollowHandler) GetFollowingList(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var query models.FollowListQuery
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

	list, total, err := h.service.GetFollowingList(uint(userID), query.Page, query.PageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageResponse(c, list, total, query.Page, query.PageSize)
}

// GetFollowerList 获取粉丝列表
// GET /api/users/:id/followers
func (h *FollowHandler) GetFollowerList(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var query models.FollowListQuery
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

	list, total, err := h.service.GetFollowerList(uint(userID), query.Page, query.PageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageResponse(c, list, total, query.Page, query.PageSize)
}

