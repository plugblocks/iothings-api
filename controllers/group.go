package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type GroupController struct{}

func NewGroupController() GroupController {
	return GroupController{}
}

func (gtc GroupController) GetGroups(c *gin.Context) {
	groups, err := store.GetAllGroups(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (gtc GroupController) GetGroupById(c *gin.Context) {
	id := c.Param("id")

	group, err := store.GetGroupById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, group)
}

func (gtc GroupController) CreateGroup(c *gin.Context) {
	group := &models.Group{}

	if err := c.BindJSON(group); err != nil {
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

func (gtc GroupController) EditGroup(c *gin.Context) {
	group := &models.Group{}
	id := c.Param("id")

	if err := c.BindJSON(group); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateGroup(c, id, params.M{"$set": group}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, group)
}

func (fc GroupController) DeleteGroup(c *gin.Context) {
	err := store.DeleteGroup(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)

}
