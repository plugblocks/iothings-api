package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateWarehouse(c context.Context, record *models.Warehouse) error {
	return FromContext(c).CreateWarehouse(Current(c).OrganizationId, record)
}

func GetAllWarehouses(c context.Context) ([]models.Warehouse, error) {
	return FromContext(c).GetAllWarehouses(Current(c).OrganizationId)
}

func GetWarehouseById(c context.Context, id string) (*models.Warehouse, error) {
	return FromContext(c).GetWarehouseById(Current(c).OrganizationId, id)
}

func UpdateWarehouse(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateWarehouse(Current(c).OrganizationId, id, params)
}

func DeleteWarehouse(c context.Context, id string) error {
	return FromContext(c).DeleteWarehouse(Current(c).OrganizationId, id)
}
