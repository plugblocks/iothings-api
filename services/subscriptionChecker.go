package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
	PlanType:        "onpremise",
	PlanExpiration:  "1538388000", //1/10/2018
	PlanCreditMails: "100000",
	PlanCreditTexts: "100",
	PlanCreditWifi:  "2000",
*/

func CheckSubscription(conf *viper.Viper) {
	fmt.Println("Checking subscription")
	remoteCheckerUrl := "https://adminapi.plugblocks.com/"
	client := &http.Client{}
	data := map[string]string{
		"client":           conf.GetString("mail_sender_name"),
		"plan_expiration":  strconv.Itoa(conf.GetInt("plan_expiration")),
		"plan_credit_mail": strconv.Itoa(conf.GetInt("plan_credit_mail")),
		"plan_credit_text": strconv.Itoa(conf.GetInt("plan_credit_text")),
		"plan_credit_wifi": strconv.Itoa(conf.GetInt("plan_credit_wifi")),
	}
	upstream, _ := json.Marshal(data)
	checkerReq, _ := http.NewRequest("POST", remoteCheckerUrl+"consumption/", bytes.NewBuffer(upstream))
	checkerReq.Header.Set("Content-Type", "application/json")
	checkerReq.Header.Set("Authorization", "Bearer "+conf.GetString("rsa_private"))
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
	if checkerResp.StatusCode == http.StatusOK {
		fmt.Println("Subscription check OK")
		conf.Set("", subscription.PlanExpiration)
	} else if checkerResp.StatusCode == http.StatusExpectationFailed {
		fmt.Println("Subscription check failed, one of credits is empty")
	} else if checkerResp.StatusCode == http.StatusPaymentRequired {
		fmt.Println("Subscription check failed, stopping API")
		conf.Set("", subscription.Active)
	}
}
