package repository

import (
	"context"
	"gin-app-start/internal/model"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, id uint) (*model.Device, error)
	GetByName(ctx context.Context, name string) (*model.Device, error)
	GetByType(ctx context.Context, deviceType string) ([]*model.Device, error)
	Update(ctx context.Context, device *model.Device) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*model.Device, int64, error)
}

type deviceRepository struct {
	*BaseRepository[model.Device]
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{
		BaseRepository: &BaseRepository[model.Device]{},
	}
}

func (r *deviceRepository) GetByName(ctx context.Context, name string) (*model.Device, error) {
	var device model.Device
	err := r.GetDB().WithContext(ctx).Where("name = ?", name).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetByType(ctx context.Context, deviceType string) ([]*model.Device, error) {
	var devices []*model.Device
	err := r.GetDB().WithContext(ctx).Where("type = ?", deviceType).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) List(ctx context.Context, offset, limit int) ([]*model.Device, int64, error) {
	var devices []*model.Device
	var total int64

	if err := r.GetDB().WithContext(ctx).Model(&model.Device{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.GetDB().WithContext(ctx).Offset(offset).Limit(limit).Find(&devices).Error
	return devices, total, err
}
