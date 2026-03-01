package service

import "gin-app-start/internal/repository"

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	UserService
	DeviceService
}

var (
	userRepository   = repository.RepoGroupApp.UserRepository
	deviceRepository = repository.RepoGroupApp.DeviceRepository
)

func init() {
	ServiceGroupApp.UserService = &userService{}
	ServiceGroupApp.DeviceService = &deviceService{}
}
