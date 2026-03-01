package middleware

import (
	"gin-app-start/internal/repository"
	"gin-app-start/pkg/response"

	"github.com/gin-gonic/gin"
)

// PermissionCheck 权限检查中间件
func PermissionCheck(permRepo repository.PermissionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Error(c, 10003, "未授权访问")
			c.Abort()
			return
		}

		// 获取请求的方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path

		// 检查用户是否有权限访问该接口
		hasPermission, err := permRepo.CheckUserPermission(c.Request.Context(), userID, method, path)
		if err != nil {
			response.Error(c, 50000, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Error(c, 10003, "无权限访问该资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleCheck 角色检查中间件
func RoleCheck(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDs := GetRoleIDs(c)
		if len(roleIDs) == 0 {
			response.Error(c, 10003, "未分配角色")
			c.Abort()
			return
		}

		// 这里简化处理，实际应该查询角色编码
		// 可以根据需要扩展
		c.Next()
	}
}
