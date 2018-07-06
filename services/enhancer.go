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
	"strconv"
	"strings"
)

func ResolveWifiPosition(contxt *gin.Context, msg *sigfox.Message) (bool, *sigfox.Location, *models.Observation) {
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

	fmt.Print("WiFis: SSID1: ", ssid1, "\t SSID2:", ssid2, "\t")

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
	//var wifiLoc sigfox.Location
	wifiLoc := &sigfox.Location{}
	wifiLoc.Latitude = resp.Location.Lat
	wifiLoc.Longitude = resp.Location.Lng
	wifiLoc.Radius = resp.Accuracy
	wifiLoc.FrameNumber = msg.FrameNumber
	wifiLoc.SpotIt = false
	wifiLoc.GPS = false
	wifiLoc.WiFi = true

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

func DecodeSensitV2Message(contxt *gin.Context, msg *sigfox.Message) (bool, *models.Observation) {
	obs := &models.Observation{}
	defp := &models.DefaultProperty{"sensit", "sensor"}
	device, err := store.GetDeviceFromSigfoxId(contxt, msg.SigfoxId)
	if err != nil {
		fmt.Println("Enhancer Sigfox Device ID not found", err)
		return false, nil
	}

	//Decoder itself
	if len(msg.Data) <= 12 { //8 exactly, 4 bytes
		fmt.Println("Sensit Uplink Message")

		parsed, err := strconv.ParseUint(msg.Data, 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		data := fmt.Sprintf("%08b", parsed)
		/*byte1 := data[0:8]
		byte2 := data[8:16]
		byte3 := data[16:24]
		byte4 := data[24:32]*/

		if len(data) == 25 { //Low battery MSB
			fmt.Println("Sensit Low battery")
			//TODO: Handle low battery bit shift
			return false, nil
		}

		//Byte 1
		mode, _ := strconv.ParseInt(data[5:8], 2, 8)
		timeframe, _ := strconv.ParseInt(data[3:5], 2, 8)
		eventType, _ := strconv.ParseInt(data[1:3], 2, 8)
		batteryMsb := data[0:1]

		//Byte 2
		temperatureMsb := data[8:12]
		batteryLsb := data[12:16]
		battData := []string{batteryMsb, batteryLsb}
		battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 8)
		batVal := (float64(battery) * 0.05) + 2.7

		//Byte 3
		temperature := int64(0)
		tempVal := float64(0)

		reedSwitch := false
		if mode == 0 || mode == 1 {
			temperatureLsb := data[18:24]
			tempData := []string{temperatureMsb, temperatureLsb}
			temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
			tempVal = (float64(temperature) - 200) / 8
			if data[17] == 1 {
				reedSwitch = true
			}
		} else {
			temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
			tempVal = (float64(temperature) - 200) / 8
		}

		modeStr := ""
		swRev := ""
		humidity := float64(0.0)
		light := float64(0.0)

		switch mode {
		case 0:
			modeStr = "Button"
			majorSwRev, _ := strconv.ParseInt(data[24:28], 2, 8)
			minorSwRev, _ := strconv.ParseInt(data[28:32], 2, 8)
			swRev = fmt.Sprintf("%d.%d", majorSwRev, minorSwRev)
		case 1:
			modeStr = "Temperature + Humidity"
			humi, _ := strconv.ParseInt(data[24:32], 2, 16)
			humidity = float64(humi) * 0.5
		case 2:
			modeStr = "Light"
			lightVal, _ := strconv.ParseInt(data[18:24], 2, 8)
			lightMulti, _ := strconv.ParseInt(data[17:18], 2, 8)
			light = float64(lightVal) * 0.01
			if lightMulti == 1 {
				light = light * 8
			}
		case 3:
			modeStr = "Door"
		case 4:
			modeStr = "Move"
		case 5:
			modeStr = "Reed switch"
		default:
			modeStr = ""
		}

		timeStr := ""
		switch timeframe {
		case 0:
			timeStr = "10 mins"
		case 1:
			timeStr = "1 hour"
		case 2:
			timeStr = "6 hours"
		case 3:
			timeStr = "24 hours"
		default:
			timeStr = ""
		}

		typeStr := ""
		switch eventType {
		case 0:
			typeStr = "Regular, no alert"
		case 1:
			typeStr = "Button call"
		case 2:
			typeStr = "Alert"
		case 3:
			typeStr = "New mode"
		default:
			timeStr = ""
		}

		eventData := schema_org.QuantitativeValue{defp, "eventType", "", typeStr}
		modeData := schema_org.QuantitativeValue{defp, "mode", "", modeStr}
		timeData := schema_org.QuantitativeValue{defp, "timeframe", "seconds", timeStr}
		batData := schema_org.QuantitativeValue{defp, "battery", "volt", batVal}
		tempData := schema_org.QuantitativeValue{defp, "temperature", "celsius", tempVal} //Precision differs w/mode
		obs.Values = append(obs.Values, eventData, modeData, timeData, batData, tempData)
		obs.Timestamp = msg.Timestamp
		obs.DeviceId = device.Id

		switch mode {
		case 0:
			swRevData := schema_org.QuantitativeValue{defp, "softwareRevision", "", swRev}
			obs.Values = append(obs.Values, swRevData)
		case 1:
			//fmt.Println(humidity, "% RH")
			humiData := schema_org.QuantitativeValue{defp, "humidity", "percent", humidity}
			obs.Values = append(obs.Values, humiData)
		case 2:
			//fmt.Println(light, "lux")
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			lightData := schema_org.QuantitativeValue{defp, "light", "lux", light}
			alertData := schema_org.QuantitativeValue{defp, "alert", "", alerts}
			obs.Values = append(obs.Values, lightData, alertData)
		case 3, 4, 5:
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			alertData := schema_org.QuantitativeValue{defp, "alert", "", alerts}
			obs.Values = append(obs.Values, alertData)
		}
		if reedSwitch {
			reedData := schema_org.QuantitativeValue{defp, "reedSwitch", "", true}
			obs.Values = append(obs.Values, reedData)
		}
	} else { //len: 24 exactly, 12 bytes
		fmt.Println("Sensit Daily Downlink Message")
		//TODO: Decode sensit downlink message
		dlData := schema_org.QuantitativeValue{defp, "downlink", "", true}
		obs.Values = append(obs.Values, dlData)
	}

	return true, obs
}
