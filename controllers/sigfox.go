package controllers

import (
	"github.com/gin-gonic/gin"
)

type SigfoxController struct {
}

func NewSigfoxController() SigfoxController {
	return SigfoxController{}
}

func (sc SigfoxController) CreateMessage(c *gin.Context) {
	//sigfoxMessage := &models.SigfoxMessage{}
	//
	//err := c.BindJSON(sigfoxMessage)
	//if err != nil {
	//	c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
	//	return
	//}
	//
	//err = store.CreateMessage(c, sigfoxMessage)
	//if err != nil {
	//	c.Error(err)
	//	c.Abort()
	//	return
	//}
	//
	//c.JSON(http.StatusCreated, sigfoxMessage)
}

func (sc SigfoxController) CreateLocation(c *gin.Context) {
	//location := &models.Location{}
	//
	//err := c.BindJSON(location)
	//if err != nil {
	//	c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
	//	return
	//}
	//
	//err = store.CreateLocation(c, location)
	//if err != nil {
	//	c.Error(err)
	//	c.Abort()
	//	return
	//}
	//
	//c.JSON(http.StatusCreated, location)
}
