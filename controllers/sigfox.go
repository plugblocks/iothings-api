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
		/*c.Abort()
		return*/
	}

	if sigfoxMessage.Resolver == "wifi" {
		res, sigfoxLocation, observation := services.ResolveWifiPosition(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while resolving WiFi computed location")
			return
		}
		fmt.Println("at: ", observation.Timestamp, "\tValues:", observation.Values)
		sigfoxLocation.SigfoxId = sigfoxMessage.SigfoxId

		err = store.CreateSigfoxLocation(c, sigfoxLocation)
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

	} else if sigfoxMessage.Resolver == "sensitv2" {
		res, observation := services.DecodeSensitV2Message(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while parsing Sensit v2")
			return
		}
		fmt.Println("Resolved Sensit v2 Frame, containing: ", observation)

		err = store.CreateObservation(c, observation)
		if err != nil {
			fmt.Println("Error while storing Sensit v2 Sigfox Observation")
			c.Error(err)
			c.Abort()
			return
		}
	} else if sigfoxMessage.Resolver == "sensitv3" {
		res, observation := services.DecodeSensitV3Message(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while parsing Sensit v3")
			return
		}
		fmt.Println("Resolved Sensit v3 Frame, containing: ", observation)

		err = store.CreateObservation(c, observation)
		if err != nil {
			fmt.Println("Error while storing Sensit v3 Sigfox Observation")
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

	res, observation := services.SigfoxSpotit(c, location)
	if res == false {
		fmt.Println("Error while analyzing Spotit location")
		return
	}
	fmt.Println("Resolved Spotit, containing: ", observation)

	err = store.CreateObservation(c, observation)
	if err != nil {
		fmt.Println("Error while storing Sigfox Location Observation")
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, observation)
}

func (sc SigfoxController) CreateSigfoxLocationLegacy(c *gin.Context) {
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

func (sc SigfoxController) GetSigfoxLocations(c *gin.Context) {
	locations, err := store.GetSigfoxLocations(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (sc SigfoxController) GetGeoJSON(c *gin.Context) {
	geoJsonStruct, err := store.GetGeoJSON(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}
