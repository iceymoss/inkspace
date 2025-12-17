package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	token, user, err := h.service.Login(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user.ToResponse())
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	user, err := h.service.UpdateUser(userID.(uint), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, user.ToResponse())
}

// ChangePassword 修改密码
// PUT /api/profile/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req models.PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	var query models.UserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("绑定用户查询参数失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	users, total, err := h.service.GetUserList(&query)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 管理后台用户列表需要包含角色、状态等完整信息，使用 UserResponse
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	utils.PageResponse(c, userResponses, total, query.Page, query.PageSize)
}

// SearchUsers 根据关键字搜索用户（用户名或昵称）
// GET /api/users/search?keyword=xxx&limit=10
func (h *UserHandler) SearchUsers(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		utils.BadRequest(c, "keyword不能为空")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	users, err := h.service.SearchUsers(keyword, limit)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	// 将精确匹配的用户放入 top，其余放入 users
	top := make([]*models.PublicUserResponse, 0)
	others := make([]*models.PublicUserResponse, 0)

	for _, user := range users {
		resp := user.ToPublicResponse()
		if user.Username == keyword || user.Nickname == keyword {
			top = append(top, resp)
		} else {
			others = append(others, resp)
		}
	}

	utils.Success(c, gin.H{
		"top":   top,
		"users": others,
	})
}

// UpdateUserStatus 更新用户状态
// PUT /api/admin/users/:id/status
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateUserStatus(uint(id), req.Status); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

// UpdateUserRole 更新用户角色
// PUT /api/admin/users/:id/role
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateUserRole(uint(id), req.Role); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

// DeleteUser 删除用户
// DELETE /api/admin/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetUserProfile 获取用户主页信息（公开API，只返回公开信息）
// GET /api/users/:id
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 只返回公开信息，不包含Email、Role、Status等敏感信息
	utils.Success(c, user.ToPublicResponse())
}
