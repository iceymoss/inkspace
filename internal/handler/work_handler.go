package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")
	roleStr := "user"
	if role != nil {
		roleStr = role.(string)
	}

	work, err := h.service.Update(uint(id), &req, userID.(uint), roleStr)
	if err != nil {
		// 权限错误返回403
		if err.Error() == "无权限修改此作品" {
			utils.Error(c, 403, err.Error())
		} else {
			utils.Error(c, 400, err.Error())
		}
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

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")
	roleStr := "user"
	if role != nil {
		roleStr = role.(string)
	}

	if err := h.service.Delete(uint(id), userID.(uint), roleStr); err != nil {
		// 权限错误返回403
		if err.Error() == "无权限删除此作品" {
			utils.Error(c, 403, err.Error())
		} else {
			utils.Error(c, 400, err.Error())
		}
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

	// 获取当前用户信息（可选认证）
	userID, userExists := c.Get("user_id")
	role, _ := c.Get("role")
	roleStr := "user"
	if role != nil {
		roleStr = role.(string)
	}

	// 权限检查：
	// - 草稿（status=0）：只有作者或管理员可以查看
	// - 待审核（status=2）：只有作者或管理员可以查看
	// - 审核不通过（status=3）：只有作者或管理员可以查看
	// - 已发布（status=1）：所有人可以查看
	if work.Status != 1 {
		// 如果不是管理员，且不是作者，则无权限查看
		if roleStr != "admin" && (!userExists || userID.(uint) != work.AuthorID) {
			utils.NotFound(c, "作品不存在")
			return
		}
	}

	// Increment view count only if not skipping (skip_view=true means skip increment)
	// 只有已发布的作品才增加浏览量
	skipView := c.Query("skip_view") == "true"
	if !skipView && work.Status == 1 {
		// 同步增加浏览量，确保返回的数据包含最新的浏览量
		if err := h.service.IncrementViewCount(uint(id)); err == nil {
			// 更新返回数据中的浏览量（+1）
			work.ViewCount++
		}
	}

	utils.Success(c, h.toResponse(work))
}

// GetEdit 获取作品详情用于编辑（需要认证，且只允许作者或管理员访问）
// GET /api/works/:id/edit
func (h *WorkHandler) GetEdit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	// 必须登录
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	role, _ := c.Get("role")
	roleStr := "user"
	if role != nil {
		roleStr = role.(string)
	}

	// 查询作品，但使用WHERE条件确保权限（非管理员只能查询自己的作品）
	var work models.Work
	query := database.DB.Where("id = ?", uint(id))

	// 非管理员只能查询自己的作品
	if roleStr != "admin" {
		query = query.Where("author_id = ?", userID.(uint))
	}

	if err := query.Preload("Author").First(&work).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 403, "作品不存在或无权限编辑")
		} else {
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	// 返回作品数据（不增加浏览量）
	utils.Success(c, h.toResponse(&work))
}

func (h *WorkHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	workType := c.Query("type")
	sortBy := c.DefaultQuery("sort", "") // 排序方式：hot, time, view, like
	keyword := c.Query("keyword")        // 标题或描述关键字

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	} else {
		// 如果未指定状态，根据路由判断：
		// - 管理后台（/admin/works）：显示所有状态（status = nil）
		// - 用户端（/api/works）：只显示已发布的作品（status = 1）
		if !strings.Contains(c.Request.URL.Path, "/admin/") {
			// 用户端，默认只显示已发布的作品
			s := 1
			status = &s
		}
		// 管理后台，status 保持为 nil，显示所有状态
	}

	works, total, err := h.service.GetList(page, pageSize, workType, status, sortBy, keyword)
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
		IsPublished:   work.Status == 1, // 只有已发布（status=1）的作品才允许评论
		IsRecommend:   work.IsRecommend,
		AuditMessage:  work.AuditMessage,
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

// UpdateWorkStatus 更新作品审核状态（管理后台使用）
// PUT /api/admin/works/:id/status
func (h *WorkHandler) UpdateWorkStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Status       int    `json:"status" binding:"required"` // 1=通过, 3=拒绝
		AuditMessage string `json:"audit_message"`             // 审核消息（可选，用于记录审核通过或拒绝的原因）
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 验证状态值
	if req.Status != 1 && req.Status != 3 {
		utils.BadRequest(c, "无效的状态值，只能是1（通过）或3（拒绝）")
		return
	}

	if err := h.service.UpdateWorkStatus(uint(id), req.Status, req.AuditMessage); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "审核成功", nil)
}

// GetHotWorks 获取热门作品（支持分页）
// GET /api/works/hot?page=1&page_size=10
func (h *WorkHandler) GetHotWorks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 兼容旧的limit参数
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ := strconv.Atoi(limitStr)
		if limit > 0 {
			pageSize = limit
			page = 1
		}
	}

	works, total, err := h.service.GetHotWorks(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	workResponses := make([]*models.WorkResponse, len(works))
	for i, work := range works {
		workResponses[i] = h.toResponse(work)
	}

	// 如果使用limit参数，返回数组格式（兼容旧接口）
	if c.Query("limit") != "" {
		utils.Success(c, workResponses)
	} else {
		// 使用分页格式
		utils.PageResponse(c, workResponses, total, page, pageSize)
	}
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

// GetUserWorks 获取指定用户的作品列表（公开访问）
// GET /api/users/:id/works
func (h *WorkHandler) GetUserWorks(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	workType := c.DefaultQuery("type", "all")
	sortBy := c.DefaultQuery("sort_by", "latest") // 默认最新排序

	works, total, err := h.service.GetUserWorks(uint(userID), page, pageSize, workType, sortBy)
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
