package controllers

import (
	"fmt"
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

func (ObservationController) CreateObservation(c *gin.Context) {
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

func (ObservationController) GetDeviceObservations(c *gin.Context) {
	var params models.ObservationQueryParams
	if c.ShouldBind(&params) == nil {
		if params.Limit == 0 {
			params.Limit = 10
		}
		order := "-timestamp"
		if params.Order != "" {
			order = params.Order
		}

		observations, err := store.GetDeviceObservations(c, c.Param("deviceId"), params.Resolver, order, params.Limit)

		fmt.Println("len: ", len(observations))

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, observations)
	} else {
		c.JSON(http.StatusInternalServerError, "Device observation error")
	}
}
func (ObservationController) GetFleetObservations(c *gin.Context) {
	var params models.ObservationQueryParams
	if c.ShouldBind(&params) == nil {
		if params.Limit == 0 {
			params.Limit = 10
		}
		order := "-timestamp"
		if params.Order == "" {
			order = params.Order
		}
		observations, err := store.GetFleetObservations(c, c.Param("fleetId"), params.Resolver, order, params.Limit)

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, observations)
	} else {
		c.JSON(http.StatusInternalServerError, "Fleet observation error")
	}
}

func (ObservationController) DeleteObservation(c *gin.Context) {
	err := store.DeleteObservation(c, c.Param("id"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
