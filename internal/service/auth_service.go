package service

import (
	"context"
	"gin-app-start/internal/dto"
	"gin-app-start/internal/repository"
	"gin-app-start/pkg/errors"
	"gin-app-start/pkg/jwt"
	"gin-app-start/pkg/logger"
	"gin-app-start/pkg/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	GetUserInfo(ctx context.Context, userID uint) (*dto.UserInfo, error)
	ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "用户名或密码错误")
		}
		logger.Error("Failed to query user", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "查询用户失败", err)
	}

	// 验证密码
	if !utils.ValidatePassword(req.Password, user.Salt, user.Password) {
		return nil, errors.NewBusinessError(10002, "用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.NewBusinessError(10003, "用户已被禁用")
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		logger.Error("Failed to get user roles", zap.Error(err))
		return nil, errors.WrapBusinessError(10005, "获取用户角色失败", err)
	}

	// 提取角色ID
	roleIDs := make([]uint, len(roles))
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
		roleNames[i] = role.Name
	}

	// 生成 Token (7天有效期)
	token, err := jwt.GenerateToken(user.ID, user.Username, roleIDs, 7*24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.WrapBusinessError(50000, "生成Token失败", err)
	}

	logger.Info("User login successfully",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
	)

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		UserInfo: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Status:   user.Status,
			Roles:    roleNames,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// 解析旧 Token
	claims, err := jwt.ParseToken(req.Token)
	if err != nil {
		return nil, errors.NewBusinessError(10003, "Token无效或已过期")
	}

	// 查询用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "用户不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询用户失败", err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.NewBusinessError(10003, "用户已被禁用")
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "获取用户角色失败", err)
	}

	roleIDs := make([]uint, len(roles))
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
		roleNames[i] = role.Name
	}

	// 生成新 Token
	newToken, err := jwt.GenerateToken(user.ID, user.Username, roleIDs, 7*24*time.Hour)
	if err != nil {
		return nil, errors.WrapBusinessError(50000, "生成Token失败", err)
	}

	return &dto.LoginResponse{
		Token:     newToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		UserInfo: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Status:   user.Status,
			Roles:    roleNames,
		},
	}, nil
}

func (s *authService) GetUserInfo(ctx context.Context, userID uint) (*dto.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(10002, "用户不存在")
		}
		return nil, errors.WrapBusinessError(10005, "查询用户失败", err)
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, errors.WrapBusinessError(10005, "获取用户角色失败", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Status:   user.Status,
		Roles:    roleNames,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBusinessError(10002, "用户不存在")
		}
		return errors.WrapBusinessError(10005, "查询用户失败", err)
	}

	// 验证旧密码
	if !utils.ValidatePassword(req.OldPassword, user.Salt, user.Password) {
		return errors.NewBusinessError(10001, "原密码错误")
	}

	// 生成新密码
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		logger.Error("Failed to generate salt", zap.Error(err))
		return errors.WrapBusinessError(50000, "生成盐值失败", err)
	}
	hashedPassword := utils.HashPassword(req.NewPassword, salt)

	user.Password = hashedPassword
	user.Salt = salt

	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update password", zap.Error(err))
		return errors.WrapBusinessError(10005, "修改密码失败", err)
	}

	logger.Info("User password changed",
		zap.Uint("user_id", userID),
	)

	return nil
}
