package controllers

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/store"
)

type GroupController struct {

}
func NewGroupController() GroupController {
	return GroupController{}
}

func (gc GroupController) CreateGroup(c *gin.Context) {
	group := &models.Group{}

	err := c.BindJSON(group)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateGroup(c, group); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, group)
}