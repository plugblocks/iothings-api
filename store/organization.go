package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
)

func CreateOrganization(c context.Context, record *models.Organization) error {
	return FromContext(c).CreateOrganization(record)
}

func GetAllOrganizations(c context.Context) ([]models.Organization, error) {
	return FromContext(c).GetAllOrganizations()
}

func GetOrganizationById(c context.Context, id string) (*models.Organization, error) {
	return FromContext(c).GetOrganizationById(id)
}

func UpdateOrganization(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateOrganization(id, params)
}

func DeleteOrganization(c context.Context, id string) error {
	return FromContext(c).DeleteOrganization(id)
}
