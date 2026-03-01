package controller

import "gin-app-start/internal/service"

// ControllerGroup 控制器层组
type ControllerGroup struct {
	HealthController     *HealthController
	UserController       *UserController
	DeviceController     *DeviceController
	AuthController       *AuthController
	RoleController       *RoleController
	PermissionController *PermissionController
	MenuController       *MenuController
}

// NewControllerGroup 创建控制器层组
func NewControllerGroup(svc *service.ServiceGroup) *ControllerGroup {
	return &ControllerGroup{
		HealthController:     NewHealthController(),
		UserController:       NewUserController(svc.UserService),
		DeviceController:     NewDeviceController(svc.DeviceService),
		AuthController:       NewAuthController(svc.AuthService),
		RoleController:       NewRoleController(svc.RoleService),
		PermissionController: NewPermissionController(svc.PermissionService),
		MenuController:       NewMenuController(svc.MenuService),
	}
}
