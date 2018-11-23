package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
)

type NotificationController struct {
}

func NewNotificationController() NotificationController {
	return NotificationController{}
}

func (nc NotificationController) SendAlertMail(c *gin.Context, user *models.User, device *models.Device, observation *models.Observation) {
	appName := config.GetString(c, "mail_sender_name")
	subject := "Alert for device for" + appName
	templateLink := "./templates/html/mail_alert.html"

	s := services.GetEmailSender(c)
	data := models.EmailData{ReceiverMail: user.Email, ReceiverName: user.Firstname + " " + user.Lastname, Subject: subject, ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}

	subscription, err := store.GetOrganizationSubscription(c, user.OrganizationId)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.SendEmailFromTemplate(c, subscription, &data, templateLink)
}
