package middleware

import (
	"gin-app-start/pkg/jwt"
	"gin-app-start/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 10003, "未提供认证令牌")
			c.Abort()
			return
		}

		// 检查 Token 格式 (Bearer token)
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, 10003, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Error(c, 10003, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_ids", claims.RoleIDs)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			return uid
		}
	}
	return 0
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// GetRoleIDs 从上下文中获取角色ID列表
func GetRoleIDs(c *gin.Context) []uint {
	if roleIDs, exists := c.Get("role_ids"); exists {
		if ids, ok := roleIDs.([]uint); ok {
			return ids
		}
	}
	return []uint{}
}
