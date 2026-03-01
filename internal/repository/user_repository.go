package repository

import (
	"context"
	"gin-app-start/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error
	GetUserWithRoles(ctx context.Context, userID uint) (*model.User, error)
}

type userRepository struct {
	*BaseRepository[model.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: &BaseRepository[model.User]{},
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.GetDB().WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.GetDB().WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.GetDB().WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.GetDB().WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var user model.User
	if err := r.GetDB().WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	// 清空现有角色
	if err := r.GetDB().WithContext(ctx).Model(&user).Association("Roles").Clear(); err != nil {
		return err
	}

	// 分配新角色
	if len(roleIDs) > 0 {
		var roles []model.Role
		if err := r.GetDB().WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
			return err
		}
		return r.GetDB().WithContext(ctx).Model(&user).Association("Roles").Append(roles)
	}

	return nil
}

func (r *userRepository) GetUserWithRoles(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	err := r.GetDB().WithContext(ctx).Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
