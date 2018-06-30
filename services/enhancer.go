package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/schema_org"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/store"
	"googlemaps.github.io/maps"
	"log"
)

func ResolveWifiPosition(contxt *gin.Context, msg *sigfox.Message) (bool, *sigfox.Location, *models.Observation) {
	fmt.Print("WiFi frame: \t\t\t")

	if len(msg.Data) <= 12 {
		fmt.Println("Only one WiFi, frame don't resolve for privacy issues")
		return false, nil, nil
	}

	ssid1 := ""
	for i := 0; i <= 10; i += 2 {
		if i == 10 {
			ssid1 += fmt.Sprint(string(msg.Data[i : i+2]))
		} else {
			ssid1 += fmt.Sprint(string(msg.Data[i:i+2]), ":")
		}
	}
	ssid2 := ""
	for i := 12; i <= 22; i += 2 {
		if i == 22 {
			ssid2 += fmt.Sprint(string(msg.Data[i : i+2]))
		} else {
			ssid2 += fmt.Sprint(string(msg.Data[i:i+2]), ":")
		}
	}

	fmt.Print("SSID1: ", ssid1, "\t SSID2:", ssid2, "\t\t\t")

	googleApiKey := config.GetString(contxt, "google_api_key")

	c, err := maps.NewClient(maps.WithAPIKey(googleApiKey))
	if err != nil {
		log.Fatalf("API connection fatal error: %s", err)
	}
	r := &maps.GeolocationRequest{
		ConsiderIP: false,
		WiFiAccessPoints: []maps.WiFiAccessPoint{{
			MACAddress: ssid1,
		}, {
			MACAddress: ssid2,
		}},
	}
	resp, err := c.Geolocate(context.Background(), r)
	if err != nil {
		fmt.Println("Google Maps Geolocation Request, Position:", err)
		return false, nil, nil
	}

	//Else, position is resolved
	fmt.Println("Google Maps Geolocation resolved")
	//var wifiLoc sigfox.Location
	wifiLoc := &sigfox.Location{}
	wifiLoc.Latitude = resp.Location.Lat
	wifiLoc.Longitude = resp.Location.Lng
	wifiLoc.Radius = resp.Accuracy
	wifiLoc.FrameNumber = msg.FrameNumber
	wifiLoc.SpotIt = false
	wifiLoc.GPS = false
	wifiLoc.WiFi = true

	fmt.Println(resp)
	fmt.Println(wifiLoc)

	obs := &models.Observation{}
	defp := &models.DefaultProperty{"wifi", "location"}
	device, err := store.GetDeviceFromSigfoxId(contxt, msg.SigfoxId)
	if err != nil {
		fmt.Println("Enhancer Sigfox Device ID not found", err)
		return false, nil, nil
	}
	latVal := schema_org.QuantitativeValue{defp, "latitude", "degrees", resp.Location.Lat}
	lngVal := schema_org.QuantitativeValue{defp, "longitude", "degrees", resp.Location.Lng}
	accVal := schema_org.QuantitativeValue{defp, "accuracy", "meters", resp.Accuracy}
	obs.Values = append(obs.Values, latVal, lngVal, accVal)
	obs.Timestamp = msg.Timestamp
	obs.DeviceId = device.Id
	return true, wifiLoc, obs
}
