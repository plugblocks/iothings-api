package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateOrder(organizationId string, order *models.Order) error {
	session := db.Session.Copy()
	defer session.Close()
	orders := db.C(models.OrdersCollection).With(session)

	order.Id = bson.NewObjectId().Hex()
	order.OrganizationId = organizationId

	err := orders.Insert(order)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "order_creation_failed", "Failed to insert the order in the database", err)
	}

	return nil
}

func (db *mongo) GetAllOrders(organizationId string) ([]models.Order, error) {
	session := db.Session.Copy()
	defer session.Close()
	orderCollection := db.C(models.OrdersCollection).With(session)

	orders := []models.Order{}
	err := orderCollection.Find(bson.M{"organization_id": organizationId}).All(&orders)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_orders_failed", "Failed to get the orders: "+err.Error(), err)
	}

	return orders, nil
}

func (db *mongo) GetOrderById(organizationId string, id string) (*models.Order, error) {
	session := db.Session.Copy()
	defer session.Close()
	orderCollection := db.C(models.OrdersCollection).With(session)

	order := &models.Order{}
	err := orderCollection.Find(bson.M{"_id": id, "organization_id": organizationId}).One(order)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "order_not_found", "Could not find the order", err)
	}

	return order, nil
}

func (db *mongo) UpdateOrder(organizationId string, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	orders := db.C(models.OrdersCollection).With(session)

	err := orders.Update(bson.M{"_id": id, "organization_id": organizationId}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "order_update_failed", "Failed to update the orders: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteOrder(organizationId string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	orders := db.C(models.OrdersCollection).With(session)

	err := orders.Remove(bson.M{"_id": id, "organization_id": organizationId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "order_delete_failed", "Failed to delete the order", err)
	}

	return nil
}
