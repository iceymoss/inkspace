package handler

import (
	"errors"
	"strconv"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

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

	article, err := h.service.Update(uint(id), &req, userID.(uint), roleStr)
	if err != nil {
		// 权限错误返回403
		if err.Error() == "无权限修改" {
			utils.Error(c, 403, err.Error())
		} else {
			utils.Error(c, 400, err.Error())
		}
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
		if err.Error() == "无权限删除" {
			utils.Error(c, 403, err.Error())
		} else {
			utils.Error(c, 400, err.Error())
		}
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

	// 权限检查：如果是草稿(status=0)或私有(status=2)，只有作者或管理员可以查看
	if article.Status == 0 || article.Status == 2 {
		userID, exists := c.Get("user_id")
		role, _ := c.Get("role")
		roleStr := "user"
		if role != nil {
			roleStr = role.(string)
		}

		// 如果不是管理员，且不是作者，则无权限查看
		if roleStr != "admin" && (!exists || userID.(uint) != article.AuthorID) {
			utils.Error(c, 403, "无权限查看此文章")
			return
		}
	}

	// Increment view count
	go h.service.IncrementViewCount(uint(id))

	utils.Success(c, article.ToResponse())
}

// GetEdit 获取文章详情用于编辑（需要认证，且只允许作者或管理员访问）
// GET /api/articles/:id/edit
func (h *ArticleHandler) GetEdit(c *gin.Context) {
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

	// 查询文章，但使用WHERE条件确保权限（非管理员只能查询自己的文章）
	var article models.Article
	query := database.DB.Where("id = ?", uint(id))

	// 非管理员只能查询自己的文章
	if roleStr != "admin" {
		query = query.Where("author_id = ?", userID.(uint))
	}

	if err := query.Preload("Category").Preload("Tags").Preload("Author").First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 403, "文章不存在或无权限编辑")
		} else {
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	// 返回文章数据（不增加浏览量）
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
		query.PageSize = 20
	}

	// 检查是否是管理后台请求（通过路径判断）
	// 如果是 /api/admin/articles，则显示所有状态
	if c.Request.URL.Path == "/api/admin/articles" {
		query.ShowAll = true
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

// GetUserArticles 获取用户的文章列表
// GET /api/users/:id/articles
// 数据安全：只返回公开的文章（status=1），不返回草稿和私有文章
// 支持排序参数：sort_by=latest（最新）或 sort_by=hot（最热）
func (h *ArticleHandler) GetUserArticles(c *gin.Context) {
	authorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var query models.ArticleListQuery
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

	// 获取当前登录用户ID（可能不存在）
	currentUserID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		currentUserID = uid.(uint)
	}

	// 设置作者ID过滤
	query.AuthorID = uint(authorID)

	// 数据安全：如果查看的是自己的文章，显示所有状态；如果是别人的文章，只显示公开的
	if currentUserID == uint(authorID) {
		// 查看自己的文章，显示所有状态（草稿、私有、公开）
		// 不设置Status，让GetList显示所有状态
	} else {
		// 查看别人的文章，只显示公开的（status=1）
		status := 1
		query.Status = &status
	}

	// 处理排序参数：latest -> time, hot -> hot
	if query.SortBy == "latest" {
		query.SortBy = "time"
	}
	// 如果没有指定排序，默认使用最新排序
	if query.SortBy == "" {
		query.SortBy = "time"
	}
	query.SortOrder = "desc"

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

// GetRecommended 获取推荐文章
// GET /api/articles/recommended
func (h *ArticleHandler) GetRecommended(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}
	if limit > 10 {
		limit = 10
	}

	articles, err := h.service.GetRecommended(limit)
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

	utils.Success(c, articleResponses)
}

// SetRecommend 设置文章推荐状态
// PUT /api/admin/articles/:id/recommend
func (h *ArticleHandler) SetRecommend(c *gin.Context) {
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

// GetHotArticles 获取热门文章
// GET /api/articles/hot
func (h *ArticleHandler) GetHotArticles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "6")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 6
	}
	if limit > 20 {
		limit = 20
	}

	articles, err := h.service.GetHotArticles(limit)
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

	utils.Success(c, articleResponses)
}
