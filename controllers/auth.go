package controllers

import (
	"fmt"
	"net/http"

	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func NewAuthController() AuthController {
	return AuthController{}
}

func (ac AuthController) UserAuthentication(c *gin.Context) {
	userInput := models.User{}
	if err := c.Bind(&userInput); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	user, err := store.FindUser(c, params.M{"email": userInput.Email})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_does_not_exist", "User does not exist", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("incorrect_password", "Password is not correct", err))
		return
	}

	if !user.Active {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_needs_activation", "User needs to be activated via email", nil))
		return
	}

	//Read base64 private key
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	accessToken, err := helpers.GenerateAccessToken(encodedKey, user.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("token_generation_failed", "Could not generate the access token", err))
		return
	}

	fmt.Println("User authenticated: ", user)

	c.JSON(http.StatusOK, gin.H{"token": accessToken, "user": user.Sanitize()})
}
