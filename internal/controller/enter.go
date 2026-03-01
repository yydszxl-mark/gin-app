package controller

import "gin-app-start/internal/service"

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	HealthController
	UserController
	DeviceController
}

var (
	userService   = service.ServiceGroupApp.UserService
	deviceService = service.ServiceGroupApp.DeviceService
)
