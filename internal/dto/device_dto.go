package dto

type CreateDeviceRequest struct {
	Name    string `json:"name" binding:"required" example:"iPhone 13"`
	Type    string `json:"type" binding:"required" example:"phone"`
	Content string `json:"content"  example:"{}"`
}

type UpdateDeviceRequest struct {
	Name    string `json:"name" binding:"required" example:"iPhone 13"`
	Content string `json:"content"  example:"{}"`
}
