package dto

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name" binding:"required,min=2,max=64"`
	Title     string `json:"title" binding:"required,min=2,max=64"`
	Icon      string `json:"icon" binding:"max=128"`
	Path      string `json:"path" binding:"max=255"`
	Component string `json:"component" binding:"max=255"`
	Redirect  string `json:"redirect" binding:"max=255"`
	Type      string `json:"type" binding:"required,oneof=menu button"`
	Hidden    bool   `json:"hidden"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status" binding:"oneof=0 1"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name" binding:"omitempty,min=2,max=64"`
	Title     string `json:"title" binding:"omitempty,min=2,max=64"`
	Icon      string `json:"icon" binding:"max=128"`
	Path      string `json:"path" binding:"max=255"`
	Component string `json:"component" binding:"max=255"`
	Redirect  string `json:"redirect" binding:"max=255"`
	Type      string `json:"type" binding:"omitempty,oneof=menu button"`
	Hidden    bool   `json:"hidden"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// MenuResponse 菜单响应
type MenuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Name      string         `json:"name"`
	Title     string         `json:"title"`
	Icon      string         `json:"icon"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Redirect  string         `json:"redirect"`
	Type      string         `json:"type"`
	Hidden    bool           `json:"hidden"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
	Children  []MenuResponse `json:"children,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// MenuTreeResponse 菜单树响应
type MenuTreeResponse struct {
	ID       uint               `json:"id"`
	ParentID uint               `json:"parent_id"`
	Title    string             `json:"title"`
	Children []MenuTreeResponse `json:"children,omitempty"`
}
