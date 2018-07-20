package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
)

type GeolocationController struct {
}

func NewGeolocationController() GeolocationController {
	return GeolocationController{}
}

func (gc GeolocationController) CreateGeolocation(c *gin.Context) {
	geolocation := &models.Geolocation{}

	err := c.BindJSON(geolocation)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateGeolocation(c, geolocation); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, geolocation)
}

/*func (gc GeolocationController) GetGeolocation(c *gin.Context) {
	device, err := store.GetGeolocation(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}*/

func (gc GeolocationController) DeleteGeolocation(c *gin.Context) {
	err := store.DeleteDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

/*func (gc GeolocationController) GetDeviceGeoJSON(c *gin.Context) {
	geoJsonStruct, err := store.GetDeviceGeoJSON(c, c.Param("id"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}


func (gc GeolocationController) GetFleetGeoJSON(c *gin.Context) {
	id := c.Param("id")

	geoJsonStruct, err := store.GetFleetGeoJSON(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}*/
