package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
)

type DeviceController struct {
}

func NewDeviceController() DeviceController {
	return DeviceController{}
}

func (dc DeviceController) CreateDevice(c *gin.Context) {
	device := &models.Device{}

	err := c.BindJSON(device)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateDevice(c, device); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, device)
}

func (dc DeviceController) GetDevices(c *gin.Context) {
	devices, err := store.GetDevices(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
}

func (dc DeviceController) UpdateDevice(c *gin.Context) {
	newDevice := models.Device{}

	err := c.BindJSON(&newDevice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	oldDevice, err := store.GetDevice(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("device_not_found", "Failed to find device id", err))
		return
	}

	changes := params.M{"$set": params.M{"organization_id": newDevice.OrganizationId, "customer_id": newDevice.CustomerId,
		"name": newDevice.Name, "ble_mac": newDevice.BleMac, "wifi_mac": newDevice.WifiMac, "sigfox_id": newDevice.SigfoxId,
		"last_access": oldDevice.LastAccess, "active": oldDevice.Active}}

	err = store.UpdateDevice(c, c.Param("id"), changes)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (dc DeviceController) DeleteDevice(c *gin.Context) {
	err := store.DeleteDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (dc DeviceController) GetDevice(c *gin.Context) {
	device, err := store.GetDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

func (dc DeviceController) GetDeviceGeoJSON(c *gin.Context) {
	geoJsonStruct, err := store.GetDeviceGeoJSON(c, c.Param("id"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}
