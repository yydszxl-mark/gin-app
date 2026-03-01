package model

// Menu 菜单模型
type Menu struct {
	BaseModel
	ParentID  uint   `gorm:"default:0;index" json:"parent_id"`                 // 父菜单ID，0表示顶级菜单
	Name      string `gorm:"size:64;not null" json:"name" binding:"required"`  // 菜单名称
	Title     string `gorm:"size:64;not null" json:"title" binding:"required"` // 菜单标题（显示名称）
	Icon      string `gorm:"size:128" json:"icon"`                             // 菜单图标
	Path      string `gorm:"size:255" json:"path"`                             // 路由路径
	Component string `gorm:"size:255" json:"component"`                        // 组件路径
	Redirect  string `gorm:"size:255" json:"redirect"`                         // 重定向路径
	Type      string `gorm:"size:32;not null;default:'menu'" json:"type"`      // 类型：menu-菜单，button-按钮
	Hidden    bool   `gorm:"default:false" json:"hidden"`                      // 是否隐藏
	Sort      int    `gorm:"default:0" json:"sort"`                            // 排序
	Status    int    `gorm:"default:1" json:"status"`                          // 状态：1-启用，0-禁用
	Children  []Menu `gorm:"-" json:"children,omitempty"`                      // 子菜单（不存储在数据库）
	Roles     []Role `gorm:"many2many:role_menus;" json:"roles,omitempty"`     // 拥有该菜单的角色
}

func (Menu) TableName() string {
	return "menus"
}
