package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateOrder(c context.Context, record *models.Order) error {
	return FromContext(c).CreateOrder(Current(c).OrganizationId, record)
}

func GetAllOrders(c context.Context) ([]models.Order, error) {
	return FromContext(c).GetAllOrders(Current(c).OrganizationId)
}

func GetOrderById(c context.Context, id string) (*models.Order, error) {
	return FromContext(c).GetOrderById(Current(c).OrganizationId, id)
}

func UpdateOrder(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateOrder(Current(c).OrganizationId, id, params)
}

func DeleteOrder(c context.Context, id string) error {
	return FromContext(c).DeleteOrder(Current(c).OrganizationId, id)
}
