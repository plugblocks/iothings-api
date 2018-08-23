package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
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

	if !services.GetEmailSender(c).CheckMailCredit(c) {
		return
	}
	s := services.GetEmailSender(c)
	data := models.EmailData{User: user, Subject: subject, ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
	s.SendEmailFromTemplate(&data, templateLink)
}
