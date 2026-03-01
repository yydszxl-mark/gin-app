package main

import (
	"context"
	"fmt"
	"gin-app-start/internal/config"
	"gin-app-start/internal/model"
	"gin-app-start/internal/router"
	"gin-app-start/pkg/database"
	"gin-app-start/pkg/logger"
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

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	initZapLog(cfg)
	initDB(cfg)
	initRedis(cfg)
	initRouter(cfg)
}

func initRouter(cfg *config.Config) {
	r := router.SetupRouter(cfg)

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

func initDB(cfg *config.Config) {
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

	//创建表
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(&model.User{}, &model.Device{}); err != nil {
			logger.Fatal("Database migration failed", zap.Error(err))
		}
		logger.Info("Database migration completed")
	}
}

func initZapLog(conf *config.Config) {
	if err := logger.Init(conf.Server.Mode, conf.Log.FilePath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Application starting", zap.String("version", Version), zap.String("mode", conf.Server.Mode))
}
