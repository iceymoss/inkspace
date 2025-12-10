package handler

import (
	"strconv"

	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.GetUserList(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	utils.PageResponse(c, userResponses, total, page, pageSize)
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

