package main

import (
	"context"
	"fmt"
	"gin-app-start/internal/config"
	"gin-app-start/internal/model"
	"gin-app-start/pkg/database"
	"gin-app-start/pkg/logger"
	"gin-app-start/pkg/utils"
	"log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化脚本：创建默认管理员用户、角色、权限和菜单
func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Server.Mode, cfg.Log.FilePath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 连接数据库
	db, err := database.NewPostgresDB(&database.PostgresConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxLifetime:  cfg.Database.MaxLifetime,
		LogLevel:     cfg.Database.LogLevel,
	})
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}

	logger.Info("Starting initialization...")

	ctx := context.Background()

	// 1. 创建管理员角色
	adminRole := createAdminRole(ctx, db)
	logger.Info("Admin role created", zap.Uint("role_id", adminRole.ID))

	// 2. 创建开发者角色
	devRole := createDeveloperRole(ctx, db)
	logger.Info("Developer role created", zap.Uint("role_id", devRole.ID))

	// 3. 创建权限
	permissions := createPermissions(ctx, db)
	logger.Info("Permissions created", zap.Int("count", len(permissions)))

	// 4. 创建菜单
	menus := createMenus(ctx, db)
	logger.Info("Menus created", zap.Int("count", len(menus)))

	// 5. 分配所有权限给管理员角色
	assignPermissionsToRole(ctx, db, adminRole.ID, permissions)
	logger.Info("Permissions assigned to admin role")

	// 6. 分配部分权限给开发者角色
	devPermissions := filterDevPermissions(permissions)
	assignPermissionsToRole(ctx, db, devRole.ID, devPermissions)
	logger.Info("Permissions assigned to developer role")

	// 7. 分配所有菜单给管理员角色
	assignMenusToRole(ctx, db, adminRole.ID, menus)
	logger.Info("Menus assigned to admin role")

	// 8. 分配部分菜单给开发者角色
	devMenus := filterDevMenus(menus)
	assignMenusToRole(ctx, db, devRole.ID, devMenus)
	logger.Info("Menus assigned to developer role")

	// 9. 创建管理员用户
	adminUser := createAdminUser(ctx, db)
	logger.Info("Admin user created", zap.Uint("user_id", adminUser.ID))

	// 10. 分配管理员角色给管理员用户
	assignRolesToUser(ctx, db, adminUser.ID, []uint{adminRole.ID})
	logger.Info("Admin role assigned to admin user")

	logger.Info("Initialization completed successfully!")
	fmt.Println("\n===========================================")
	fmt.Println("初始化完成！")
	fmt.Println("===========================================")
	fmt.Println("管理员账号:")
	fmt.Println("  用户名: admin")
	fmt.Println("  密码: admin123")
	fmt.Println("===========================================")
}

func createAdminRole(ctx context.Context, db *gorm.DB) *model.Role {
	role := &model.Role{
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员，拥有所有权限",
		Sort:        1,
		Status:      1,
	}

	// 检查是否已存在
	var existing model.Role
	if err := db.Where("code = ?", role.Code).First(&existing).Error; err == nil {
		return &existing
	}

	if err := db.Create(role).Error; err != nil {
		logger.Fatal("Failed to create admin role", zap.Error(err))
	}
	return role
}

func createDeveloperRole(ctx context.Context, db *gorm.DB) *model.Role {
	role := &model.Role{
		Name:        "开发者",
		Code:        "developer",
		Description: "开发人员角色",
		Sort:        2,
		Status:      1,
	}

	var existing model.Role
	if err := db.Where("code = ?", role.Code).First(&existing).Error; err == nil {
		return &existing
	}

	if err := db.Create(role).Error; err != nil {
		logger.Fatal("Failed to create developer role", zap.Error(err))
	}
	return role
}

func createPermissions(ctx context.Context, db *gorm.DB) []model.Permission {
	permissions := []model.Permission{
		// 用户管理权限
		{Name: "查看用户列表", Code: "user:list", Type: "api", Method: "GET", Path: "/api/v1/users", Sort: 1, Status: 1},
		{Name: "查看用户详情", Code: "user:get", Type: "api", Method: "GET", Path: "/api/v1/users/:id", Sort: 2, Status: 1},
		{Name: "创建用户", Code: "user:create", Type: "api", Method: "POST", Path: "/api/v1/users", Sort: 3, Status: 1},
		{Name: "更新用户", Code: "user:update", Type: "api", Method: "PUT", Path: "/api/v1/users/:id", Sort: 4, Status: 1},
		{Name: "删除用户", Code: "user:delete", Type: "api", Method: "DELETE", Path: "/api/v1/users/:id", Sort: 5, Status: 1},
		{Name: "分配角色", Code: "user:assign-roles", Type: "api", Method: "POST", Path: "/api/v1/users/:id/roles", Sort: 6, Status: 1},

		// 角色管理权限
		{Name: "查看角色列表", Code: "role:list", Type: "api", Method: "GET", Path: "/api/v1/roles", Sort: 11, Status: 1},
		{Name: "查看角色详情", Code: "role:get", Type: "api", Method: "GET", Path: "/api/v1/roles/:id", Sort: 12, Status: 1},
		{Name: "创建角色", Code: "role:create", Type: "api", Method: "POST", Path: "/api/v1/roles", Sort: 13, Status: 1},
		{Name: "更新角色", Code: "role:update", Type: "api", Method: "PUT", Path: "/api/v1/roles/:id", Sort: 14, Status: 1},
		{Name: "删除角色", Code: "role:delete", Type: "api", Method: "DELETE", Path: "/api/v1/roles/:id", Sort: 15, Status: 1},
		{Name: "分配权限", Code: "role:assign-permissions", Type: "api", Method: "POST", Path: "/api/v1/roles/:id/permissions", Sort: 16, Status: 1},
		{Name: "分配菜单", Code: "role:assign-menus", Type: "api", Method: "POST", Path: "/api/v1/roles/:id/menus", Sort: 17, Status: 1},

		// 权限管理权限
		{Name: "查看权限列表", Code: "permission:list", Type: "api", Method: "GET", Path: "/api/v1/permissions", Sort: 21, Status: 1},
		{Name: "查看权限详情", Code: "permission:get", Type: "api", Method: "GET", Path: "/api/v1/permissions/:id", Sort: 22, Status: 1},
		{Name: "创建权限", Code: "permission:create", Type: "api", Method: "POST", Path: "/api/v1/permissions", Sort: 23, Status: 1},
		{Name: "更新权限", Code: "permission:update", Type: "api", Method: "PUT", Path: "/api/v1/permissions/:id", Sort: 24, Status: 1},
		{Name: "删除权限", Code: "permission:delete", Type: "api", Method: "DELETE", Path: "/api/v1/permissions/:id", Sort: 25, Status: 1},

		// 菜单管理权限
		{Name: "查看菜单列表", Code: "menu:list", Type: "api", Method: "GET", Path: "/api/v1/menus", Sort: 31, Status: 1},
		{Name: "查看菜单详情", Code: "menu:get", Type: "api", Method: "GET", Path: "/api/v1/menus/:id", Sort: 32, Status: 1},
		{Name: "创建菜单", Code: "menu:create", Type: "api", Method: "POST", Path: "/api/v1/menus", Sort: 33, Status: 1},
		{Name: "更新菜单", Code: "menu:update", Type: "api", Method: "PUT", Path: "/api/v1/menus/:id", Sort: 34, Status: 1},
		{Name: "删除菜单", Code: "menu:delete", Type: "api", Method: "DELETE", Path: "/api/v1/menus/:id", Sort: 35, Status: 1},

		// 设备管理权限
		{Name: "查看设备列表", Code: "device:list", Type: "api", Method: "GET", Path: "/api/v1/devices", Sort: 41, Status: 1},
		{Name: "查看设备详情", Code: "device:get", Type: "api", Method: "GET", Path: "/api/v1/devices/:id", Sort: 42, Status: 1},
		{Name: "创建设备", Code: "device:create", Type: "api", Method: "POST", Path: "/api/v1/devices", Sort: 43, Status: 1},
		{Name: "更新设备", Code: "device:update", Type: "api", Method: "PUT", Path: "/api/v1/devices/:id", Sort: 44, Status: 1},
		{Name: "删除设备", Code: "device:delete", Type: "api", Method: "DELETE", Path: "/api/v1/devices/:id", Sort: 45, Status: 1},
	}

	var result []model.Permission
	for _, perm := range permissions {
		var existing model.Permission
		if err := db.Where("code = ?", perm.Code).First(&existing).Error; err == nil {
			result = append(result, existing)
			continue
		}

		if err := db.Create(&perm).Error; err != nil {
			logger.Error("Failed to create permission", zap.Error(err), zap.String("code", perm.Code))
			continue
		}
		result = append(result, perm)
	}

	return result
}

func createMenus(ctx context.Context, db *gorm.DB) []model.Menu {
	menus := []model.Menu{
		// 一级菜单
		{ParentID: 0, Name: "dashboard", Title: "仪表盘", Icon: "dashboard", Path: "/dashboard", Component: "Dashboard", Type: "menu", Sort: 1, Status: 1},
		{ParentID: 0, Name: "system", Title: "系统管理", Icon: "setting", Path: "/system", Component: "Layout", Redirect: "/system/user", Type: "menu", Sort: 2, Status: 1},
		{ParentID: 0, Name: "device", Title: "设备管理", Icon: "laptop", Path: "/device", Component: "Layout", Redirect: "/device/list", Type: "menu", Sort: 3, Status: 1},
	}

	var result []model.Menu
	for _, menu := range menus {
		var existing model.Menu
		if err := db.Where("name = ? AND parent_id = ?", menu.Name, menu.ParentID).First(&existing).Error; err == nil {
			result = append(result, existing)
			continue
		}

		if err := db.Create(&menu).Error; err != nil {
			logger.Error("Failed to create menu", zap.Error(err), zap.String("name", menu.Name))
			continue
		}
		result = append(result, menu)
	}

	// 二级菜单（需要获取父菜单ID）
	var systemMenu model.Menu
	db.Where("name = ?", "system").First(&systemMenu)

	var deviceMenu model.Menu
	db.Where("name = ?", "device").First(&deviceMenu)

	subMenus := []model.Menu{
		{ParentID: systemMenu.ID, Name: "user", Title: "用户管理", Icon: "user", Path: "/system/user", Component: "system/user/index", Type: "menu", Sort: 1, Status: 1},
		{ParentID: systemMenu.ID, Name: "role", Title: "角色管理", Icon: "team", Path: "/system/role", Component: "system/role/index", Type: "menu", Sort: 2, Status: 1},
		{ParentID: systemMenu.ID, Name: "permission", Title: "权限管理", Icon: "lock", Path: "/system/permission", Component: "system/permission/index", Type: "menu", Sort: 3, Status: 1},
		{ParentID: systemMenu.ID, Name: "menu", Title: "菜单管理", Icon: "menu", Path: "/system/menu", Component: "system/menu/index", Type: "menu", Sort: 4, Status: 1},
		{ParentID: deviceMenu.ID, Name: "device-list", Title: "设备列表", Icon: "unordered-list", Path: "/device/list", Component: "device/list/index", Type: "menu", Sort: 1, Status: 1},
	}

	for _, menu := range subMenus {
		var existing model.Menu
		if err := db.Where("name = ? AND parent_id = ?", menu.Name, menu.ParentID).First(&existing).Error; err == nil {
			result = append(result, existing)
			continue
		}

		if err := db.Create(&menu).Error; err != nil {
			logger.Error("Failed to create sub menu", zap.Error(err), zap.String("name", menu.Name))
			continue
		}
		result = append(result, menu)
	}

	return result
}

func createAdminUser(ctx context.Context, db *gorm.DB) *model.User {
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		logger.Fatal("Failed to generate salt", zap.Error(err))
	}
	password := utils.HashPassword("admin123", salt)

	user := &model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Phone:    "13800138000",
		Password: password,
		Salt:     salt,
		Status:   1,
	}

	var existing model.User
	if err := db.Where("username = ?", user.Username).First(&existing).Error; err == nil {
		return &existing
	}

	if err := db.Create(user).Error; err != nil {
		logger.Fatal("Failed to create admin user", zap.Error(err))
	}
	return user
}

func assignPermissionsToRole(ctx context.Context, db *gorm.DB, roleID uint, permissions []model.Permission) {
	var role model.Role
	if err := db.First(&role, roleID).Error; err != nil {
		logger.Fatal("Failed to find role", zap.Error(err))
	}

	if err := db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		logger.Fatal("Failed to assign permissions", zap.Error(err))
	}
}

func assignMenusToRole(ctx context.Context, db *gorm.DB, roleID uint, menus []model.Menu) {
	var role model.Role
	if err := db.First(&role, roleID).Error; err != nil {
		logger.Fatal("Failed to find role", zap.Error(err))
	}

	if err := db.Model(&role).Association("Menus").Replace(menus); err != nil {
		logger.Fatal("Failed to assign menus", zap.Error(err))
	}
}

func assignRolesToUser(ctx context.Context, db *gorm.DB, userID uint, roleIDs []uint) {
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		logger.Fatal("Failed to find user", zap.Error(err))
	}

	var roles []model.Role
	if err := db.Find(&roles, roleIDs).Error; err != nil {
		logger.Fatal("Failed to find roles", zap.Error(err))
	}

	if err := db.Model(&user).Association("Roles").Replace(roles); err != nil {
		logger.Fatal("Failed to assign roles", zap.Error(err))
	}
}

func filterDevPermissions(permissions []model.Permission) []model.Permission {
	var result []model.Permission
	for _, perm := range permissions {
		// 开发者只有查看和设备管理权限
		if perm.Code == "user:list" || perm.Code == "user:get" ||
			perm.Code == "role:list" || perm.Code == "role:get" ||
			perm.Code == "permission:list" || perm.Code == "permission:get" ||
			perm.Code == "menu:list" || perm.Code == "menu:get" ||
			perm.Type == "api" && (perm.Path == "/api/v1/devices" || perm.Path == "/api/v1/devices/:id") {
			result = append(result, perm)
		}
	}
	return result
}

func filterDevMenus(menus []model.Menu) []model.Menu {
	var result []model.Menu
	for _, menu := range menus {
		// 开发者只能看到仪表盘和设备管理
		if menu.Name == "dashboard" || menu.Name == "device" || menu.Name == "device-list" {
			result = append(result, menu)
		}
	}
	return result
}
