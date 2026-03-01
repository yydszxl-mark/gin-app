package controller

import (
	"gin-app-start/pkg/errors"
	"gin-app-start/pkg/logger"
	"gin-app-start/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// handleServiceError 处理服务层错误（全局函数）
func handleServiceError(c *gin.Context, err error) {
	var bizErr *errors.BusinessError
	if e, ok := err.(*errors.BusinessError); ok {
		bizErr = e
		response.Error(c, bizErr.Code, bizErr.Message)
	} else {
		logger.Error("Unknown error", zap.Error(err))
		response.Error(c, 50000, "Internal server error")
	}
}
