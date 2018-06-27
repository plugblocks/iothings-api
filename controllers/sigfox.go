package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type SigfoxController struct {
}

func NewSigfoxController() SigfoxController {
	return SigfoxController{}
}

func (sc SigfoxController) CreateSigfoxMessage(c *gin.Context) {
	sigfoxMessage := &sigfox.Message{}

	err := c.BindJSON(sigfoxMessage)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	err = store.CreateSigfoxMessage(c, sigfoxMessage)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if sigfoxMessage.Type == "wifi" {
		res, sigfoxLocation, observation := services.ResolveWifiPosition(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while resolving WiFi computed location")
			return
		}
		fmt.Println("Resolved WiFi Frame, contaning: ", sigfoxLocation)
		sigfoxLocation.SigfoxId = sigfoxMessage.SigfoxId

		err = store.CreateSigfoxLocation(c, sigfoxLocation)
		fmt.Println("WiFi Sigfox Location created")
		if err != nil {
			fmt.Println("Error while creating WiFi Sigfox Location")
			c.Error(err)
			c.Abort()
			return
		}

		err = store.CreateObservation(c, observation)
		if err != nil {
			fmt.Println("Error while storing WiFi Sigfox Observation")
			c.Error(err)
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusCreated, sigfoxMessage)
}

func (sc SigfoxController) CreateSigfoxLocation(c *gin.Context) {
	location := &sigfox.Location{}

	err := c.BindJSON(location)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	err = store.CreateSigfoxLocation(c, location)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, location)
}
