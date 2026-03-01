package model

// Role 角色模型
type Role struct {
	BaseModel
	Name        string       `gorm:"size:64;uniqueIndex;not null" json:"name" binding:"required"` // 角色名称
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code" binding:"required"` // 角色编码
	Description string       `gorm:"size:255" json:"description"`                                 // 角色描述
	Sort        int          `gorm:"default:0" json:"sort"`                                       // 排序
	Status      int          `gorm:"default:1" json:"status"`                                     // 状态：1-启用，0-禁用
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`    // 角色拥有的权限
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`                // 角色拥有的菜单
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`                // 拥有该角色的用户
}

func (Role) TableName() string {
	return "roles"
}
