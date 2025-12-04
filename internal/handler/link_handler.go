package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

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

