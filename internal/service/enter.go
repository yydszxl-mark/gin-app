package service

import "gin-app-start/internal/repository"

// ServiceGroup 服务层组
type ServiceGroup struct {
	UserService       UserService
	DeviceService     DeviceService
	AuthService       AuthService
	RoleService       RoleService
	PermissionService PermissionService
	MenuService       MenuService
}

// NewServiceGroup 创建服务层组
func NewServiceGroup(repo *repository.RepositoryGroup) *ServiceGroup {
	return &ServiceGroup{
		UserService:       NewUserService(repo.UserRepo),
		DeviceService:     NewDeviceService(repo.DeviceRepo),
		AuthService:       NewAuthService(repo.UserRepo, repo.RoleRepo),
		RoleService:       NewRoleService(repo.RoleRepo),
		PermissionService: NewPermissionService(repo.PermissionRepo),
		MenuService:       NewMenuService(repo.MenuRepo),
	}
}
