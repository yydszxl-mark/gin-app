package router

import (
	"gin-app-start/internal/config"
	"gin-app-start/internal/controller"
	"gin-app-start/internal/middleware"
	"gin-app-start/internal/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB, ctrl *controller.ControllerGroup) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	if cfg.Server.LimitNum > 0 {
		router.Use(middleware.RateLimit(cfg.Server.LimitNum))
	}

	// 初始化 Repository（用于权限检查中间件）
	permRepo := repository.NewPermissionRepository(db)

	router.GET("/health", ctrl.HealthController.HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 认证相关路由（无需认证）
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", ctrl.AuthController.Login)
		auth.POST("/refresh", ctrl.AuthController.RefreshToken)
	}

	// API v1 路由组（需要认证）
	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.JWTAuth())
	{
		// 当前用户信息（无需额外权限）
		apiV1.GET("/user/info", ctrl.AuthController.GetUserInfo)
		apiV1.POST("/user/change-password", ctrl.AuthController.ChangePassword)
		apiV1.POST("/user/logout", ctrl.AuthController.Logout)

		// 用户管理（需要权限检查）
		users := apiV1.Group("/users")
		users.Use(middleware.PermissionCheck(permRepo)) // 启用权限检查
		{
			users.POST("", ctrl.UserController.CreateUser)
			users.GET("/:id", ctrl.UserController.GetUser)
			users.PUT("/:id", ctrl.UserController.UpdateUser)
			users.DELETE("/:id", ctrl.UserController.DeleteUser)
			users.GET("", ctrl.UserController.ListUsers)
			users.POST("/:id/roles", ctrl.UserController.AssignRoles)
		}

		// 角色管理（需要权限检查）
		roles := apiV1.Group("/roles")
		roles.Use(middleware.PermissionCheck(permRepo)) // 启用权限检查
		{
			roles.POST("", ctrl.RoleController.CreateRole)
			roles.GET("/:id", ctrl.RoleController.GetRole)
			roles.PUT("/:id", ctrl.RoleController.UpdateRole)
			roles.DELETE("/:id", ctrl.RoleController.DeleteRole)
			roles.GET("", ctrl.RoleController.ListRoles)
			roles.POST("/:id/permissions", ctrl.RoleController.AssignPermissions)
			roles.POST("/:id/menus", ctrl.RoleController.AssignMenus)
			roles.GET("/:id/permissions", ctrl.RoleController.GetRolePermissions)
			roles.GET("/:id/menus", ctrl.RoleController.GetRoleMenus)
		}

		// 权限管理（需要权限检查）
		permissions := apiV1.Group("/permissions")
		permissions.Use(middleware.PermissionCheck(permRepo)) // 启用权限检查
		{
			permissions.POST("", ctrl.PermissionController.CreatePermission)
			permissions.GET("/:id", ctrl.PermissionController.GetPermission)
			permissions.PUT("/:id", ctrl.PermissionController.UpdatePermission)
			permissions.DELETE("/:id", ctrl.PermissionController.DeletePermission)
			permissions.GET("", ctrl.PermissionController.ListPermissions)
			permissions.GET("/user", ctrl.PermissionController.GetUserPermissions)
		}

		// 菜单管理（需要权限检查）
		menus := apiV1.Group("/menus")
		menus.Use(middleware.PermissionCheck(permRepo)) // 启用权限检查
		{
			menus.POST("", ctrl.MenuController.CreateMenu)
			menus.GET("/:id", ctrl.MenuController.GetMenu)
			menus.PUT("/:id", ctrl.MenuController.UpdateMenu)
			menus.DELETE("/:id", ctrl.MenuController.DeleteMenu)
			menus.GET("", ctrl.MenuController.ListMenus)
			menus.GET("/tree", ctrl.MenuController.GetMenuTree)
			menus.GET("/user/tree", ctrl.MenuController.GetUserMenuTree)
		}

		// 设备管理（仅需要认证，不需要权限检查）
		devices := apiV1.Group("/devices")
		{
			devices.POST("", ctrl.DeviceController.CreateDevice)
			devices.GET("/:id", ctrl.DeviceController.GetDevice)
			devices.PUT("/:id", ctrl.DeviceController.UpdateDevice)
			devices.DELETE("/:id", ctrl.DeviceController.DeleteDevice)
			devices.GET("", ctrl.DeviceController.ListDevices)
		}
	}

	return router
}
