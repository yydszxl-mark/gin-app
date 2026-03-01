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

type PermissionService interface {
	CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest) (*model.Permission, error)
	GetPermission(ctx context.Context, id uint) (*model.Permission, error)
	UpdatePermission(ctx context.Context, id uint, req *dto.UpdatePermissionRequest) (*model.Permission, error)
	DeletePermission(ctx context.Context, id uint) error
	ListPermissions(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error)
}

type permissionService struct {
	permRepo repository.PermissionRepository
}

func NewPermissionService(permRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permRepo: permRepo,
	}
}

func (s *permissionService) CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest) (*model.Permission, error) {
	// 检查权限编码是否已存在
	existingPerm, err := s.permRepo.GetByCode(ctx, req.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.WrapBusinessError(10005, "查询权限失败", err)
	}
	if existingPerm != nil {
		return nil, errors.NewBusinessError(10004, "权限编码已存在")
	}

	permission := &model.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Method:      req.Method,
		Path:        req.Path,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := s.permRepo.Create(ctx, permission); err != nil {
		logger.Error("Failed to create permission", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "创建权限失败", err)
	}

	logger.Info("Permission created successfully", zap.Uint("permission_id", permission.ID))
	return permission, nil
}

func (s *permissionService) GetPermission(ctx context.Context, id uint) (*model.Permission, error) {
	permission, err := s.permRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "权限不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询权限失败", err)
	}
	return permission, nil
}

func (s *permissionService) UpdatePermission(ctx context.Context, id uint, req *dto.UpdatePermissionRequest) (*model.Permission, error) {
	permission, err := s.permRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "权限不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询权限失败", err)
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.Method != "" {
		permission.Method = req.Method
	}
	if req.Path != "" {
		permission.Path = req.Path
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	permission.Sort = req.Sort
	if req.Status != 0 {
		permission.Status = req.Status
	}

	if err := s.permRepo.Update(ctx, permission); err != nil {
		logger.Error("Failed to update permission", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "更新权限失败", err)
	}

	logger.Info("Permission updated successfully", zap.Uint("permission_id", permission.ID))
	return permission, nil
}

func (s *permissionService) DeletePermission(ctx context.Context, id uint) error {
	permission, err := s.permRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "权限不存在")
		}
		return errors.WrapBusinessError(10005, "查询权限失败", err)
	}

	if err := s.permRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete permission", zap.Error(err))
		return errors.WrapBusinessError(10005, "删除权限失败", err)
	}

	logger.Info("Permission deleted successfully", zap.Uint("permission_id", permission.ID))
	return nil
}

func (s *permissionService) ListPermissions(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permRepo.List(ctx, offset, pageSize)
	if err != nil {
		logger.Error("Failed to list permissions", zap.Error(err))
		return nil, 0, errors.WrapBusinessError(10005, "查询权限列表失败", err)
	}

	return permissions, total, nil
}

func (s *permissionService) GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error) {
	permissions, err := s.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "查询用户权限失败", err)
	}
	return permissions, nil
}
