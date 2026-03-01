package repository

import "gorm.io/gorm"

// RepositoryGroup 仓储层组
type RepositoryGroup struct {
	UserRepo       UserRepository
	DeviceRepo     DeviceRepository
	RoleRepo       RoleRepository
	PermissionRepo PermissionRepository
	MenuRepo       MenuRepository
}

// NewRepositoryGroup 创建仓储层组
func NewRepositoryGroup(db *gorm.DB) *RepositoryGroup {
	return &RepositoryGroup{
		UserRepo:       NewUserRepository(db),
		DeviceRepo:     NewDeviceRepository(db),
		RoleRepo:       NewRoleRepository(db),
		PermissionRepo: NewPermissionRepository(db),
		MenuRepo:       NewMenuRepository(db),
	}
}
