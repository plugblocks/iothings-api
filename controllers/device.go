package controllers

import (
	"fmt"
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

func (dc DeviceController) GetAvailableDevices(c *gin.Context) {
	devices, err := store.GetAvailableDevices(c)

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

	if newDevice.OrganizationId == "" {
		newDevice.OrganizationId = oldDevice.OrganizationId
	}
	if newDevice.CustomerId == "" {
		newDevice.CustomerId = oldDevice.CustomerId
	}
	if newDevice.LastAccess == 0 {
		newDevice.LastAccess = oldDevice.LastAccess
	}
	newDevice.Active = oldDevice.Active
	/*changes := params.M{"$set": params.M{"organization_id": newDevice.OrganizationId, "customer_id": newDevice.CustomerId,
	"name": newDevice.Name, "ble_mac": newDevice.BleMac, "wifi_mac": newDevice.WifiMac, "sigfox_id": newDevice.SigfoxId,
	"last_access": oldDevice.LastAccess, "active": oldDevice.Active}}*/

	err = store.UpdateDevice(c, c.Param("id"), params.M{"$set": newDevice})
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
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("device_delete_failed", "Failed to delete the device", err))
		return
	}

	err = store.DeleteDeviceObservations(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("device_observations_delete_failed", "Failed to delete the device observations", err))
		return
	}

	err = store.DeleteDeviceGeolocations(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("device_geolocations_delete_failed", "Failed to delete the device geolocations", err))
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

func (dc DeviceController) GetDeviceLastLocation(c *gin.Context) {
	lastLoc, err := store.GetDeviceGeolocation(c, c.Param("id"), c.Param("source"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Get last location error")
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, lastLoc)
}

func (dc DeviceController) GetDeviceGeolocations(c *gin.Context) {
	var params models.GeolocationQueryParams

	if c.ShouldBind(&params) == nil {
		fmt.Println("params: ", params)
		if params.Limit == 0 {
			params.Limit = 100
		}
		if params.EndTime == 0 {
			params.EndTime = 2147483646 //Max uint32
		}
		if params.StartTime > params.EndTime {
			c.JSON(http.StatusInternalServerError, "Fleets geolocations query error, endTime > startTime in query")
		}
		deviceLocations, err := store.GetDeviceGeolocations(c, c.Param("id"), params.Source, params.Limit, params.StartTime, params.EndTime)

		fmt.Println("deviceLocations len: ", len(deviceLocations))

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, deviceLocations)
	} else {
		c.JSON(http.StatusInternalServerError, "GetDeviceLocations GeolocationQueryParams bind error")
	}
}

func (dc DeviceController) GetDeviceGeoJSON(c *gin.Context) {
	var params models.GeolocationQueryParams

	if c.ShouldBind(&params) == nil {
		fmt.Println("params: ", params)
		if params.Limit == 0 {
			params.Limit = 100
		}
		if params.EndTime == 0 {
			params.EndTime = 2147483646 //Max uint32
		}
		if params.StartTime > params.EndTime {
			c.JSON(http.StatusInternalServerError, "Fleets geolocations query error, endTime > startTime in query")
		}
		if params.Source == "" {
			params.Source = "wifi"
		}
		geoJsonStruct, err := store.GetDeviceGeoJSON(c, c.Param("id"), params.Source, params.Limit, params.StartTime, params.EndTime)

		fmt.Println("len: ", len(geoJsonStruct.Features))

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, geoJsonStruct)
	} else {
		c.JSON(http.StatusInternalServerError, "GetDeviceGeoJSON GeolocationQueryParams bind error")
	}
}
