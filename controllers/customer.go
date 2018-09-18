package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type CustomerController struct{}

func NewCustomerController() CustomerController {
	return CustomerController{}
}

func (oc CustomerController) GetCustomers(c *gin.Context) {
	customers, err := store.GetAllCustomers(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (oc CustomerController) GetCustomerById(c *gin.Context) {
	id := c.Param("id")

	customer, err := store.GetCustomerById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (oc CustomerController) CreateCustomer(c *gin.Context) {
	customer := &models.Customer{}

	if err := c.BindJSON(customer); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateCustomer(c, customer); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (oc CustomerController) EditCustomer(c *gin.Context) {
	customer := &models.Customer{}
	id := c.Param("id")

	if err := c.BindJSON(customer); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateCustomer(c, id, params.M{"$set": customer}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, customer)
}

func (oc CustomerController) DeleteCustomer(c *gin.Context) {
	err := store.DeleteCustomer(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
