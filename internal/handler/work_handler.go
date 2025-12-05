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

	// 获取当前用户ID和角色
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}
	role, _ := c.Get("role")

	work, err := h.service.Create(&req, userID.(uint), role.(string))
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

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	work, err := h.service.Update(uint(id), &req, userID.(uint), role.(string))
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

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if err := h.service.Delete(uint(id), userID.(uint), role.(string)); err != nil {
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
	workType := c.Query("type")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	works, total, err := h.service.GetList(page, pageSize, workType, status)
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

// GetRecommended 获取推荐作品
// GET /api/works/recommended
func (h *WorkHandler) GetRecommended(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}
	if limit > 10 {
		limit = 10
	}

	works, err := h.service.GetRecommended(limit)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	workResponses := make([]*models.WorkResponse, len(works))
	for i, work := range works {
		workResponses[i] = h.toResponse(work)
	}

	utils.Success(c, workResponses)
}

func (h *WorkHandler) toResponse(work *models.Work) *models.WorkResponse {
	var images []models.PhotoItem
	if work.Images != "" {
		json.Unmarshal([]byte(work.Images), &images)
	}

	var metadata map[string]interface{}
	if work.Metadata != "" {
		json.Unmarshal([]byte(work.Metadata), &metadata)
	}

	resp := &models.WorkResponse{
		ID:            work.ID,
		Title:         work.Title,
		Type:          work.Type,
		Metadata:      metadata,
		DailyQuota:    work.DailyQuota,
		Description:   work.Description,
		Cover:         work.Cover,
		Images:        images,
		Link:          work.Link,
		GithubURL:     work.GithubURL,
		DemoURL:       work.DemoURL,
		TechStack:     work.TechStack,
		AuthorID:      work.AuthorID,
		Sort:          work.Sort,
		ViewCount:     work.ViewCount,
		CommentCount:  work.CommentCount,
		LikeCount:     work.LikeCount,
		FavoriteCount: work.FavoriteCount,
		Status:        work.Status,
		IsRecommend:   work.IsRecommend,
		CreatedAt:     work.CreatedAt,
		UpdatedAt:     work.UpdatedAt,
	}

	if work.Author != nil {
		resp.Author = work.Author.ToResponse()
	}

	return resp
}

// SetRecommend 设置作品推荐状态
// PUT /api/admin/works/:id/recommend
func (h *WorkHandler) SetRecommend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		IsRecommend bool `json:"is_recommend"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.SetRecommend(uint(id), req.IsRecommend); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "设置成功", nil)
}

// GetHotWorks 获取热门作品
// GET /api/works/hot
func (h *WorkHandler) GetHotWorks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "4")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 4
	}
	if limit > 10 {
		limit = 10
	}

	works, err := h.service.GetHotWorks(limit)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	workResponses := make([]*models.WorkResponse, len(works))
	for i, work := range works {
		workResponses[i] = h.toResponse(work)
	}

	utils.Success(c, workResponses)
}

// GetMyWorks 获取我的作品列表
// GET /api/works/my
func (h *WorkHandler) GetMyWorks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	workType := c.Query("type")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	works, total, err := h.service.GetMyWorks(userID.(uint), page, pageSize, workType, status)
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

// GetQuotaUsage 获取今日配额使用情况
// GET /api/works/quota
func (h *WorkHandler) GetQuotaUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	usage, err := h.service.GetTodayQuotaUsage(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"used":      usage,
		"limit":     3,
		"remaining": 3 - usage,
	})
}
