package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"gitlab.com/plugblocks/iothings-api/utils"
	"net/http"
	"time"
	"unsafe"
)

type SigfoxController struct {
}

func NewSigfoxController() SigfoxController {
	return SigfoxController{}
}

func CreateObservationsAndLocations(c *gin.Context, deviceId string, obs1 *models.Observation, obs2 *models.Observation, loc1 *models.Geolocation, loc2 *models.Geolocation) {
	fmt.Println("Sizes:", unsafe.Sizeof(obs1), unsafe.Sizeof(obs2), unsafe.Sizeof(loc1), unsafe.Sizeof(loc2))
	if obs1 != nil && obs1.DeviceId != "" {
		err := store.CreateObservation(c, obs1)
		if err != nil {
			fmt.Println("Error while storing Observation 1")
			c.Error(err)
			c.Abort()
			return
		}
	}

	if obs2 != nil && obs2.DeviceId != "" {
		err := store.CreateObservation(c, obs2)
		if err != nil {
			fmt.Println("Error while storing Observation 2")
			c.Error(err)
			c.Abort()
			return
		}
	}

	if loc1 != nil && loc1.Timestamp != 0 {
		loc1.DeviceId = deviceId

		err := store.CreateGeolocation(c, loc1)
		if err != nil {
			fmt.Println("Error while creating WiFi Geolocation 1 from Sigfox")
			c.Error(err)
			c.Abort()
			return
		}
	}

	if loc2 != nil && loc2.Timestamp != 0 {
		loc2.DeviceId = deviceId

		err := store.CreateGeolocation(c, loc2)
		if err != nil {
			fmt.Println("Error while creating WiFi Geolocation 2 from Sigfox")
			c.Error(err)
			c.Abort()
			return
		}
	}
}

func (SigfoxController) CreateSigfoxMessage(c *gin.Context) {
	sigfoxMessage := &sigfox.Message{}

	err := c.BindJSON(sigfoxMessage)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	device, err := store.GetDeviceFromSigfoxId(c, sigfoxMessage.SigfoxId)
	if err != nil {
		fmt.Println("Sigfox Device ID not found", err)
		device := models.Device{bson.NewObjectId().Hex(), "", "", "Sigfox Device: " + sigfoxMessage.SigfoxId,
			"", "", sigfoxMessage.SigfoxId, time.Now().Unix(), true}
		utils.CheckErr(store.CreateDevice(c, &device))
		return
	}

	err = store.CreateSigfoxMessage(c, sigfoxMessage)
	if err != nil {
		c.Error(err)
	}

	var loc1 *models.Geolocation
	var loc2 *models.Geolocation
	var obs1 *models.Observation
	var obs2 *models.Observation
	if sigfoxMessage.Resolver == "wifi" {
		res, googleGeoloc, hereGeoloc, googleObs, hereObs := services.ResolveWifiPosition(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while enhancing WiFi computed location")
			return
		}
		loc1, loc2 = googleGeoloc, hereGeoloc
		obs1, obs2 = googleObs, hereObs
		//fmt.Println("at: ", obs1.Timestamp, "\tValues:", obs1.Values)
	} else if sigfoxMessage.Resolver == "sensitv2" {
		res, observation := services.DecodeSensitV2Message(c, sigfoxMessage.SigfoxId, sigfoxMessage.Data, sigfoxMessage.Timestamp)
		if res == false {
			fmt.Println("Error while enhancing Sensit v2")
			return
		}
		obs1 = observation
		fmt.Println("Resolved Sensit v2 Frame, containing: ", observation)
	} else if sigfoxMessage.Resolver == "sensitv3" {
		res, observation := services.DecodeSensitV3Message(c, sigfoxMessage.SigfoxId, sigfoxMessage.Data, sigfoxMessage.Timestamp)
		if res == false {
			fmt.Println("Error while enhancing Sensit v3")
			return
		}
		obs1 = observation
		fmt.Println("Resolved Sensit v3 Frame, containing: ", observation)
	} else if sigfoxMessage.Resolver == "wisol" {
		res, lo1, lo2, ob1, ob2 := services.Wisol(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while enhancing Wisol frame")
			return
		}
		obs1, obs2, loc1, loc2 = ob1, ob2, lo1, lo2
		fmt.Println("Resolved wisol Frame, containing: ", obs1, obs2)
	}

	CreateObservationsAndLocations(c, device.Id, obs1, obs2, loc1, loc2)

	c.JSON(http.StatusCreated, sigfoxMessage)
}

func (SigfoxController) CreateSigfoxDataAdvancedMessage(c *gin.Context) {
	sigfoxMessage := &sigfox.MessageDataAdvanced{}

	err := c.BindJSON(sigfoxMessage)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	err = store.CreateSigfoxDataAdvancedMessage(c, sigfoxMessage)
	if err != nil {
		c.Error(err)
	}

	device, err := store.GetDeviceFromSigfoxId(c, sigfoxMessage.SigfoxId)
	if err != nil {
		fmt.Println("Sigfox Device ID", sigfoxMessage.SigfoxId, "not found", err)
		/*device := models.Device{bson.NewObjectId().Hex(), "", "", "Sigfox Device: " + sigfoxMessage.SigfoxId,
			"", "", sigfoxMessage.SigfoxId, time.Now().Unix(), true}
		utils.CheckErr(store.CreateDevice(c, &device))*/
		return
	}

	var loc1 *models.Geolocation
	var loc2 *models.Geolocation
	var obs1 *models.Observation
	var obs2 *models.Observation
	if sigfoxMessage.Resolver == "wifi" {
		res, googleGeoloc, hereGeoloc, googleObs, hereObs := services.ResolveWifiPosition(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while enhancing WiFi computed location")
			return
		}
		loc1, loc2 = googleGeoloc, hereGeoloc
		obs1, obs2 = googleObs, hereObs
		//fmt.Println("at: ", obs1.Timestamp, "\tValues:", obs1.Values)
	} else if sigfoxMessage.Resolver == "sensitv2" {
		res, observation := services.DecodeSensitV2Message(c, sigfoxMessage.SigfoxId, sigfoxMessage.Data, sigfoxMessage.Timestamp)
		if res == false {
			fmt.Println("Error while enhancing Sensit v2")
			return
		}
		obs1 = observation
		fmt.Println("Resolved Sensit v2 Frame, containing: ", observation)
	} else if sigfoxMessage.Resolver == "sensitv3" {
		res, observation := services.DecodeSensitV3Message(c, sigfoxMessage.SigfoxId, sigfoxMessage.Data, sigfoxMessage.Timestamp)
		if res == false {
			fmt.Println("Error while enhancing Sensit v3")
			return
		}
		obs1 = observation
		fmt.Println("Resolved Sensit v3 Frame, containing: ", observation)
	} else if sigfoxMessage.Resolver == "wisol" {
		res, lo1, lo2, ob1, ob2 := services.Wisol(c, sigfoxMessage)
		if res == false {
			fmt.Println("Error while enhancing Wisol frame")
			return
		}
		obs1, obs2, loc1, loc2 = ob1, ob2, lo1, lo2
		fmt.Println("Resolved wisol Frame, containing: ", obs1, obs2)
	}

	CreateObservationsAndLocations(c, device.Id, obs1, obs2, loc1, loc2)

	c.JSON(http.StatusCreated, sigfoxMessage)

}

func (SigfoxController) CreateSigfoxLocation(c *gin.Context) {
	sigfoxLocation := &sigfox.Location{}

	err := c.BindJSON(sigfoxLocation)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	device, err := store.GetDeviceFromSigfoxId(c, sigfoxLocation.SigfoxId)
	if err != nil {
		fmt.Println("Wifi Enhancer Sigfox Device ID not found", err)
		return
	}

	res, geoLocation, observation := services.SigfoxSpotit(c, sigfoxLocation)
	if res == false {
		fmt.Println("Error while analyzing Spotit location")
		return
	}
	fmt.Println("Resolved Spotit, containing: ", observation)

	geoLocation.DeviceId = device.Id

	err = store.CreateGeolocation(c, geoLocation)
	if err != nil {
		fmt.Println("Error while creating WiFi Geolocation from Sigfox")
		c.Error(err)
		c.Abort()
		return
	}

	err = store.CreateObservation(c, observation)
	if err != nil {
		fmt.Println("Error while storing Sigfox Location Observation")
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, observation)
}

func (SigfoxController) CreateSigfoxLocationLegacy(c *gin.Context) {
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

func (SigfoxController) GetSigfoxLocations(c *gin.Context) {
	locations, err := store.GetSigfoxLocations(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (SigfoxController) GetGeoJSON(c *gin.Context) {
	geoJsonStruct, err := store.GetGeoJSON(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, geoJsonStruct)
}
