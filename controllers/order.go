package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type OrderController struct{}

func NewOrderController() OrderController {
	return OrderController{}
}

func (oc OrderController) GetOrders(c *gin.Context) {
	orders, err := store.GetAllOrders(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (oc OrderController) GetOrderById(c *gin.Context) {
	id := c.Param("id")

	order, err := store.GetOrderById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oc OrderController) CreateOrder(c *gin.Context) {
	order := &models.Order{}

	if err := c.BindJSON(order); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateOrder(c, order); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (oc OrderController) EditOrder(c *gin.Context) {
	order := &models.Order{}
	id := c.Param("id")

	if err := c.BindJSON(order); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateOrder(c, id, params.M{"$set": order}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, order)
}

func (oc OrderController) DeleteOrder(c *gin.Context) {
	err := store.DeleteOrder(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
