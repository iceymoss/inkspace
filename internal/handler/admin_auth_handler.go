package handler

import (
	"mysite/internal/models"
	"mysite/internal/service"
	"mysite/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminAuthHandler struct {
	service *service.UserService
}

func NewAdminAuthHandler() *AdminAuthHandler {
	return &AdminAuthHandler{
		service: service.NewUserService(),
	}
}

// AdminLogin 管理员登录
// POST /api/admin/auth/login
func (h *AdminAuthHandler) AdminLogin(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 先验证用户
	_, user, err := h.service.Login(&req)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	// 验证是否是管理员
	if user.Role != "admin" {
		utils.Error(c, 403, "该账号不是管理员")
		return
	}

	// 生成管理员专用Token（使用独立的secret）
	adminToken, err := utils.GenerateAdminToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "生成Token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": adminToken,
		"user":  user.ToResponse(),
	})
}

// GetAdminProfile 获取管理员信息
// GET /api/admin/auth/profile
func (h *AdminAuthHandler) GetAdminProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 验证是否是管理员
	if user.Role != "admin" {
		utils.Error(c, 403, "该账号不是管理员")
		return
	}

	utils.Success(c, user.ToResponse())
}

// AdminLogout 管理员退出
// POST /api/admin/auth/logout
func (h *AdminAuthHandler) AdminLogout(c *gin.Context) {
	// 由于使用JWT，服务端不需要特殊处理
	// 客户端清除token即可
	utils.Success(c, nil)
}
