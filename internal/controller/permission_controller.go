package controller

import (
	"gin-app-start/internal/dto"
	"gin-app-start/internal/middleware"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permService service.PermissionService
}

func NewPermissionController(permService service.PermissionService) *PermissionController {
	return &PermissionController{
		permService: permService,
	}
}

// CreatePermission 创建权限
func (ctrl *PermissionController) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	permission, err := ctrl.permService.CreatePermission(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, permission)
}

// GetPermission 获取权限详情
func (ctrl *PermissionController) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的权限ID")
		return
	}

	permission, err := ctrl.permService.GetPermission(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
func (ctrl *PermissionController) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的权限ID")
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	permission, err := ctrl.permService.UpdatePermission(c.Request.Context(), uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
func (ctrl *PermissionController) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的权限ID")
		return
	}

	if err := ctrl.permService.DeletePermission(c.Request.Context(), uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ListPermissions 权限列表
func (ctrl *PermissionController) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	permissions, total, err := ctrl.permService.ListPermissions(c.Request.Context(), page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.SuccessWithPage(c, permissions, total, page, pageSize)
}

// GetUserPermissions 获取当前用户权限
func (ctrl *PermissionController) GetUserPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, 10003, "未授权访问")
		return
	}

	permissions, err := ctrl.permService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, permissions)
}
