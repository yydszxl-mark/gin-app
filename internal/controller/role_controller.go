package controller

import (
	"gin-app-start/internal/dto"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService service.RoleService
}

func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// CreateRole 创建角色
func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	role, err := ctrl.roleService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
func (ctrl *RoleController) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	role, err := ctrl.roleService.GetRole(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	role, err := ctrl.roleService.UpdateRole(c.Request.Context(), uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	if err := ctrl.roleService.DeleteRole(c.Request.Context(), uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ListRoles 角色列表
func (ctrl *RoleController) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	roles, total, err := ctrl.roleService.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.SuccessWithPage(c, roles, total, page, pageSize)
}

// AssignPermissions 分配权限
func (ctrl *RoleController) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	var req dto.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	if err := ctrl.roleService.AssignPermissions(c.Request.Context(), uint(id), &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "分配权限成功"})
}

// AssignMenus 分配菜单
func (ctrl *RoleController) AssignMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	var req dto.AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	if err := ctrl.roleService.AssignMenus(c.Request.Context(), uint(id), &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "分配菜单成功"})
}

// GetRolePermissions 获取角色权限
func (ctrl *RoleController) GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	permissions, err := ctrl.roleService.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, permissions)
}

// GetRoleMenus 获取角色菜单
func (ctrl *RoleController) GetRoleMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的角色ID")
		return
	}

	menus, err := ctrl.roleService.GetRoleMenus(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menus)
}
