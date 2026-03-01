package repository

import (
	"context"
	"gin-app-start/internal/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *model.Permission) error
	GetByID(ctx context.Context, id uint) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*model.Permission, int64, error)
	GetByType(ctx context.Context, permType string) ([]model.Permission, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error)
	CheckUserPermission(ctx context.Context, userID uint, method, path string) (bool, error)
}

type permissionRepository struct {
	*BaseRepository[model.Permission]
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		BaseRepository: &BaseRepository[model.Permission]{},
	}
}

func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.GetDB().WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) List(ctx context.Context, offset, limit int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	if err := r.GetDB().WithContext(ctx).Model(&model.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.GetDB().WithContext(ctx).
		Order("sort ASC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&permissions).Error

	return permissions, total, err
}

func (r *permissionRepository) GetByType(ctx context.Context, permType string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.GetDB().WithContext(ctx).Where("type = ? AND status = 1", permType).Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	// 通过用户的角色获取权限
	err := r.GetDB().WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1", userID).
		Group("permissions.id").
		Find(&permissions).Error

	return permissions, err
}

func (r *permissionRepository) CheckUserPermission(ctx context.Context, userID uint, method, path string) (bool, error) {
	var count int64

	err := r.GetDB().WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.method = ? AND permissions.path = ? AND permissions.status = 1", userID, method, path).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
