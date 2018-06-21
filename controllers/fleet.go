package controllers

import (
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

func (gtc FleetController) GetFleets(c *gin.Context) {
	fleets, err := store.GetAllFleets(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleets)
}

func (gtc FleetController) GetFleetById(c *gin.Context) {
	id := c.Param("id")

	fleet, err := store.GetFleetById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleet)
}

func (gtc FleetController) CreateFleet(c *gin.Context) {
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

func (gtc FleetController) EditFleet(c *gin.Context) {
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
