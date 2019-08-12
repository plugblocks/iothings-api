package services

//TODO: Store API keys on fleet/device

/*
func SigfoxAPIGetDeviceMessages(c *gin.Context) (sigfoxMessages *sigfox.SfxApiMessages) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://backend.sigfox.com/api/devices/"+c.Param("sigfoxId")+"/messages", nil)
	//TODO: fix config interface bug

	var key models.APIKey
	key = config.RetrieveSigfoxAPIKey(c.Param("sigfoxId"))

	req.SetBasicAuth(key.SigfoxKey, key.SigfoxSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var messages = new(sigfox.SfxApiMessages)
	err = json.Unmarshal(bodyResp, &messages)
	if err != nil {
		panic(err.Error())
	}
	return messages
}*/
/*
func SigfoxAPIMessagesInflating(c *gin.Context) {
	//TODO: fix config interface bug
	remoteThingsUrl := "http://localhost:4000/v1/"
	thingsAuthParams := map[string]string{"email": "admin@plugblocks.fr", "password": "adchapwd"}

	//Authenticate to things-api : POST /auth
	//Create Device if not exists
	//Iterate through data to fire POST request (delay) to not hit rate limits

	//thingsAuthReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	sigfoxAPIMessages := SigfoxAPIGetDeviceMessages(c)
	// Step 1: Authenticate User
	jsonValThingsAuth, _ := json.Marshal(thingsAuthParams)
	thingsAuthReq, _ := http.NewRequest("POST", remoteThingsUrl+"auth/", bytes.NewBuffer(jsonValThingsAuth))
	//thingsAuthReq, _ := http.NewRequest("POST", "http://requestbin.fullcontact.com/1ewrhf51", bytes.NewBuffer(jsonValue))
	thingsAuthReq.Header.Set("Content-Type", "application/json")
	authRes, _ := client.Do(thingsAuthReq)
	thingsApiAuthResp, err := ioutil.ReadAll(authRes.Body)

	var authResp = new(models.AuthResp)
	err = json.Unmarshal(thingsApiAuthResp, &authResp)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("User: ", authResp.User.Firstname, authResp.User.Lastname, " succesfully authenticated")

	// Step 2: Create device and checks if it already exists
	thingsDeviceParams := map[string]string{"name": "Sigfox Device ID: " + c.Param("sigfoxId"), "sigfox_id": c.Param("sigfoxId")}
	jsonValDevice, _ := json.Marshal(thingsDeviceParams)
	thingsDeviceReq, _ := http.NewRequest("POST", remoteThingsUrl+"devices/", bytes.NewBuffer(jsonValDevice))
	thingsDeviceReq.Header.Set("Content-Type", "application/json")
	thingsDeviceReq.Header.Set("Authorization", "Bearer "+authResp.Token)
	deviceResp, _ := client.Do(thingsDeviceReq)
	thingsDeviceResp, err := ioutil.ReadAll(deviceResp.Body)

	var thingsDev = new(models.ThingsDevice)
	err = json.Unmarshal(thingsDeviceResp, &thingsDev)
	if err != nil {
		panic(err.Error())
	}
	if deviceResp.StatusCode == http.StatusCreated {
		fmt.Println("Device successfully existed")
	} else if deviceResp.StatusCode == http.StatusConflict {
		fmt.Println("Device already existed")
	}

	//fmt.Println("Device creation code: "+ string(deviceResp.StatusCode) + " Device ID: "+ thingsDev.Id)
	//fmt.Println("Device "+thingsDev.Id+" successfully created")

	// Step 3: Iterate through data to fire POST request (delay) to not hit rate limits

	fmt.Println(sigfoxAPIMessages.Messages)
	for i, mes := range sigfoxAPIMessages.Messages {
		fmt.Print("mes" + strconv.Itoa(i) + ": " + mes.SigfoxId + " " + strconv.Itoa(int(mes.Time)) + " " + strconv.Itoa(int(mes.SequenceNbr)) + " " +
			mes.Data + " " + mes.Snr + " " + fmt.Sprintf("%f", mes.CompLoc.Latitude) + " " + fmt.Sprintf("%f", mes.CompLoc.Longitude))
		fmt.Println(mes.CompLoc.Radius)
		fmt.Println(mes.CompLoc.Source)

		messageToSend := models.ThingMessage{}
		sigSnr, err := strconv.ParseFloat(mes.Snr, 64)
		if err != nil {
			panic(err.Error())
		}

		//IMPV: Use mes.LinkQuality = GOOD, AVERAGE ... as a RSSI

		messageToSend.SigfoxId = mes.SigfoxId
		messageToSend.Timestamp = mes.Time
		messageToSend.FrameNumber = mes.SequenceNbr
		messageToSend.Data = mes.Data
		messageToSend.Snr = sigSnr
		messageToSend.Resolver = c.Param("resolver")//wifi, sensit ...

		jsonValMessage, _ := json.Marshal(messageToSend)
		thingsMessageReq, _ := http.NewRequest("POST", remoteThingsUrl+"sigfox/message", bytes.NewBuffer(jsonValMessage))
		thingsMessageReq.Header.Set("Content-Type", "application/json")
		//thingsMessageReq.Header.Set("Authorization", "Bearer "+authResp.Token)
		messageResp, err := client.Do(thingsMessageReq)
		//TODO fix bug by reusing client instead of creating a new one
		if err != nil {
			panic(err.Error())
		}

		if messageResp.StatusCode == http.StatusCreated {
			fmt.Println("Message successfully created")
		} else if messageResp.StatusCode == http.StatusPartialContent {
			fmt.Println("Message created but Sigfox device Id not exists")
		} else {
			fmt.Println("Unknown message Error")
		}

		locationToSend := models.ThingLocation{mes.SigfoxId, mes.SequenceNbr, mes.Time, mes.CompLoc.Latitude,
			mes.CompLoc.Longitude, float64(mes.CompLoc.Radius), true, false, false}

		jsonValLocation, _ := json.Marshal(locationToSend)
		thingsLocationReq, _ := http.NewRequest("POST", remoteThingsUrl+"sigfox/location", bytes.NewBuffer(jsonValLocation))
		thingsLocationReq.Header.Set("Content-Type", "application/json")
		locationResp, _ := client.Do(thingsLocationReq)

		if locationResp.StatusCode == http.StatusCreated {
			fmt.Println("Location successfully created")
		} else if locationResp.StatusCode == http.StatusPartialContent {
			fmt.Println("Location created but Sigfox device Id not exists")
		} else {
			fmt.Println("Unknown message Error")
		}
	}

	// Step 4: Logout

	//c.JSON(http.StatusOK, authResp.User.Email)
	c.JSON(http.StatusOK, "Okay")
}*/
