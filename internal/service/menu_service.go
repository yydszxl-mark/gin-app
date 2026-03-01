package service

import (
	"context"
	"gin-app-start/internal/dto"
	"gin-app-start/internal/model"
	"gin-app-start/internal/repository"
	"gin-app-start/pkg/errors"
	"gin-app-start/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MenuService interface {
	CreateMenu(ctx context.Context, req *dto.CreateMenuRequest) (*model.Menu, error)
	GetMenu(ctx context.Context, id uint) (*model.Menu, error)
	UpdateMenu(ctx context.Context, id uint, req *dto.UpdateMenuRequest) (*model.Menu, error)
	DeleteMenu(ctx context.Context, id uint) error
	ListMenus(ctx context.Context) ([]model.Menu, error)
	GetMenuTree(ctx context.Context) ([]dto.MenuResponse, error)
	GetUserMenuTree(ctx context.Context, userID uint) ([]dto.MenuResponse, error)
}

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (s *menuService) CreateMenu(ctx context.Context, req *dto.CreateMenuRequest) (*model.Menu, error) {
	menu := &model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Title:     req.Title,
		Icon:      req.Icon,
		Path:      req.Path,
		Component: req.Component,
		Redirect:  req.Redirect,
		Type:      req.Type,
		Hidden:    req.Hidden,
		Sort:      req.Sort,
		Status:    req.Status,
	}

	if err := s.menuRepo.Create(ctx, menu); err != nil {
		logger.Error("Failed to create menu", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "创建菜单失败", err)
	}

	logger.Info("Menu created successfully", zap.Uint("menu_id", menu.ID))
	return menu, nil
}

func (s *menuService) GetMenu(ctx context.Context, id uint) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "菜单不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询菜单失败", err)
	}
	return menu, nil
}

func (s *menuService) UpdateMenu(ctx context.Context, id uint, req *dto.UpdateMenuRequest) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "菜单不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询菜单失败", err)
	}

	menu.ParentID = req.ParentID
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Title != "" {
		menu.Title = req.Title
	}
	menu.Icon = req.Icon
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Redirect = req.Redirect
	if req.Type != "" {
		menu.Type = req.Type
	}
	menu.Hidden = req.Hidden
	menu.Sort = req.Sort
	if req.Status != 0 {
		menu.Status = req.Status
	}

	if err := s.menuRepo.Update(ctx, menu); err != nil {
		logger.Error("Failed to update menu", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "更新菜单失败", err)
	}

	logger.Info("Menu updated successfully", zap.Uint("menu_id", menu.ID))
	return menu, nil
}

func (s *menuService) DeleteMenu(ctx context.Context, id uint) error {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "菜单不存在")
		}
		return errors.WrapBusinessError(10005, "查询菜单失败", err)
	}

	// 检查是否有子菜单
	children, err := s.menuRepo.GetByParentID(ctx, id)
	if err != nil {
		return errors.WrapBusinessError(10005, "查询子菜单失败", err)
	}
	if len(children) > 0 {
		return errors.NewBusinessError(10001, "存在子菜单，无法删除")
	}

	if err := s.menuRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete menu", zap.Error(err))
		return errors.WrapBusinessError(10005, "删除菜单失败", err)
	}

	logger.Info("Menu deleted successfully", zap.Uint("menu_id", menu.ID))
	return nil
}

func (s *menuService) ListMenus(ctx context.Context) ([]model.Menu, error) {
	menus, err := s.menuRepo.List(ctx)
	if err != nil {
		logger.Error("Failed to list menus", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "查询菜单列表失败", err)
	}
	return menus, nil
}

func (s *menuService) GetMenuTree(ctx context.Context) ([]dto.MenuResponse, error) {
	menus, err := s.menuRepo.GetMenuTree(ctx)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "查询菜单树失败", err)
	}

	return s.buildMenuTree(menus, 0), nil
}

func (s *menuService) GetUserMenuTree(ctx context.Context, userID uint) ([]dto.MenuResponse, error) {
	menus, err := s.menuRepo.GetUserMenus(ctx, userID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "查询用户菜单失败", err)
	}

	return s.buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func (s *menuService) buildMenuTree(menus []model.Menu, parentID uint) []dto.MenuResponse {
	var tree []dto.MenuResponse

	for _, menu := range menus {
		if menu.ParentID == parentID {
			menuResp := dto.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Title:     menu.Title,
				Icon:      menu.Icon,
				Path:      menu.Path,
				Component: menu.Component,
				Redirect:  menu.Redirect,
				Type:      menu.Type,
				Hidden:    menu.Hidden,
				Sort:      menu.Sort,
				Status:    menu.Status,
				CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: menu.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			// 递归查找子菜单
			children := s.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menuResp.Children = children
			}

			tree = append(tree, menuResp)
		}
	}

	return tree
}
