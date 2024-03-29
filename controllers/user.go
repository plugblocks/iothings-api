package controllers

import (
	"net/http"

	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"

	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/services"
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

	appName := config.GetString(c, "mail_sender_name")
	subject := "Bienvenue sur " + appName + "! Merci de confirmer votre compte"
	templateLink := "./templates/html/mail_utilisateur_activation.html"

	if err := store.CreateUser(c, user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	s := services.GetEmailSender(c)
	data := models.EmailData{ReceiverMail: user.Email, ReceiverName: user.Firstname + " " + user.Lastname, User: user, Subject: subject, ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}

	err := s.SendEmailFromTemplate(c, &data, templateLink)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_credit_spent", "Your mail credit is spent", err))
		return
	}

	c.JSON(http.StatusCreated, user.Sanitize())
}

func (uc UserController) DeleteUser(c *gin.Context) {
	err := store.DeleteUser(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uc UserController) ChangeLanguage(c *gin.Context) {
	if err := store.ChangeLanguage(c, c.Param("id"), c.Param("language")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	//c.JSON(http.StatusOK, nil)

	/*user, err := store.FindUserById(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	vars := gin.H{
		"User": user,
		"AppName": config.GetString(c, "mail_sender_name"),
		"AppUrl": config.GetString(c, "front_url"),
	}

	c.HTML(http.StatusOK, "./templates/html/page_account_activated.html", vars)*/

	c.Redirect(http.StatusMovedPermanently, "https://"+config.GetString(c, "front_url"))
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

func (uc UserController) AssignOrganization(c *gin.Context) {
	userId := c.Param("id")
	organizationId := c.Param("organization_id")

	err := store.AssignOrganization(c, userId, organizationId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uc UserController) GetUserOrganization(c *gin.Context) {
	user := store.Current(c)
	var err error

	if user.Admin {
		user, err = store.FindUserById(c, c.Param("id"))
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	organization, err := store.GetUserOrganization(c, user)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, organization)
}
