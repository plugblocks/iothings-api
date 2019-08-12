package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateCustomer(c context.Context, record *models.Customer) error {
	return FromContext(c).CreateCustomer(Current(c).OrganizationId, record)
}

func GetAllCustomers(c context.Context) ([]models.Customer, error) {
	return FromContext(c).GetAllCustomers(Current(c).OrganizationId)
}

func GetCustomerById(c context.Context, id string) (*models.Customer, error) {
	return FromContext(c).GetCustomerById(Current(c).OrganizationId, id)
}

func UpdateCustomer(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateCustomer(Current(c).OrganizationId, id, params)
}

func DeleteCustomer(c context.Context, id string) error {
	return FromContext(c).DeleteCustomer(Current(c).OrganizationId, id)
}
