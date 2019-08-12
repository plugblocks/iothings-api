package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store"
)

type SubscriptionController struct {
}

func NewSubscriptionController() SubscriptionController {
	return SubscriptionController{}
}

func (sc SubscriptionController) CreateSubscription(c *gin.Context) {
	subscription := &models.Subscription{}

	err := c.BindJSON(subscription)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateSubscription(c, subscription); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (sc SubscriptionController) GetSubscriptions(c *gin.Context) {
	subscriptions, err := store.GetSubscriptions(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (sc SubscriptionController) GetSubscription(c *gin.Context) {
	subscription, err := store.GetSubscription(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (sc SubscriptionController) UpdateSubscription(c *gin.Context) {
	newSubscription := models.Subscription{}

	err := c.BindJSON(&newSubscription)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	oldSubscription, err := store.GetSubscription(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("subscription_not_found", "Failed to find subscription id", err))
		return
	}

	if newSubscription.OrganizationId == "" {
		newSubscription.OrganizationId = oldSubscription.OrganizationId
	}

	newSubscription.Active = oldSubscription.Active
	/*changes := params.M{"$set": params.M{"organization_id": newSubscription.OrganizationId, "customer_id": newSubscription.CustomerId,
	"name": newSubscription.Name, "ble_mac": newSubscription.BleMac, "wifi_mac": newSubscription.WifiMac, "sigfox_id": newSubscription.SigfoxId,
	"last_access": oldSubscription.LastAccess, "active": oldSubscription.Active}}*/

	err = store.UpdateSubscription(c, c.Param("id"), params.M{"$set": newSubscription})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (sc SubscriptionController) DeleteSubscription(c *gin.Context) {
	err := store.DeleteSubscription(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("subscription_delete_failed", "Failed to delete the subscription", err))
		return
	}

	c.JSON(http.StatusOK, nil)
}
