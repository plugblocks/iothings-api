package controllers

import (
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
	"gitlab.com/plugblocks/iothings-api/models"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
)

type OrganizationController struct{}

func NewOrganizationController() OrganizationController {
	return OrganizationController{}
}

func (oc OrganizationController) GetOrganizations(c *gin.Context) {
	organizations, err := store.GetAllOrganizations(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, organizations)
}

func (oc OrganizationController) GetOrganizationById(c *gin.Context) {
	id := c.Param("id")

	organization, err := store.GetOrganizationById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (oc OrganizationController) CreateOrganization(c *gin.Context) {
	organization := &models.Organization{}

	if err := c.BindJSON(organization); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateOrganization(c, organization); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, organization)
}

func (oc OrganizationController) UpdateOrganization(c *gin.Context) {
	organization := &models.Organization{}
	id := c.Param("id")

	if err := c.BindJSON(organization); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateOrganization(c, id, params.M{"$set": organization}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, organization)
}

func (fc OrganizationController) DeleteOrganization(c *gin.Context) {
	err := store.DeleteOrganization(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)

}
