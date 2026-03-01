package router

import (
	"gin-app-start/internal/config"
	"gin-app-start/internal/controller"
	"gin-app-start/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	if cfg.Server.LimitNum > 0 {
		router.Use(middleware.RateLimit(cfg.Server.LimitNum))
	}
	app := controller.ApiGroupApp
	router.GET("/health", app.HealthController.HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.POST("", app.UserController.CreateUser)
			users.GET("/:id", app.UserController.GetUser)
			users.PUT("/:id", app.UserController.UpdateUser)
			users.DELETE("/:id", app.UserController.DeleteUser)
			users.GET("", app.UserController.ListUsers)
		}

		devices := apiV1.Group("/devices")
		{
			devices.POST("", app.DeviceController.CreateDevice)
			devices.GET("/:id", app.DeviceController.GetDevice)
			devices.PUT("/:id", app.DeviceController.UpdateDevice)
			devices.DELETE("/:id", app.DeviceController.DeleteDevice)
			devices.GET("", app.DeviceController.ListDevices)
		}
	}

	return router
}
