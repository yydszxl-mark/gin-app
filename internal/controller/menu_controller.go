package controller

import (
	"gin-app-start/internal/dto"
	"gin-app-start/internal/middleware"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService service.MenuService
}

func NewMenuController(menuService service.MenuService) *MenuController {
	return &MenuController{
		menuService: menuService,
	}
}

// CreateMenu 创建菜单
func (ctrl *MenuController) CreateMenu(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	menu, err := ctrl.menuService.CreateMenu(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menu)
}

// GetMenu 获取菜单详情
func (ctrl *MenuController) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的菜单ID")
		return
	}

	menu, err := ctrl.menuService.GetMenu(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menu)
}

// UpdateMenu 更新菜单
func (ctrl *MenuController) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的菜单ID")
		return
	}

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	menu, err := ctrl.menuService.UpdateMenu(c.Request.Context(), uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menu)
}

// DeleteMenu 删除菜单
func (ctrl *MenuController) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 10001, "无效的菜单ID")
		return
	}

	if err := ctrl.menuService.DeleteMenu(c.Request.Context(), uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ListMenus 菜单列表
func (ctrl *MenuController) ListMenus(c *gin.Context) {
	menus, err := ctrl.menuService.ListMenus(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menus)
}

// GetMenuTree 获取菜单树
func (ctrl *MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := ctrl.menuService.GetMenuTree(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menuTree)
}

// GetUserMenuTree 获取当前用户菜单树
func (ctrl *MenuController) GetUserMenuTree(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, 10003, "未授权访问")
		return
	}

	menuTree, err := ctrl.menuService.GetUserMenuTree(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, menuTree)
}
