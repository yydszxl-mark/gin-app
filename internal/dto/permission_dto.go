package dto

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Code        string `json:"code" binding:"required,min=2,max=128"`
	Type        string `json:"type" binding:"required,oneof=menu button api"`
	Method      string `json:"method" binding:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Path        string `json:"path" binding:"max=255"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=64"`
	Type        string `json:"type" binding:"omitempty,oneof=menu button api"`
	Method      string `json:"method" binding:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Path        string `json:"path" binding:"max=255"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
