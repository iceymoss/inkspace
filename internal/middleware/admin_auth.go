package middleware

import (
	"github.com/iceymoss/inkspace/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "未登录")
			c.Abort()
			return
		}

		// 验证 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, 401, "Token格式错误")
			c.Abort()
			return
		}

		// 解析管理员token（使用独立的secret）
		claims, err := utils.ParseAdminToken(parts[1])
		if err != nil {
			utils.Error(c, 401, "Token无效或已过期")
			c.Abort()
			return
		}

		// 验证是否是管理员
		if claims.Role != "admin" {
			utils.Error(c, 403, "需要管理员权限")
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
