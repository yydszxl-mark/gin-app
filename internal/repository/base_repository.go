package repository

import (
	"context"
	"gin-app-start/pkg/database"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return database.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := database.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return database.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(new(T), id).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, offset, limit int) ([]*T, error) {
	var entities []*T
	err := database.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}

func (r *BaseRepository[T]) GetDB() *gorm.DB {
	return database.DB
}
