package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateSubscription(subscription *models.Subscription) error {
	session := db.Session.Copy()
	defer session.Close()
	subscriptions := db.C(models.SubscriptionsCollection).With(session)

	subscription.BeforeCreate()

	err := subscriptions.Insert(subscription)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "subscription_creation_failed", "Failed to insert the subscription in the database", err)
	}

	return nil
}

func (db *mongo) GetSubscriptions() ([]models.Subscription, error) {
	session := db.Session.Copy()
	defer session.Close()
	subscriptionCollection := db.C(models.SubscriptionsCollection).With(session)

	subscriptions := []models.Subscription{}
	err := subscriptionCollection.Find(bson.M{}).All(&subscriptions)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_subscriptions_failed", "Failed to get the subscriptions: "+err.Error(), err)
	}

	return subscriptions, nil
}

func (db *mongo) GetSubscription(id string) (*models.Subscription, error) {
	session := db.Session.Copy()
	defer session.Close()
	subscriptionCollection := db.C(models.SubscriptionsCollection).With(session)

	subscription := &models.Subscription{}
	err := subscriptionCollection.Find(bson.M{"_id": id}).One(subscription)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "subscription_not_found", "Could not find the subscription", err)
	}

	return subscription, nil
}

func (db *mongo) UpdateSubscription(id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	subscriptions := db.C(models.SubscriptionsCollection).With(session)

	err := subscriptions.Update(bson.M{"_id": id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "subscription_update_failed", "Failed to update the subscriptions: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteSubscription(id string) error {
	session := db.Session.Copy()
	defer session.Close()
	subscriptions := db.C(models.SubscriptionsCollection).With(session)

	err := subscriptions.Remove(bson.M{"_id": id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "subscription_delete_failed", "Failed to delete the subscription", err)
	}

	return nil
}
