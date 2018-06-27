package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type ObservationController struct {
}

func NewObservationController() ObservationController {
	return ObservationController{}
}

func (oc ObservationController) CreateObservation(c *gin.Context) {
	observation := &models.Observation{}

	err := c.BindJSON(observation)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateObservation(c, observation); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, observation)
}

func (oc ObservationController) GetDeviceObservations(c *gin.Context) {
	observations, err := store.GetDeviceObservations(c, c.Param("customerId"), c.Param("deviceId"), c.Param("type"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, observations)
}

func (oc ObservationController) GetDeviceLatestObservation(c *gin.Context) {
	observation, err := store.GetDeviceLatestObservation(c, c.Param("customerId"), c.Param("deviceId"), c.Param("type"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, observation)
}
func (oc ObservationController) GetFleetObservations(c *gin.Context) {
	observations, err := store.GetFleetObservations(c, c.Param("id"), c.Param("type"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, observations)
}
func (oc ObservationController) GetFleetLatestObservation(c *gin.Context) {
	observation, err := store.GetFleetLatestObservation(c, c.Param("id"), c.Param("type"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, observation)
}