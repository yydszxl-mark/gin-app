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

type DeviceService interface {
	CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) (*model.Device, error)
	GetDevice(ctx context.Context, id uint) (*model.Device, error)
	GetDeviceByName(ctx context.Context, name string) (*model.Device, error)
	GetDevicesByType(ctx context.Context, deviceType string) ([]*model.Device, error)
	UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) (*model.Device, error)
	DeleteDevice(ctx context.Context, id uint) error
	ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error)
}

type deviceService struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceService(deviceRepo repository.DeviceRepository) DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
	}
}

func (s *deviceService) CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) (*model.Device, error) {
	existingDevice, err := s.deviceRepo.GetByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("Failed to query device", zap.Error(err), zap.String("name", req.Name))
		return nil, errors.WrapBusinessError(20010, "Failed to query device", err)
	}

	if existingDevice != nil {
		return nil, errors.NewBusinessError(20011, "Device name already exists")
	}

	device := &model.Device{
		Name:    req.Name,
		Type:    req.Type,
		Content: req.Content,
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		logger.Error("Failed to create device", zap.Error(err), zap.String("name", req.Name))
		return nil, errors.WrapBusinessError(20012, "Failed to create device", err)
	}

	logger.Info("Device created successfully",
		zap.String("name", device.Name),
		zap.Uint("device_id", device.ID),
	)

	return device, nil
}

func (s *deviceService) GetDevice(ctx context.Context, id uint) (*model.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(20013, "Device not found")
		}
		logger.Error("Failed to get device", zap.Error(err), zap.Uint("device_id", id))
		return nil, errors.WrapBusinessError(20014, "Failed to get device", err)
	}
	return device, nil
}

func (s *deviceService) GetDeviceByName(ctx context.Context, name string) (*model.Device, error) {
	device, err := s.deviceRepo.GetByName(ctx, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(20013, "Device not found")
		}
		logger.Error("Failed to get device", zap.Error(err), zap.String("name", name))
		return nil, errors.WrapBusinessError(20015, "Failed to get device", err)
	}
	return device, nil
}

func (s *deviceService) GetDevicesByType(ctx context.Context, deviceType string) ([]*model.Device, error) {
	devices, err := s.deviceRepo.GetByType(ctx, deviceType)
	if err != nil {
		logger.Error("Failed to get devices by type", zap.Error(err), zap.String("type", deviceType))
		return nil, errors.WrapBusinessError(20016, "Failed to get devices by type", err)
	}
	return devices, nil
}

func (s *deviceService) UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) (*model.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(20013, "Device not found")
		}
		return nil, errors.WrapBusinessError(20017, "Failed to get device", err)
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Content != "" {
		device.Content = req.Content
	}

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		logger.Error("Failed to update device", zap.Error(err), zap.Uint("device_id", id))
		return nil, errors.WrapBusinessError(20018, "Failed to update device", err)
	}

	logger.Info("Device updated successfully", zap.Uint("device_id", id))
	return device, nil
}

func (s *deviceService) DeleteDevice(ctx context.Context, id uint) error {
	if err := s.deviceRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete device", zap.Error(err), zap.Uint("device_id", id))
		return errors.WrapBusinessError(20019, "Failed to delete device", err)
	}

	logger.Info("Device deleted successfully", zap.Uint("device_id", id))
	return nil
}

func (s *deviceService) ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	devices, total, err := s.deviceRepo.List(ctx, offset, pageSize)
	if err != nil {
		logger.Error("Failed to get device list", zap.Error(err))
		return nil, 0, errors.WrapBusinessError(20020, "Failed to get device list", err)
	}

	return devices, total, nil
}
