package repository

import (
	"context"
	"gin-app-start/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id uint) (*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*model.Role, int64, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error)
	GetRoleMenus(ctx context.Context, roleID uint) ([]model.Menu, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	AssignMenus(ctx context.Context, roleID uint, menuIDs []uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error)
}

type roleRepository struct {
	*BaseRepository[model.Role]
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: &BaseRepository[model.Role]{},
	}
}

func (r *roleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.GetDB().WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.GetDB().WithContext(ctx).Preload("Permissions").Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context, offset, limit int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	if err := r.GetDB().WithContext(ctx).Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.GetDB().WithContext(ctx).
		Order("sort ASC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&roles).Error

	return roles, total, err
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var role model.Role
	err := r.GetDB().WithContext(ctx).Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (r *roleRepository) GetRoleMenus(ctx context.Context, roleID uint) ([]model.Menu, error) {
	var role model.Role
	err := r.GetDB().WithContext(ctx).Preload("Menus").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Menus, nil
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var role model.Role
	if err := r.GetDB().WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}

	// 清空现有权限
	if err := r.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}

	// 分配新权限
	if len(permissionIDs) > 0 {
		var permissions []model.Permission
		if err := r.GetDB().WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}
		return r.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Append(permissions)
	}

	return nil
}

func (r *roleRepository) AssignMenus(ctx context.Context, roleID uint, menuIDs []uint) error {
	var role model.Role
	if err := r.GetDB().WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}

	// 清空现有菜单
	if err := r.GetDB().WithContext(ctx).Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}

	// 分配新菜单
	if len(menuIDs) > 0 {
		var menus []model.Menu
		if err := r.GetDB().WithContext(ctx).Find(&menus, menuIDs).Error; err != nil {
			return err
		}
		return r.GetDB().WithContext(ctx).Model(&role).Association("Menus").Append(menus)
	}

	return nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error) {
	var user model.User
	err := r.GetDB().WithContext(ctx).Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}
