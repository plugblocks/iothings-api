package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type AlertController struct{}

func NewAlertController() AlertController {
	return AlertController{}
}

func (ac AlertController) GetAlerts(c *gin.Context) {
	alerts, err := store.GetAlerts(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (ac AlertController) CreateAlert(c *gin.Context) {
	alert := &models.Alert{}

	err := c.BindJSON(alert)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateAlert(c, alert); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, alert)
}

func (ac AlertController) GetAlert(c *gin.Context) {
	alert, err := store.GetAlert(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, alert)
}

func (ac AlertController) GetFleetAlerts(c *gin.Context) {
	alerts, err := store.GetFleetAlerts(c, c.Param("fleetId"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (ac AlertController) GetDeviceAlerts(c *gin.Context) {
	alerts, err := store.GetDeviceAlerts(c, c.Param("deviceId"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (ac AlertController) UpdateAlert(c *gin.Context) {
	alert := &models.Alert{}

	if err := c.BindJSON(alert); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateAlert(c, c.Param("id"), params.M{"$set": alert}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, alert)
}

func (ac AlertController) DeleteAlert(c *gin.Context) {
	err := store.DeleteAlert(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (ac AlertController) TextTest(c *gin.Context) {
	user, err := store.FindUserById(c, c.Param("userId"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	s := services.GetTextSender(c)
	data := models.TextData{PhoneNumber: user.Phone, Subject: "Text Alert", Message: "You just received an alert, check at: https://demo.plugblocks.com"}
	err = s.SendText(data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("text_send_error", "Text sending error", err))
		return
	}

	c.JSON(http.StatusOK, nil)
}
