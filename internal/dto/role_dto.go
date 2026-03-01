package dto

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Code        string `json:"code" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=64"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// AssignMenusRequest 分配菜单请求
type AssignMenusRequest struct {
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
