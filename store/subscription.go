package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateSubscription(c context.Context, record *models.Subscription) error {
	return FromContext(c).CreateSubscription(record)
}

func GetSubscriptions(c context.Context) ([]models.Subscription, error) {
	return FromContext(c).GetSubscriptions()
}

func GetSubscription(c context.Context, id string) (*models.Subscription, error) {
	return FromContext(c).GetSubscription(id)
}

func UpdateSubscription(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateSubscription(id, params)
}

func DeleteSubscription(c context.Context, id string) error {
	return FromContext(c).DeleteSubscription(id)
}
