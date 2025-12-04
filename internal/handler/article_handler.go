package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{
		service: service.NewArticleService(),
	}
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	article, err := h.service.Create(&req, userID.(uint))
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, article.ToResponse())
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	article, err := h.service.Update(uint(id), &req, userID.(uint), role.(string))
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, article.ToResponse())
}

func (h *ArticleHandler) Delete(c *gin.Context) {
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

func (h *ArticleHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	article, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// Increment view count
	go h.service.IncrementViewCount(uint(id))

	utils.Success(c, article.ToResponse())
}

func (h *ArticleHandler) GetList(c *gin.Context) {
	var query models.ArticleListQuery
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

	articles, total, err := h.service.GetList(&query)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	articleResponses := make([]*models.ArticleResponse, len(articles))
	for i, article := range articles {
		resp := article.ToResponse()
		// Don't return full content in list
		resp.Content = ""
		articleResponses[i] = resp
	}

	utils.PageResponse(c, articleResponses, total, query.Page, query.PageSize)
}

func (h *ArticleHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.IncrementLikeCount(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "点赞成功", nil)
}

