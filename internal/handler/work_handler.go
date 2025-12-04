package handler

import (
	"encoding/json"
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type WorkHandler struct {
	service *service.WorkService
}

func NewWorkHandler() *WorkHandler {
	return &WorkHandler{
		service: service.NewWorkService(),
	}
}

func (h *WorkHandler) Create(c *gin.Context) {
	var req models.WorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	work, err := h.service.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, h.toResponse(work))
}

func (h *WorkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.WorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	work, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, h.toResponse(work))
}

func (h *WorkHandler) Delete(c *gin.Context) {
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

func (h *WorkHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	work, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "作品不存在")
		return
	}

	// Increment view count
	go h.service.IncrementViewCount(uint(id))

	utils.Success(c, h.toResponse(work))
}

func (h *WorkHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	works, total, err := h.service.GetList(page, pageSize, status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	workResponses := make([]*models.WorkResponse, len(works))
	for i, work := range works {
		workResponses[i] = h.toResponse(work)
	}

	utils.PageResponse(c, workResponses, total, page, pageSize)
}

func (h *WorkHandler) toResponse(work *models.Work) *models.WorkResponse {
	var images []string
	if work.Images != "" {
		json.Unmarshal([]byte(work.Images), &images)
	}

	return &models.WorkResponse{
		ID:          work.ID,
		Title:       work.Title,
		Description: work.Description,
		Cover:       work.Cover,
		Images:      images,
		Link:        work.Link,
		Sort:        work.Sort,
		ViewCount:   work.ViewCount,
		Status:      work.Status,
		CreatedAt:   work.CreatedAt,
		UpdatedAt:   work.UpdatedAt,
	}
}

