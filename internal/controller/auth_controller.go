package controller

import (
	"gin-app-start/internal/dto"
	"gin-app-start/internal/middleware"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	result, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新 Token
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	result, err := ctrl.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, result)
}

// GetUserInfo 获取当前用户信息
func (ctrl *AuthController) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, 10003, "未授权访问")
		return
	}

	userInfo, err := ctrl.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, userInfo)
}

// ChangePassword 修改密码
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, 10003, "未授权访问")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数绑定失败: "+err.Error())
		return
	}

	if err := ctrl.authService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, nil)
}

// Logout 用户登出
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 这里可以实现 Token 黑名单等逻辑
	// 目前简单返回成功
	response.Success(c, gin.H{"message": "登出成功"})
}
