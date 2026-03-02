package controller

import (
	"gin-app-start/internal/dto"
	"gin-app-start/internal/service"
	"gin-app-start/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	deviceService service.DeviceService
}

func NewDeviceController(deviceService service.DeviceService) *DeviceController {
	return &DeviceController{
		deviceService: deviceService,
	}
}

// CreateDevice godoc
//
//	@Summary		Create a new device
//	@Description	Create a new device with name, type and content
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateDeviceRequest	true	"Device information"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/devices [post]
func (ctrl *DeviceController) CreateDevice(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "Parameter binding failed: "+err.Error())
		return
	}

	device, err := ctrl.deviceService.CreateDevice(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, device)
}

// GetDevice godoc
//
//	@Summary		Get device by ID
//	@Description	Get device information by device ID
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Device ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/devices/{id} [get]
func (ctrl *DeviceController) GetDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 10001, "Invalid device ID")
		return
	}

	device, err := ctrl.deviceService.GetDevice(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, device)
}

// UpdateDevice godoc
//
//	@Summary		Update device information
//	@Description	Update device information by device ID
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Device ID"
//	@Param			request	body		dto.UpdateDeviceRequest	true	"Device information to update"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/devices/{id} [put]
func (ctrl *DeviceController) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 10001, "Invalid device ID")
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "Parameter binding failed: "+err.Error())
		return
	}

	device, err := ctrl.deviceService.UpdateDevice(c.Request.Context(), uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, device)
}

// DeleteDevice godoc
//
//	@Summary		Delete device
//	@Description	Delete device by device ID
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Device ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/devices/{id} [delete]
func (ctrl *DeviceController) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 10001, "Invalid device ID")
		return
	}

	if err := ctrl.deviceService.DeleteDevice(c.Request.Context(), uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	response.SuccessWithMessage(c, "Deleted successfully", nil)
}

// ListDevices godoc
//
//	@Summary		List devices
//	@Description	Get paginated list of devices
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(1)
//	@Param			page_size	query		int	false	"Page size"		default(10)
//	@Success		200			{object}	response.Response
//	@Failure		500			{object}	response.Response
//	@Router			/api/v1/devices [get]

func (ctrl *DeviceController) ListDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	devices, total, err := ctrl.deviceService.ListDevices(c.Request.Context(), page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.SuccessWithPage(c, devices, total, page, pageSize)
}
