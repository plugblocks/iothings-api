package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/store"
	"googlemaps.github.io/maps"
)

// These are all a mix of enhancers / resolvers. TODO: Split this in pieces.
func ResolveWifiPosition(contxt *gin.Context, msg *sigfox.Message) (bool, *models.Geolocation, *models.Observation) {
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

	//fmt.Print("WiFis: SSID1: ", ssid1, "\t SSID2:", ssid2, "\t")

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
		fmt.Println("Google WiFi Geolocation: ", err, "ssid1: ", ssid1, "ssid2: ", ssid2)
		return false, nil, nil
	}

	device, err := store.GetDeviceFromSigfoxId(contxt, msg.SigfoxId)
	if err != nil {
		fmt.Println("Wifi Enhancer Sigfox Device ID not found", err)
		return false, nil, nil
	}

	//Else, position is resolved
	//var wifiLoc sigfox.Location
	wifiLoc := &models.Geolocation{}
	wifiLoc.DeviceId = device.Id
	wifiLoc.Timestamp = msg.Timestamp
	wifiLoc.Source = "wifi"
	wifiLoc.Latitude = resp.Location.Lat
	wifiLoc.Longitude = resp.Location.Lng
	wifiLoc.Radius = resp.Accuracy

	obs := &models.Observation{}
	defp := &models.SemanticProperty{Context: "wifi", Type: "location"}
	latVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "latitude", UnitText: "degrees", Value: resp.Location.Lat}
	lngVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "longitude", UnitText: "degrees", Value: resp.Location.Lng}
	accVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "accuracy", UnitText: "meters", Value: resp.Accuracy}
	obs.Values = append(obs.Values, latVal, lngVal, accVal)
	obs.Timestamp = msg.Timestamp
	obs.DeviceId = device.Id
	obs.Resolver = "wifi"

	return true, wifiLoc, obs
}

func DecodeSensitV2Message(contxt *gin.Context, msg *sigfox.Message) (bool, *models.Observation) {
	obs := &models.Observation{}
	defp := &models.SemanticProperty{Context: "sensit", Type: "sensor"}
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
		battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 16)
		batVal := (float64(battery) * 0.05) + 2.7
		batVal = math.Round(batVal*100) / 100
		//Byte 3
		temperature := int64(0)
		tempVal := float32(0)

		reedSwitch := false
		if mode == 0 || mode == 1 {
			temperatureLsb := data[18:24]
			tempData := []string{temperatureMsb, temperatureLsb}
			temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
			tempVal = (float32(temperature) - 200) / 8
			if data[17] == 1 {
				reedSwitch = true
			}
		} else {
			temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
			tempVal = (float32(temperature) - 200) / 8
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
			modeStr = "Vibration"
		case 5:
			modeStr = "Magnet"
		default:
			modeStr = ""
		}

		timeVal := 0
		timeUnit := ""
		switch timeframe {
		case 0:
			timeVal = 10
			timeUnit = "minutes"
		case 1:
			timeVal = 1
			timeUnit = "hour"
		case 2:
			timeVal = 6
			timeUnit = "hours"
		case 3:
			timeVal = 24
			timeUnit = "hours"
		default:
			timeVal = 10
			timeUnit = "minutes"
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
			typeStr = ""
		}

		eventData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "eventType", UnitText: "", Value: typeStr}
		modeData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "mode", UnitText: "", Value: modeStr}
		timeData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "timeframe", UnitText: timeUnit, Value: timeVal}
		batData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "battery", UnitText: "volt", Value: batVal}
		tempData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "temperature", UnitText: "celsius", Value: tempVal} //Precision differs w/mode
		obs.Values = append(obs.Values, eventData, modeData, timeData, batData, tempData)

		switch mode {
		case 0:
			swRevData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "softwareRevision", UnitText: "", Value: swRev}
			obs.Values = append(obs.Values, swRevData)
		case 1:
			//fmt.Println(humidity, "% RH")
			humiData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "humidity", UnitText: "percent", Value: humidity}
			obs.Values = append(obs.Values, humiData)
		case 2:
			//fmt.Println(light, "lux")
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			lightData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "light", UnitText: "lux", Value: light}
			alertData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "alert", UnitText: "", Value: alerts}
			obs.Values = append(obs.Values, lightData, alertData)
		case 3, 4, 5:
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			alertData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "alert", UnitText: "", Value: alerts}
			obs.Values = append(obs.Values, alertData)
		}
		if reedSwitch {
			reedData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "reedSwitch", UnitText: "", Value: true}
			obs.Values = append(obs.Values, reedData)
		}
	} else { //len: 24 exactly, 12 bytes
		fmt.Println("Sensit Daily Downlink Message")
		//TODO: Decode sensit downlink message
		dlData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "downlink", UnitText: "", Value: true}
		obs.Values = append(obs.Values, dlData)
	}
	obs.Timestamp = msg.Timestamp
	obs.DeviceId = device.Id
	obs.Resolver = "sensitv2"

	return true, obs
}

func DecodeSensitV3Message(contxt *gin.Context, msg *sigfox.Message) (bool, *models.Observation) {
	obs := &models.Observation{}
	defp := &models.SemanticProperty{Context: "sensit", Type: "sensor"}
	device, err := store.GetDeviceFromSigfoxId(contxt, msg.SigfoxId)
	if err != nil {
		fmt.Println("Enhancer Sigfox Device ID not found", err)
		return false, nil
	}

	modeStr := ""
	fwRevVal := ""
	humiVal := float32(0.0)
	tempVal := float32(0.0)
	lightVal := float32(0.0)

	fmt.Println("len(msg.Data):", len(msg.Data))
	//Decoder itself
	if len(msg.Data) <= 12 { //8 exactly, 4 bytes
		fmt.Println("Sensit Uplink Message")

		parsed, err := strconv.ParseUint(msg.Data, 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		data := fmt.Sprintf("%08b", parsed)

		fmt.Println("len(data):", len(data))
		if len(data) == 25 { //Low battery MSB
			fmt.Println("Sensit Low battery")
			//TODO: Handle low battery bit shift
			return false, nil
		}

		//Byte 1 : 5b Battery & 3b reserved (0b110)
		battery, _ := strconv.ParseInt(data[0:5], 2, 8)
		batVal := (float64(battery) * 0.05) + 2.7
		batVal = math.Round(batVal*100) / 100
		// reserved, _ := strconv.ParseInt(data[5:8], 2, 8) //Should be 0b110

		//Byte 2 : 5b Mode, 1b Alert Button, 2b data
		mode, _ := strconv.ParseInt(data[8:13], 2, 8)
		buttonStr := ""
		if data[13:14] == "0" {
			buttonStr = "Not pressed"
		} else {
			buttonStr = "Pressed"
		}

		evtVal := ""
		switch mode {
		case 0:
			modeStr = "Standby"
			fwRevMaj, _ := strconv.ParseInt(data[16:20], 2, 8)
			fwRevMinJoin := []string{data[20:24], data[24:26]}
			fwRevMin, _ := strconv.ParseInt(strings.Join(fwRevMinJoin, ""), 2, 16)
			fwRevPatch, _ := strconv.ParseInt(data[26:32], 2, 8)
			fwRevVal = fmt.Sprintf("%d.%d.%d", fwRevMaj, fwRevMin, fwRevPatch)
		case 1:
			modeStr = "Temperature + Humidity"
			tempTab := []string{data[14:16], data[16:24]}
			tempJoin := strings.Join(tempTab, "")
			fmt.Println("tempJoin:", tempJoin)
			temp, _ := strconv.ParseInt(tempJoin, 2, 16)
			tempVal = (float32(temp) - 200) / 8
			fmt.Println("MSB:", data[14:16], "\tLSB:", data[16:24])
			fmt.Println("temp:", temp, "\t tempVal:", tempVal)
			humi, _ := strconv.ParseInt(data[24:32], 2, 16)
			humiVal = float32(humi) * 0.5
		case 2:
			modeStr = "Light"
			lightJoin := []string{data[16:24], data[24:32]}
			light, _ := strconv.ParseInt(strings.Join(lightJoin, ""), 2, 16)
			lightVal = float32(light) / 96
		case 3:
			modeStr = "Door"
			evtJoin := []string{data[16:24], data[24:32]}
			eventCount, _ := strconv.ParseInt(strings.Join(evtJoin, ""), 2, 16)
			switch eventCount {
			case 1:
				evtVal = "Calibration not done"
			case 3:
				evtVal = "Door closed"
			case 4:
				evtVal = "Door open"
			}
		case 4:
			modeStr = "Vibration"
			evtJoin := []string{data[16:24], data[24:32]}
			eventCount, _ := strconv.ParseInt(strings.Join(evtJoin, ""), 2, 16)
			switch eventCount {
			case 0:
				evtVal = "No vibration detected"
			case 1:
				evtVal = "Vibration detected"
			}
		case 5:
			modeStr = "Magnet"
			evtJoin := []string{data[16:24], data[24:32]}
			eventCount, _ := strconv.ParseInt(strings.Join(evtJoin, ""), 2, 16)
			switch eventCount {
			case 0:
				evtVal = "No magnet detected"
			case 1:
				evtVal = "Magnet detected"
			}
		default:
			modeStr = ""
		}

		modeData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "mode", UnitText: "", Value: modeStr}
		butData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "button", UnitText: "", Value: buttonStr}
		batData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "battery", UnitText: "volt", Value: batVal}
		obs.Values = append(obs.Values, modeData, butData, batData)
		if evtVal != "" { //Modes 3,4,5
			eventData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "event", UnitText: "", Value: evtVal}
			obs.Values = append(obs.Values, eventData)
		}

		switch mode {
		case 0:
			swRevData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "softwareRevision", UnitText: "", Value: fwRevVal}
			obs.Values = append(obs.Values, swRevData)
		case 1:
			//fmt.Println(humidity, "% RH")
			tempData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "temperature", UnitText: "celsius", Value: tempVal}
			humiData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "humidity", UnitText: "percent", Value: humiVal}
			obs.Values = append(obs.Values, tempData, humiData)
		case 2:
			//fmt.Println(light, "lux")
			lightData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "light", UnitText: "lux", Value: lightVal}
			obs.Values = append(obs.Values, lightData)
		case 3, 4, 5:
			evtData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "event", UnitText: "", Value: evtVal}
			obs.Values = append(obs.Values, evtData)
		}

	} else { //len: 24 exactly, 12 bytes
		fmt.Println("Sensit Daily Downlink Message")
		//TODO: Decode sensit downlink message
		dlData := models.QuantitativeValue{SemanticProperty: defp, Identifier: "downlink", UnitText: "", Value: true}
		obs.Values = append(obs.Values, dlData)
	}

	obs.Timestamp = msg.Timestamp
	obs.DeviceId = device.Id
	obs.Resolver = "sensitv3"
	return true, obs
}

func SigfoxSpotit(contxt *gin.Context, loc *sigfox.Location) (bool, *models.Geolocation, *models.Observation) {
	device, err := store.GetDeviceFromSigfoxId(contxt, loc.SigfoxId)
	if err != nil {
		fmt.Println("Wifi Enhancer Sigfox Device ID not found", err)
		return false, nil, nil
	}

	//Else, position is resolved
	//var wifiLoc sigfox.Location
	spotitLoc := &models.Geolocation{}
	spotitLoc.DeviceId = device.Id
	spotitLoc.Timestamp = loc.Timestamp
	spotitLoc.Source = "sigfox"
	spotitLoc.Latitude = loc.Latitude
	spotitLoc.Longitude = loc.Longitude
	spotitLoc.Radius = loc.Radius

	obs := &models.Observation{}
	defp := &models.SemanticProperty{Context: "spotit", Type: "location"}
	latVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "latitude", UnitText: "degrees", Value: loc.Latitude}
	lngVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "longitude", UnitText: "degrees", Value: loc.Longitude}
	accVal := models.QuantitativeValue{SemanticProperty: defp, Identifier: "accuracy", UnitText: "meters", Value: loc.Radius}
	obs.Values = append(obs.Values, latVal, lngVal, accVal)
	obs.Timestamp = loc.Timestamp
	obs.DeviceId = device.Id
	obs.Resolver = "spotit"

	return true, spotitLoc, obs
}

func DecodeLevelFrame(contxt *gin.Context, msg *sigfox.Message) (bool, *models.Observation) {
	//store.UpdateDeviceActivity(contxt, device.Id, 1)

	//if string(msg.Data[0:6]) == "000000"
	//humi, _ := strconv.ParseInt(msg.Data[0:32], 2, 16)

	obs := &models.Observation{}
	defp := &models.SemanticProperty{Context: "level", Type: "sensor"}
	temp := models.QuantitativeValue{SemanticProperty: defp, Identifier: "temperature", UnitText: "degrees", Value: 22}
	dist := models.QuantitativeValue{SemanticProperty: defp, Identifier: "distance", UnitText: "milimeter", Value: 1100}
	obs.Values = append(obs.Values, temp, dist)
	obs.Timestamp = msg.Timestamp
	//obs.DeviceId = device.Id
	obs.Resolver = "lockit"
	return true, obs
}

func decodeWisolGPSFrame(msg sigfox.Message) (models.Geolocation, string, int64, float64, bool) {
	fmt.Print("GPS frame: \t\t\t")
	var gpsLoc models.Geolocation
	var temperature float64
	var ori, moves int64
	var orientation string
	var status bool
	var latitude, longitude float64
	var latDeg, latMin, latSec float64
	var lngDeg, lngMin, lngSec float64

	isNorth, isEast := false, false
	if string(msg.Data[0:2]) == "4e" {
		isNorth = true
	}
	if string(msg.Data[10:12]) == "45" {
		isEast = true
	}

	if isNorth {
		fmt.Print("N:")
	} else {
		fmt.Print("S:")
	}

	valLatDeg, _ := strconv.ParseInt(msg.Data[2:4], 16, 8)
	latDeg = float64(valLatDeg)
	valLatMin, _ := strconv.ParseInt(msg.Data[4:6], 16, 8)
	latMin = float64(valLatMin)
	valLatSec, _ := strconv.ParseInt(msg.Data[6:8], 16, 8)
	latSec = float64(valLatSec)
	fmt.Print(latDeg, "° ", latMin, "m ", latSec, "s\t")

	latitude = float64(latDeg) + float64(latMin/60) + float64(latSec/3600)

	if isEast {
		fmt.Print("E:")
	} else {
		fmt.Print("W:")
	}
	valLngDeg, _ := strconv.ParseInt(msg.Data[10:12], 16, 8)
	lngDeg = float64(valLngDeg)
	valLngMin, _ := strconv.ParseInt(msg.Data[12:14], 16, 8)
	lngMin = float64(valLngMin)
	valLngSec, _ := strconv.ParseInt(msg.Data[14:16], 16, 8)
	lngSec = float64(valLngSec)
	fmt.Print(lngDeg, "° ", lngMin, "m ", lngSec, "s")

	longitude = float64(lngDeg) + float64(lngMin/60) + float64(lngSec/3600)

	fmt.Print("\t\t\t Lat: ", latitude, "\t Lng:", longitude)
	// Populating returned location
	gpsLoc.Latitude = latitude
	gpsLoc.Longitude = longitude
	gpsLoc.Timestamp = msg.Timestamp
	gpsLoc.Radius = 10
	gpsLoc.Source = "gps"

	/*if msg.Data[18:20] == "41" {
		status = true
	} else if msg.Data[18:20] == "56" {
		status = false
	}*/

	ori, _ = strconv.ParseInt(msg.Data[16:18], 16, 8)
	switch ori {
	case 1:
		orientation = "Haut"
	case 2:
		orientation = "Bas"
	case 3:
		orientation = "Droite"
	case 4:
		orientation = "Gauche"
	case 5:
		orientation = "Dos"
	case 6:
		orientation = "Ventre"
	}
	fmt.Println("Decoded orientation: ", ori, "\t : ", orientation)

	moves, _ = strconv.ParseInt(msg.Data[18:20], 16, 8)

	temperature, err := strconv.ParseFloat(msg.Data[20:22], 64)
	if err != nil {
		fmt.Println("Error while converting temperature main")
	}
	dec, err := strconv.ParseFloat(msg.Data[22:24], 64)
	if err != nil {
		fmt.Println("Error while converting temperature decimal")
	}

	temperature += dec * 0.01

	fmt.Println("\t\t", gpsLoc, "\t", temperature, '\t', status)
	return gpsLoc, orientation, moves, temperature, status
}

func CheckWifiCredit(c *gin.Context, subscription *models.Subscription) bool {
	//wifiCredit := config.GetInt(c, "plan_credit_wifi")
	wifiCredit := subscription.PlanCreditWifi
	fmt.Println("Wifi Organization credit:", wifiCredit)
	es := GetEmailSender(c)
	if wifiCredit > 0 {
		//config.Set(c, "plan_credit_wifi", wifiCredit-1)
		store.UpdateSubscription(c, subscription.Id, params.M{"$set": params.M{"plan_credit_wifi": wifiCredit - 1}})
		return true
	} else if wifiCredit == 0 {
		fmt.Println("Wifi Check Credit Organization no credit warning mails sent")
		appName := config.GetString(c, "mail_sender_name")
		subject := appName + ", your wifi token is empty, we give you 10 wifi"
		templateLink := "./templates/html/mail_token_empty.html"
		userData := models.EmailData{ReceiverMail: EmailSender.GetEmailParams(es).senderEmail, ReceiverName: EmailSender.GetEmailParams(es).senderName, Subject: subject, Body: "Wifi", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		adminData := models.EmailData{ReceiverMail: "contact@plugblocks.com", ReceiverName: "PlugBlocks Admin", Subject: subject, Body: "Wifi", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		EmailSender.SendEmailFromTemplate(es, c, subscription, &userData, templateLink)
		EmailSender.SendEmailFromTemplate(es, c, subscription, &adminData, templateLink)
		//config.Set(c, "plan_credit_wifi", -1)
		store.UpdateSubscription(c, subscription.Id, params.M{"$set": params.M{"plan_credit_wifi": wifiCredit - 1}})
		return false
	} else if wifiCredit > -10 {
		//config.Set(c, "plan_credit_wifi", -100)
		store.UpdateSubscription(c, subscription.Id, params.M{"$set": params.M{"plan_credit_wifi": wifiCredit - 100}})
		return false
	} else if wifiCredit == -100 {
		fmt.Println("Wifi Check Credit Organization no credit disable wifi sent")
		appName := config.GetString(c, "mail_sender_name")
		subject := appName + ", your wifi token is empty"
		templateLink := "./templates/html/mail_token_empty.html"
		userData := models.EmailData{ReceiverMail: EmailSender.GetEmailParams(es).senderEmail, ReceiverName: EmailSender.GetEmailParams(es).senderName, Subject: subject, Body: "Wifi", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		adminData := models.EmailData{ReceiverMail: "contact@plugblocks.com", ReceiverName: "PlugBlocks Admin", Subject: subject, Body: "Wifi", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		EmailSender.SendEmailFromTemplate(es, c, subscription, &userData, templateLink)
		EmailSender.SendEmailFromTemplate(es, c, subscription, &adminData, templateLink)
		//config.Set(c, "plan_credit_wifi", -1000)
		store.UpdateSubscription(c, subscription.Id, params.M{"$set": params.M{"plan_credit_wifi": wifiCredit - 1000}})
		return false
	}
	return false
}

func Wisol(contxt *gin.Context, sigfoxMessage *sigfox.Message) (bool, *models.Geolocation, *models.Observation) {
	device, err := store.GetDeviceFromSigfoxId(contxt, sigfoxMessage.SigfoxId)
	if err != nil {
		fmt.Println("Wifi Enhancer Sigfox Device ID not found", err)
		return false, nil, nil
	}

	geoloc := &models.Geolocation{}
	obs := &models.Observation{}
	locProp := &models.SemanticProperty{Context: "gps", Type: "location"}
	senProp := &models.SemanticProperty{Context: "gps", Type: "sensor"}

	if (string(sigfoxMessage.Data[0:2]) == "4e") || (string(sigfoxMessage.Data[0:2]) == "53") {
		if string(sigfoxMessage.Data[2:4]) != "00" {
			decodedGPSFrame, decodedOrientation, decodedMoves, decodedTemperature, _ := decodeWisolGPSFrame(*sigfoxMessage)
			geoloc = &decodedGPSFrame
			geoloc.DeviceId = device.Id

			latVal := models.QuantitativeValue{SemanticProperty: locProp, Identifier: "latitude", UnitText: "degrees", Value: geoloc.Latitude}
			lngVal := models.QuantitativeValue{SemanticProperty: locProp, Identifier: "longitude", UnitText: "degrees", Value: geoloc.Longitude}
			accVal := models.QuantitativeValue{SemanticProperty: locProp, Identifier: "accuracy", UnitText: "meters", Value: geoloc.Radius}
			tempVal := models.QuantitativeValue{SemanticProperty: senProp, Identifier: "temperature", UnitText: "celsius", Value: decodedTemperature}
			orVal := models.QuantitativeValue{SemanticProperty: senProp, Identifier: "orientation", UnitText: "", Value: decodedOrientation}
			movVal := models.QuantitativeValue{SemanticProperty: senProp, Identifier: "moves", UnitText: "", Value: decodedMoves}
			//staVal := models.QuantitativeValue{senProp, "status", "", status}
			obs.Values = append(obs.Values, latVal, lngVal, accVal, tempVal, orVal, movVal /*, staVal*/)
			obs.Timestamp = sigfoxMessage.Timestamp
			obs.DeviceId = device.Id
			obs.Resolver = "gps"

			fmt.Println("Wisol GPS Geoloc: ", geoloc, "Obs:", obs)
			return true, geoloc, obs

		} else { //No GPS, frame is empty
			sigfoxMessage.Data = "No GPS: " + sigfoxMessage.Data
			fmt.Println("Wisol No GPS Frame")
		}
	} else {
		subscription, err := store.GetOrganizationSubscription(contxt, device.OrganizationId)
		if err != nil {
			fmt.Println(err)
			return false, nil, nil
		}
		if !CheckWifiCredit(contxt, subscription) {
			return false, nil, nil
		}

		status, geoloc, obs := ResolveWifiPosition(contxt, sigfoxMessage)

		if status == false {
			fmt.Println("Error while resolving Wisol WiFi location for device: ", sigfoxMessage.SigfoxId, "at ", time.Unix(sigfoxMessage.Timestamp, 0))
			return false, nil, nil
		}

		fmt.Println("Wisol WiFi Geoloc: ", geoloc, "Obs:", obs)
		return true, geoloc, obs
	}
	return false, nil, nil
}