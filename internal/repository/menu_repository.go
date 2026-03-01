package repository

import (
	"context"
	"gin-app-start/internal/model"

	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(ctx context.Context, menu *model.Menu) error
	GetByID(ctx context.Context, id uint) (*model.Menu, error)
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]model.Menu, error)
	GetByParentID(ctx context.Context, parentID uint) ([]model.Menu, error)
	GetUserMenus(ctx context.Context, userID uint) ([]model.Menu, error)
	GetMenuTree(ctx context.Context) ([]model.Menu, error)
}

type menuRepository struct {
	*BaseRepository[model.Menu]
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		BaseRepository: &BaseRepository[model.Menu]{},
	}
}

func (r *menuRepository) List(ctx context.Context) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.GetDB().WithContext(ctx).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetByParentID(ctx context.Context, parentID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.GetDB().WithContext(ctx).
		Where("parent_id = ? AND status = 1", parentID).
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetUserMenus(ctx context.Context, userID uint) ([]model.Menu, error) {
	var menus []model.Menu

	// 通过用户的角色获取菜单
	err := r.GetDB().WithContext(ctx).
		Table("menus").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Joins("JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND menus.status = 1", userID).
		Group("menus.id").
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error

	return menus, err
}

func (r *menuRepository) GetMenuTree(ctx context.Context) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.GetDB().WithContext(ctx).
		Where("status = 1").
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}
