package handler

import (
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.TagService
}

func NewTagHandler() *TagHandler {
	return &TagHandler{
		service: service.NewTagService(),
	}
}

func (h *TagHandler) Create(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 如果是用户创建标签，自动设置user_id
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		req.UserID = &uid
	}

	tag, err := h.service.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, tag.ToResponse())
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tag, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, tag.ToResponse())
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *TagHandler) GetList(c *gin.Context) {
	// 管理后台：支持分页、筛选和排序
	if c.Request.URL.Path == "/api/admin/tags" {
		// 分页参数
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("page_size", "20")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}

		// 筛选参数
		keyword := strings.TrimSpace(c.Query("keyword"))
		scope := c.DefaultQuery("scope", "") // all/system/user

		var hasArticlesPtr *int
		hasArticlesStr := strings.TrimSpace(c.Query("has_articles"))
		if hasArticlesStr != "" {
			if v, err := strconv.Atoi(hasArticlesStr); err == nil {
				hasArticlesPtr = &v
			}
		}

		// 排序参数，例如：id_desc、article_count_desc、created_at_desc
		sort := c.DefaultQuery("sort", "")

		tags, total, err := h.service.GetAdminList(page, pageSize, keyword, scope, hasArticlesPtr, sort)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		tagResponses := make([]*models.TagResponse, len(tags))
		for i, tag := range tags {
			tagResponses[i] = tag.ToResponse()
		}

		utils.PageResponse(c, tagResponses, total, page, pageSize)
		return
	}

	// 前台：保持原来的简单列表（按文章数降序）
	tags, err := h.service.GetList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	tagResponses := make([]*models.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = tag.ToResponse()
	}

	utils.Success(c, tagResponses)
}
