package store

import (
	"gitlab.com/plugblocks/iothings-api/models"
	"golang.org/x/net/context"
)

const (
	CurrentKey = "currentUser"
	StoreKey   = "store"
)

type Setter interface {
	Set(string, interface{})
}

func Current(c context.Context) *models.User {
	return c.Value(CurrentKey).(*models.User)
}

func CurrentCustomer(c context.Context) *models.Customer {
	return c.Value(CurrentKey).(*models.Customer)
}

func ToContext(c Setter, store Store) {
	c.Set(StoreKey, store)
}

func FromContext(c context.Context) Store {
	return c.Value(StoreKey).(Store)
}
