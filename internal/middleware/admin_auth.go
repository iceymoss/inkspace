package middleware

import (
	"strings"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 验证 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 解析管理员token（使用独立的secret）
		claims, err := utils.ParseAdminToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 验证是否是管理员
		if claims.Role != "admin" {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		// JWT 中的角色可能已过期，每次访问后台都核对当前账号。
		var user models.User
		if err := database.DB.Select("id", "username", "role", "status").First(&user, claims.UserID).Error; err != nil {
			utils.Unauthorized(c, "管理员账号已失效")
			c.Abort()
			return
		}
		if user.Role != "admin" || user.Status != 1 {
			utils.Forbidden(c, "需要有效的管理员权限")
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}
