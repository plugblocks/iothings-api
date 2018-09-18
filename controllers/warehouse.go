package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
	"net/http"
)

type WarehouseController struct{}

func NewWarehouseController() WarehouseController {
	return WarehouseController{}
}

func (oc WarehouseController) GetWarehouses(c *gin.Context) {
	warehouses, err := store.GetAllWarehouses(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

func (oc WarehouseController) GetWarehouseById(c *gin.Context) {
	id := c.Param("id")

	warehouse, err := store.GetWarehouseById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

func (oc WarehouseController) CreateWarehouse(c *gin.Context) {
	warehouse := &models.Warehouse{}

	if err := c.BindJSON(warehouse); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateWarehouse(c, warehouse); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

func (oc WarehouseController) EditWarehouse(c *gin.Context) {
	warehouse := &models.Warehouse{}
	id := c.Param("id")

	if err := c.BindJSON(warehouse); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.UpdateWarehouse(c, id, params.M{"$set": warehouse}); err != nil {
		c.Error(err)
		c.Abort()
	}

	c.JSON(http.StatusOK, warehouse)
}

func (oc WarehouseController) DeleteWarehouse(c *gin.Context) {
	err := store.DeleteWarehouse(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
