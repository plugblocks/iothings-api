package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateWarehouse(organizationId string, warehouse *models.Warehouse) error {
	session := db.Session.Copy()
	defer session.Close()
	warehouses := db.C(models.WarehousesCollection).With(session)

	warehouse.Id = bson.NewObjectId().Hex()
	warehouse.OrganizationId = organizationId

	err := warehouses.Insert(warehouse)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "warehouse_creation_failed", "Failed to insert the warehouse in the database", err)
	}

	return nil
}

func (db *mongo) GetAllWarehouses(organizationId string) ([]models.Warehouse, error) {
	session := db.Session.Copy()
	defer session.Close()
	warehouseCollection := db.C(models.WarehousesCollection).With(session)

	warehouses := []models.Warehouse{}
	err := warehouseCollection.Find(bson.M{"organization_id": organizationId}).All(&warehouses)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_warehouses_failed", "Failed to get the warehouses: "+err.Error(), err)
	}

	return warehouses, nil
}

func (db *mongo) GetWarehouseById(organizationId string, id string) (*models.Warehouse, error) {
	session := db.Session.Copy()
	defer session.Close()
	warehouseCollection := db.C(models.WarehousesCollection).With(session)

	warehouse := &models.Warehouse{}
	err := warehouseCollection.Find(bson.M{"_id": id, "organization_id": organizationId}).One(warehouse)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "warehouse_not_found", "Could not find the warehouse", err)
	}

	return warehouse, nil
}

func (db *mongo) UpdateWarehouse(organizationId string, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	warehouses := db.C(models.WarehousesCollection).With(session)

	err := warehouses.Update(bson.M{"_id": id, "organization_id": organizationId}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "warehouse_update_failed", "Failed to update the warehouses: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteWarehouse(organizationId string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	warehouses := db.C(models.WarehousesCollection).With(session)

	err := warehouses.Remove(bson.M{"_id": id, "organization_id": organizationId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "warehouse_delete_failed", "Failed to delete the warehouse", err)
	}

	return nil
}
