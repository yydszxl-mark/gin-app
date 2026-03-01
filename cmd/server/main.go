package main

import (
	"context"
	"fmt"
	"gin-app-start/internal/config"
	"gin-app-start/internal/controller"
	"gin-app-start/internal/model"
	"gin-app-start/internal/repository"
	"gin-app-start/internal/router"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/database"
	"gin-app-start/pkg/jwt"
	"gin-app-start/pkg/logger"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gin-app-start/docs"

	"go.uber.org/zap"
)

//	@title			Gin App API
//	@version		1.0
//	@description	This is a RESTful API server built with Gin framework.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:9060
//	@BasePath	/

//	@schemes	http https

var Version string

func main() {
	log.Printf("Version: %s\n", Version)

	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	initZapLog(cfg)

	// 3. 初始化数据库
	db := initDB(cfg)

	// 4. 初始化 Redis
	initRedis(cfg)

	// 5. 初始化 JWT
	initJWT()

	// 6. 初始化依赖注入链
	// Repository 层
	repoGroup := repository.NewRepositoryGroup(db)

	// Service 层
	svcGroup := service.NewServiceGroup(repoGroup)

	// Controller 层
	ctrlGroup := controller.NewControllerGroup(svcGroup)

	// 7. 初始化路由
	initRouter(cfg, db, ctrlGroup)
}

func initJWT() {
	// 初始化 JWT 密钥对
	privateKeyPath := "configs/private_key.pem"
	publicKeyPath := "configs/public_key.pem"

	if err := jwt.InitJWT(privateKeyPath, publicKeyPath); err != nil {
		logger.Fatal("Failed to initialize JWT", zap.Error(err))
	}
	logger.Info("JWT initialized successfully")
}

func initRouter(cfg *config.Config, db *gorm.DB, ctrlGroup *controller.ControllerGroup) {
	r := router.SetupRouter(cfg, db, ctrlGroup)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		appURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
		swaggerURL := fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port)

		logger.Info("Server started", zap.String("url", appURL))
		logger.Info("Swagger documentation", zap.String("url", swaggerURL))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	logger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
	}

	if err := database.CloseRedis(); err != nil {
		logger.Error("Failed to close Redis", zap.Error(err))
	}
	if err := database.Close(); err != nil {
		logger.Error("Failed to close database", zap.Error(err))
	}
	logger.Sync()

	logger.Info("Server stopped")
}

func initRedis(cfg *config.Config) {
	_, err := database.NewRedisClient(&database.RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	})
	if err != nil {
		logger.Warn("Failed to initialize Redis", zap.Error(err))
	} else {
		logger.Info("Redis connected successfully")
	}
}

func initDB(cfg *config.Config) *gorm.DB {
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
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	// 自动迁移数据库表
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&model.User{},
			&model.Device{},
			&model.Role{},
			&model.Permission{},
			&model.Menu{},
		); err != nil {
			logger.Fatal("Database migration failed", zap.Error(err))
		}
		logger.Info("Database migration completed")
	}

	return db
}

func initZapLog(conf *config.Config) {
	if err := logger.Init(conf.Server.Mode, conf.Log.FilePath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Application starting", zap.String("version", Version), zap.String("mode", conf.Server.Mode))
}
