package store

import (
	"context"

	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateCustomer(c context.Context, record *models.Customer) error {
	return FromContext(c).CreateCustomer(record)
}

func DeleteCustomer(c context.Context, customerId string) error {
	return FromContext(c).DeleteCustomer(customerId)
}

func FindCustomerById(c context.Context, id string) (*models.Customer, error) {
	return FromContext(c).FindCustomerById(id)
}

func FindCustomer(c context.Context, params params.M) (*models.Customer, error) {
	return FromContext(c).FindCustomer(params)
}

func GetCustomers(c context.Context) ([]*models.Customer, error) {
	return FromContext(c).GetCustomers()
}

func ActivateCustomer(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateCustomer(activationKey, id)
}

func UpdateCustomer(c context.Context, customer *models.Customer, params params.M) error {
	return FromContext(c).UpdateCustomer(customer, params)
}

func CustomerAssignOrganization(c context.Context, customerId string, organizationId string) error {
	return FromContext(c).CustomerAssignOrganization(customerId, organizationId)
}

func GetCustomerOrganization(c context.Context, customer *models.Customer) (*models.Organization, error) {
	return FromContext(c).GetCustomerOrganization(customer)
}

func CountCustomers(c context.Context) (int, error) {
	return FromContext(c).CountCustomers()
}
