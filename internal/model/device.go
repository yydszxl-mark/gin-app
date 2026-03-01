package model

type Device struct {
	BaseModel
	Name    string `gorm:"size:64;uniqueIndex;not null" json:"name" binding:"required" example:"iPhone 13"`
	Type    string `gorm:"size:64;uniqueIndex;not null" json:"type" binding:"required" example:"phone"`
	Content string `gorm:"size:1024;not null" json:"content"  example:"iPhone 13"`
}

func (Device) TableName() string {
	return "devices"
}
