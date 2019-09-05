package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"gitlab.com/plugblocks/iothings-api/utils"
	"net/http"
	"time"
)

type SigfoxController struct {
}

func NewSigfoxController() SigfoxController {
	return SigfoxController{}
}

func CreateObservationAndLocation(c *gin.Context, deviceId string, obs *models.Observation, loc *models.Geolocation) {
	if obs != nil && obs.DeviceId != "" {
		err := store.CreateObservation(c, obs)
		if err != nil {
			fmt.Println("Error while storing Observation 1")
			c.Error(err)
			c.Abort()
			return
		}
	}

	if loc != nil && loc.Timestamp != 0 {
		loc.DeviceId = deviceId

		err := store.CreateGeolocation(c, loc)
		if err != nil {
			fmt.Println("Error while creating WiFi Geolocation 1 from Sigfox")
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

	loc1, loc2 := &models.Geolocation{}, &models.Geolocation{}
	obs1, obs2 := &models.Observation{}, &models.Observation{}

	if sigfoxMessage.Resolver == "wifi" {
		res, googleGeoloc, googleObs := services.ResolveWifiPosition(c, sigfoxMessage, "google")
		if res == false {
			fmt.Println("Error while enhancing WiFi Google computed location")
			return
		}
		res, hereGeoloc, hereObs := services.ResolveWifiPosition(c, sigfoxMessage, "here")
		if res == false {
			fmt.Println("Error while enhancing WiFi Google computed location")
			return
		}
		loc1, loc2 = googleGeoloc, hereGeoloc
		obs1, obs2 = googleObs, hereObs
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

	CreateObservationAndLocation(c, device.Id, obs1, loc1)
	CreateObservationAndLocation(c, device.Id, obs2, loc2)

	c.JSON(http.StatusCreated, sigfoxMessage)
}

func CreateObsByAdvancedMessage(sourceName string, deviceId string, sigfoxMessage *sigfox.MessageDataAdvanced) (sfxLoc *models.Geolocation, sfxObs *models.Observation) {
	sfxLoc = &models.Geolocation{DeviceId: deviceId, Timestamp: sigfoxMessage.Timestamp, Latitude: sigfoxMessage.ComputedLocation.Lat,
		Longitude: sigfoxMessage.ComputedLocation.Lng, Radius: float64(sigfoxMessage.ComputedLocation.Radius), Source: sourceName}
	defp := &models.SemanticProperty{Context: "sigfox-antennas", Type: "location"}
	latVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "latitude", UnitText: "degrees", Value: sigfoxMessage.ComputedLocation.Lat}
	lngVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "longitude", UnitText: "degrees", Value: sigfoxMessage.ComputedLocation.Lng}
	accVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "accuracy", UnitText: "meters", Value: sigfoxMessage.ComputedLocation.Radius}
	sfxObs = &models.Observation{Timestamp: sigfoxMessage.Timestamp, DeviceId: deviceId, Resolver: sourceName,
		Values: []models.QuantitativeValue{latVal, lngVal, accVal},
	}
	return sfxLoc, sfxObs
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

	sfxLoc, sfxObs := &models.Geolocation{}, &models.Observation{}

	source, status := sigfoxMessage.ComputedLocation.Source, sigfoxMessage.ComputedLocation.Status

	if source == 2 { // Regular Atlas Basic
		sfxLoc, sfxObs = CreateObsByAdvancedMessage("sigfox-antennas", device.Id, sigfoxMessage)
		CreateObservationAndLocation(c, device.Id, sfxObs, sfxLoc)
	} else if source == 6 { //Here WiFi by Sigfox
		switch status {
		case 0: // Wrong antennas, impossible to have a position
			c.JSON(http.StatusAccepted, "No geolocation returned by Sigfox")
		case 1: // Correctly resolved by Sigfox/Here
			sfxLoc, sfxObs = CreateObsByAdvancedMessage("wifi-sigfox", device.Id, sigfoxMessage)
			CreateObservationAndLocation(c, device.Id, sfxObs, sfxLoc)
		case 2: // Not resolved by Sigfox/Here
			res, googleGeoloc, googleObs := services.ResolveWifiPosition(c, sigfoxMessage, "google")
			if res == false {
				fmt.Println("Error while enhancing WiFi computed location")
				return
			}
			CreateObservationAndLocation(c, device.Id, googleObs, googleGeoloc)
			return
		case 20:
			logrus.Warnln("Invalid Wifi Payload to be resolved")
			return
		}
	}

	loc1, loc2 := &models.Geolocation{}, &models.Geolocation{}
	obs1, obs2 := &models.Observation{}, &models.Observation{}
	if sigfoxMessage.Resolver == "wifi" && source == 2 { // Contract without wifi
		res1, googleGeoloc, googleObs := services.ResolveWifiPosition(c, sigfoxMessage, "google")
		if res1 == false {
			fmt.Println("Error while enhancing WiFi Google computed location")
			return
		}

		res2, hereGeoloc, hereObs := services.ResolveWifiPosition(c, sigfoxMessage, "here")
		if res2 == false {
			fmt.Println("Error while enhancing WiFi Google computed location")
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

	CreateObservationAndLocation(c, device.Id, obs1, loc1)
	CreateObservationAndLocation(c, device.Id, obs2, loc2)

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
