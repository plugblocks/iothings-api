package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type FleetController struct{}

func NewFleetController() FleetController {
	return FleetController{}
}

func (fc FleetController) AddDeviceToFleet(c *gin.Context) {
	fleet, err := store.AddDeviceToFleet(c, c.Param("id"), c.Param("deviceId"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleet)
}

func (fc FleetController) GetFleets(c *gin.Context) {
	fleets, err := store.GetAllFleets(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleets)
}

func (fc FleetController) GetFleetById(c *gin.Context) {
	id := c.Param("id")

	fleet, err := store.GetFleetById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleet)
}

func (fc FleetController) GetDevicesFromFleet(c *gin.Context) {
	id := c.Param("id")

	devices, err := store.GetDevicesFromFleet(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
}

func (fc FleetController) GetFleetGeoJSON(c *gin.Context) {
	id := c.Param("id")

	geoJsonStruct, err := store.GetFleetGeoJSON(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}

func (fc FleetController) GetFleetsGeoJSON(c *gin.Context) {

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
		geoJsonStruct, err := store.GetFleetsGeoJSON(c, params.Source, params.Limit, params.StartTime, params.EndTime)

		fmt.Println("len: ", len(geoJsonStruct.Features))

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, geoJsonStruct)
	} else {
		c.JSON(http.StatusInternalServerError, "Fleets geolocations error")
	}

	/*geoJsonStruct, err := store.GetFleetsGeoJSON(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)*/
}

func (fc FleetController) GetUserFleetsGeoJSON(c *gin.Context) {
	geoJsonStruct, err := store.GetUserFleetsGeoJSON(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}

func (fc FleetController) CreateFleet(c *gin.Context) {
	fleet := &models.Fleet{}

	if err := c.BindJSON(fleet); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateFleet(c, fleet); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, fleet)
}

func (fc FleetController) UpdateFleet(c *gin.Context) {
	fleet := &models.Fleet{}
	id := c.Param("id")

	if err := c.BindJSON(fleet); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateFleet(c, id, params.M{"$set": fleet}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, fleet)
}

func (fc FleetController) DeleteFleet(c *gin.Context) {
	err := store.DeleteFleet(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)

}
