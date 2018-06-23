package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateOrganization(organization *models.Organization) error {
	session := db.Session.Copy()
	defer session.Close()
	organizations := db.C(models.OrganizationsCollection).With(session)

	organization.Id = bson.NewObjectId().Hex()

	err := organizations.Insert(organization)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_creation_failed", "Failed to insert the organization in the database", err)
	}

	return nil
}

func (db *mongo) GetAllOrganizations() ([]models.Organization, error) {
	session := db.Session.Copy()
	defer session.Close()
	organizationCollection := db.C(models.OrganizationsCollection).With(session)

	organizations := []models.Organization{}
	err := organizationCollection.Find(bson.M{}).All(&organizations)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_organizations_failed", "Failed to get the organizations: "+err.Error(), err)
	}

	return organizations, nil
}

func (db *mongo) GetOrganizationById(id string) (*models.Organization, error) {
	session := db.Session.Copy()
	defer session.Close()
	organizationCollection := db.C(models.OrganizationsCollection).With(session)

	organization := &models.Organization{}
	err := organizationCollection.Find(bson.M{"_id": id}).One(organization)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "organization_not_found", "Could not find the organization", err)
	}

	return organization, nil
}

func (db *mongo) UpdateOrganization(id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	organizations := db.C(models.OrganizationsCollection).With(session)

	err := organizations.Update(bson.M{"_id": id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_update_failed", "Failed to update the organizations: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteOrganization(id string) error {
	session := db.Session.Copy()
	defer session.Close()
	organizations := db.C(models.OrganizationsCollection).With(session)

	err := organizations.Remove(bson.M{"_id": id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_delete_failed", "Failed to delete the organization", err)
	}

	return nil
}
