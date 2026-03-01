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

type RoleService interface {
	CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*model.Role, error)
	GetRole(ctx context.Context, id uint) (*model.Role, error)
	UpdateRole(ctx context.Context, id uint, req *dto.UpdateRoleRequest) (*model.Role, error)
	DeleteRole(ctx context.Context, id uint) error
	ListRoles(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error)
	AssignPermissions(ctx context.Context, roleID uint, req *dto.AssignPermissionsRequest) error
	AssignMenus(ctx context.Context, roleID uint, req *dto.AssignMenusRequest) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error)
	GetRoleMenus(ctx context.Context, roleID uint) ([]model.Menu, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (s *roleService) CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*model.Role, error) {
	// 检查角色编码是否已存在
	existingRole, err := s.roleRepo.GetByCode(ctx, req.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.WrapBusinessError(10005, "查询角色失败", err)
	}
	if existingRole != nil {
		return nil, errors.NewBusinessError(10004, "角色编码已存在")
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		logger.Error("Failed to create role", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "创建角色失败", err)
	}

	logger.Info("Role created successfully", zap.Uint("role_id", role.ID))
	return role, nil
}

func (s *roleService) GetRole(ctx context.Context, id uint) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "角色不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询角色失败", err)
	}
	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, req *dto.UpdateRoleRequest) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "角色不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询角色失败", err)
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.Sort = req.Sort
	if req.Status != 0 {
		role.Status = req.Status
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		logger.Error("Failed to update role", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "更新角色失败", err)
	}

	logger.Info("Role updated successfully", zap.Uint("role_id", role.ID))
	return role, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "角色不存在")
		}
		return errors.WrapBusinessError(10005, "查询角色失败", err)
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete role", zap.Error(err))
		return errors.WrapBusinessError(10005, "删除角色失败", err)
	}

	logger.Info("Role deleted successfully", zap.Uint("role_id", role.ID))
	return nil
}

func (s *roleService) ListRoles(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	roles, total, err := s.roleRepo.List(ctx, offset, pageSize)
	if err != nil {
		logger.Error("Failed to list roles", zap.Error(err))
		return nil, 0, errors.WrapBusinessError(10005, "查询角色列表失败", err)
	}

	return roles, total, nil
}

func (s *roleService) AssignPermissions(ctx context.Context, roleID uint, req *dto.AssignPermissionsRequest) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "角色不存在")
		}
		return errors.WrapBusinessError(10005, "查询角色失败", err)
	}

	if err := s.roleRepo.AssignPermissions(ctx, roleID, req.PermissionIDs); err != nil {
		logger.Error("Failed to assign permissions", zap.Error(err))
		return errors.WrapBusinessError(10005, "分配权限失败", err)
	}

	logger.Info("Permissions assigned successfully", zap.Uint("role_id", roleID))
	return nil
}

func (s *roleService) AssignMenus(ctx context.Context, roleID uint, req *dto.AssignMenusRequest) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "角色不存在")
		}
		return errors.WrapBusinessError(10005, "查询角色失败", err)
	}

	if err := s.roleRepo.AssignMenus(ctx, roleID, req.MenuIDs); err != nil {
		logger.Error("Failed to assign menus", zap.Error(err))
		return errors.WrapBusinessError(10005, "分配菜单失败", err)
	}

	logger.Info("Menus assigned successfully", zap.Uint("role_id", roleID))
	return nil
}

func (s *roleService) GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	permissions, err := s.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "查询角色权限失败", err)
	}
	return permissions, nil
}

func (s *roleService) GetRoleMenus(ctx context.Context, roleID uint) ([]model.Menu, error) {
	menus, err := s.roleRepo.GetRoleMenus(ctx, roleID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "查询角色菜单失败", err)
	}
	return menus, nil
}
