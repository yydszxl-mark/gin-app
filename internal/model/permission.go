package model

// Permission 权限模型
type Permission struct {
	BaseModel
	Name        string `gorm:"size:64;not null" json:"name" binding:"required"`              // 权限名称
	Code        string `gorm:"size:128;uniqueIndex;not null" json:"code" binding:"required"` // 权限编码（如：user:create）
	Type        string `gorm:"size:32;not null" json:"type"`                                 // 权限类型：menu-菜单，button-按钮，api-接口
	Method      string `gorm:"size:16" json:"method"`                                        // HTTP方法：GET, POST, PUT, DELETE
	Path        string `gorm:"size:255" json:"path"`                                         // API路径
	Description string `gorm:"size:255" json:"description"`                                  // 权限描述
	Sort        int    `gorm:"default:0" json:"sort"`                                        // 排序
	Status      int    `gorm:"default:1" json:"status"`                                      // 状态：1-启用，0-禁用
	Roles       []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`           // 拥有该权限的角色
}

func (Permission) TableName() string {
	return "permissions"
}
