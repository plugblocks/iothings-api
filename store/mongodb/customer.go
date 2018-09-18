package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateCustomer(organizationId string, customer *models.Customer) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	customer.Id = bson.NewObjectId().Hex()
	customer.OrganizationId = organizationId

	err := customers.Insert(customer)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_creation_failed", "Failed to insert the customer in the database", err)
	}

	return nil
}

func (db *mongo) GetAllCustomers(organizationId string) ([]models.Customer, error) {
	session := db.Session.Copy()
	defer session.Close()
	customerCollection := db.C(models.CustomersCollection).With(session)

	customers := []models.Customer{}
	err := customerCollection.Find(bson.M{"organization_id": organizationId}).All(&customers)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_customers_failed", "Failed to get the customers: "+err.Error(), err)
	}

	return customers, nil
}

func (db *mongo) GetCustomerById(organizationId string, id string) (*models.Customer, error) {
	session := db.Session.Copy()
	defer session.Close()
	customerCollection := db.C(models.CustomersCollection).With(session)

	customer := &models.Customer{}
	err := customerCollection.Find(bson.M{"_id": id, "organization_id": organizationId}).One(customer)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customer_not_found", "Could not find the customer", err)
	}

	return customer, nil
}

func (db *mongo) UpdateCustomer(organizationId string, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	err := customers.Update(bson.M{"_id": id, "organization_id": organizationId}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_update_failed", "Failed to update the customers: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteCustomer(organizationId string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	err := customers.Remove(bson.M{"_id": id, "organization_id": organizationId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_delete_failed", "Failed to delete the customer", err)
	}

	return nil
}
