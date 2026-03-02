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

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/api/v1/auth/login [post]
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

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Get new access token using refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/api/v1/auth/refresh [post]
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

// GetUserInfo godoc
//
//	@Summary		Get current user info
//	@Description	Get authenticated user information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.UserInfo}
//	@Failure		401	{object}	response.Response
//	@Router			/api/v1/user/info [get]
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

// ChangePassword godoc
//
//	@Summary		Change password
//	@Description	Change current user password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ChangePasswordRequest	true	"Password change request"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/api/v1/user/change-password [post]
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

// Logout godoc
//
//	@Summary		User logout
//	@Description	Logout current user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response
//	@Router			/api/v1/user/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 这里可以实现 Token 黑名单等逻辑
	// 目前简单返回成功
	response.Success(c, gin.H{"message": "登出成功"})
}
