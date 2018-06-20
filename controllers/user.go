package controllers

import (
	"net/http"

	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"

	"gopkg.in/gin-gonic/gin.v1"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

func (uc UserController) GetUser(c *gin.Context) {
	user, err := store.FindUserById(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	c.JSON(http.StatusOK, user.Sanitize())
}

func (uc UserController) CreateUser(c *gin.Context) {
	user := &models.User{}

	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateUser(c, user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	//TODO: Remove business logic from controller
	appName := config.GetString(c, "sendgrid_name")
	subject := "Welcome to " + appName + "! Account confirmation"
	templateLink := "./templates/html/mail_activate_account.html"
	services.GetEmailSender(c).SendEmailFromTemplate(user, subject, templateLink)

	c.JSON(http.StatusCreated, user.Sanitize())
}

func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uc UserController) GetUsers(c *gin.Context) {
	users, err := store.GetUsers(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc UserController) ImpersonateUser(c *gin.Context) {

}
