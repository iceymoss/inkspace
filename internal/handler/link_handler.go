package handler

import (
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{
		service: service.NewLinkService(),
	}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var req models.LinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	link, err := h.service.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, link.ToResponse())
}

func (h *LinkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.LinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	link, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, link.ToResponse())
}

func (h *LinkHandler) Delete(c *gin.Context) {
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

func (h *LinkHandler) GetList(c *gin.Context) {
	// 管理后台：分页 + 筛选 + 排序
	if c.Request.URL.Path == "/api/admin/links" {
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

		keyword := strings.TrimSpace(c.Query("keyword"))

		var statusPtr *int
		if statusStr := c.Query("status"); strings.TrimSpace(statusStr) != "" {
			if s, err := strconv.Atoi(statusStr); err == nil {
				statusPtr = &s
			}
		}

		sort := c.DefaultQuery("sort", "")

		links, total, err := h.service.GetAdminList(page, pageSize, keyword, statusPtr, sort)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		linkResponses := make([]*models.LinkResponse, len(links))
		for i, link := range links {
			linkResponses[i] = link.ToResponse()
		}

		utils.PageResponse(c, linkResponses, total, page, pageSize)
		return
	}

	// 前台：保持原有逻辑，按状态过滤，返回完整列表
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	links, err := h.service.GetList(status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	linkResponses := make([]*models.LinkResponse, len(links))
	for i, link := range links {
		linkResponses[i] = link.ToResponse()
	}

	utils.Success(c, linkResponses)
}
