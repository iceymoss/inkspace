package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		service: service.NewFavoriteService(),
	}
}

// AddFavorite 收藏文章
// POST /api/articles/:id/favorite
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.service.AddFavorite(userID.(uint), uint(articleID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "收藏成功", nil)
}

// RemoveFavorite 取消收藏
// DELETE /api/articles/:id/favorite
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.service.RemoveFavorite(userID.(uint), uint(articleID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "取消收藏成功", nil)
}

// CheckFavorited 检查是否已收藏
// GET /api/articles/:id/is-favorited
func (h *FavoriteHandler) CheckFavorited(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Success(c, gin.H{"is_favorited": false})
		return
	}

	isFavorited, err := h.service.IsFavorited(userID.(uint), uint(articleID))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"is_favorited": isFavorited})
}

// GetMyFavorites 获取我的收藏列表
// GET /api/favorites
func (h *FavoriteHandler) GetMyFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var query models.FavoriteListQuery
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

	list, total, err := h.service.GetFavoriteList(userID.(uint), query.Page, query.PageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageResponse(c, list, total, query.Page, query.PageSize)
}

// GetUserFavorites 获取指定用户的收藏列表
// GET /api/users/:id/favorites
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var query models.FavoriteListQuery
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

	list, total, err := h.service.GetFavoriteList(uint(userID), query.Page, query.PageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.PageResponse(c, list, total, query.Page, query.PageSize)
}

