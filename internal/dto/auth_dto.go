package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string   `json:"token"`
	ExpiresAt int64    `json:"expires_at"`
	UserInfo  UserInfo `json:"user_info"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Avatar   string   `json:"avatar"`
	Status   int8     `json:"status"`
	Roles    []string `json:"roles"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}
