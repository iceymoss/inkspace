package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		service: service.NewCategoryService(),
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.service.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, category.ToResponse())
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, category.ToResponse())
}

func (h *CategoryHandler) Delete(c *gin.Context) {
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

func (h *CategoryHandler) GetList(c *gin.Context) {
	// 区分管理后台与前台博客的分类列表需求：
	// - 管理后台（/api/admin/categories）：需要分页，并返回 total 等分页信息
	// - 前台博客（/api/categories）：一次性返回全部分类，用于筛选下拉等

	// 管理后台走分页逻辑
	if c.Request.URL.Path == "/api/admin/categories" {
		// 解析分页参数
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("page_size", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		categories, total, err := h.service.GetListPaged(page, pageSize)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		categoryResponses := make([]*models.CategoryResponse, len(categories))
		for i, category := range categories {
			categoryResponses[i] = category.ToResponse()
		}

		utils.PageResponse(c, categoryResponses, total, page, pageSize)
		return
	}

	// 默认：前台接口返回全部分类列表
	categories, err := h.service.GetList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	categoryResponses := make([]*models.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = category.ToResponse()
	}

	utils.Success(c, categoryResponses)
}
