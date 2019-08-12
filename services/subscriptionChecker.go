package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/models"
	"io/ioutil"
	"net/http"
)

/*
	PlanType:        "onpremise",
	PlanExpiration:  "1538388000", //1/10/2018
	PlanCreditMails: "100000",
	PlanCreditTexts: "100",
	PlanCreditWifi:  "2000",
*/

func CheckSubscription(router *gin.Engine, conf *viper.Viper, ctxt *gin.Context) {
	fmt.Println("Checking subscription")
	//remoteCheckerUrl := "https://adminapi.plugblocks.com/v1/"
	remoteCheckerUrl := "http://localhost:6000/v1/"
	client := &http.Client{}

	//devicesNbr, err := store.CountDevices(ctxt)
	//fleetsNbr, err := store.CountFleets(ctxt)

	//fmt.Println("Fleets: ", fleetsNbr, "Devices:", devicesNbr)

	//Step 1: Admin API Auth
	type AuthResp struct {
		Token string      `json:"token" bson:"token" valid:"-"`
		User  models.User `json:"user" bson:"user" valid:"-"`
	}

	thingsAuthParams := map[string]string{"email": "admin@plugblocks.fr", "password": "adchapwd"}
	jsonValThingsAuth, _ := json.Marshal(thingsAuthParams)
	thingsAuthReq, _ := http.NewRequest("POST", remoteCheckerUrl+"auth/", bytes.NewBuffer(jsonValThingsAuth))
	thingsAuthReq.Header.Set("Content-Type", "application/json")
	authRes, err := client.Do(thingsAuthReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	thingsApiAuthResp, err := ioutil.ReadAll(authRes.Body)

	var authResp = new(AuthResp)
	err = json.Unmarshal(thingsApiAuthResp, &authResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Step 2: Admin API check
	data := models.Subscription{PlanType: conf.GetString("plan_type"), /*PlanExpiration:conf.GetInt("plan_expiration"),*/
		PlanCreditMails: conf.GetInt("plan_credit_mail"), PlanCreditTexts: conf.GetInt("plan_credit_text"), PlanCreditWifi: conf.GetInt("plan_credit_wifi")}

	upstream, _ := json.Marshal(data)
	checkerReq, _ := http.NewRequest("POST", remoteCheckerUrl+"telemetry/check/"+conf.GetString("client_name"), bytes.NewBuffer(upstream))
	checkerReq.Header.Set("Content-Type", "application/json")
	checkerReq.Header.Set("Authorization", "Bearer "+authResp.Token)
	checkerResp, err := client.Do(checkerReq)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkerResp.StatusCode == http.StatusNotFound {
		fmt.Println("Subscription check API KO")
		return
	}

	adminCheckerResp, err := ioutil.ReadAll(checkerResp.Body)

	var subscription = new(models.Subscription)
	err = json.Unmarshal(adminCheckerResp, &subscription)
	if err != nil {
		fmt.Println(err.Error())
	}
	conf.Set("plan_expiration", subscription.PlanExpiration)
	conf.Set("plan_credit_mail", subscription.PlanCreditMails)
	conf.Set("plan_credit_text", subscription.PlanCreditTexts)
	conf.Set("plan_credit_wifi", subscription.PlanCreditWifi)

	if checkerResp.StatusCode == http.StatusOK {
		fmt.Println("Subscription check OK")
		conf.Set("plan_expired", false)
	} else if checkerResp.StatusCode == http.StatusExpectationFailed {
		fmt.Println("Subscription check failed, one of credits is empty")
	} else if checkerResp.StatusCode == http.StatusPaymentRequired || subscription.Active == false {
		fmt.Println("Subscription check failed, stopping API")
		conf.Set("plan_expired", true)
		//endless.ListenAndServe(conf.GetString("host_address"), router)
		/*srv := &http.Server{
			Addr:    conf.GetString("host_address"),
			Handler: router,
		}
		if err := srv.Shutdown(ctxt); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")*/
	}
}
