package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateCustomer(customer *models.Customer) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	customer.Id = bson.NewObjectId().Hex()
	err := customer.BeforeCreate()
	if err != nil {
		return err
	}

	if count, _ := customers.Find(bson.M{"email": customer.Email}).Count(); count > 0 {
		return helpers.NewError(http.StatusConflict, "customer_already_exists", "Customer already exists", err)
	}

	err = customers.Insert(customer)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_creation_failed", "Failed to insert the customer in the database", err)
	}

	return nil
}

func (db *mongo) FindCustomerById(id string) (*models.Customer, error) {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	customer := &models.Customer{}
	err := customers.FindId(id).One(customer)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customer_not_found", "Customer not found", err)
	}

	return customer, err
}

func (db *mongo) FindCustomer(params params.M) (*models.Customer, error) {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	customer := &models.Customer{}

	err := customers.Find(params).One(customer)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customer_not_found", "Customer not found", err)
	}

	return customer, err
}

func (db *mongo) DeleteCustomer(customerId string) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	err := customers.Remove(bson.M{"_id": customerId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_delete_failed", "Failed to delete the customer", err)
	}

	return nil
}

func (db *mongo) ActivateCustomer(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	err := customers.Update(bson.M{"$and": []bson.M{{"_id": id}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_activation_failed", "Couldn't find the customer to activate", err)
	}
	return nil
}

func (db *mongo) UpdateCustomer(customer *models.Customer, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	customers := db.C(models.CustomersCollection).With(session)

	err := customers.UpdateId(customer.Id, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "customer_update_failed", "Failed to update the customer", err)
	}

	return nil
}

func (db *mongo) GetCustomers() ([]*models.Customer, error) {
	session := db.Session.Copy()
	defer session.Close()

	customers := db.C(models.CustomersCollection).With(session)

	list := []*models.Customer{}
	err := customers.Find(params.M{}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customers_not_found", "Customers not found", err)
	}

	return list, nil
}

func (db *mongo) CustomerAssignOrganization(customer_id string, organization_id string) error {
	session := db.Session.Copy()
	defer session.Close()

	customers := db.C(models.CustomersCollection).With(session)
	err := customers.Update(bson.M{"_id": customer_id}, bson.M{"$set": bson.M{"organization_id": organization_id}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_change_failed", "Couldn't change the customer organization", err)
	}
	return nil
}

func (db *mongo) GetCustomerOrganization(customer *models.Customer) (*models.Organization, error) {
	session := db.Session.Copy()
	defer session.Close()
	organizations := db.C(models.OrganizationsCollection).With(session)

	organization := &models.Organization{}

	err := organizations.Find(bson.M{"_id": customer.OrganizationId}).One(organization)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "organization_not_found", "Organization not found", err)
	}

	return organization, err
}

func (db *mongo) CountCustomers() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	customers := db.C(models.CustomersCollection).With(session)

	nbr, err := customers.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "customers_not_found", "Customers not found", err)
	}
	return nbr, nil
}
