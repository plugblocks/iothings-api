package server

import (
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store/mongodb"
)

func (a *API) SetupSeeds() error {
	store := mongodb.New(a.Database)

	organization := &models.Organization{
		Name:   "IoThings",
		Active: true,
	}
	store.CreateOrganization(organization)

	user := &models.User{
		Firstname:      "Adrien",
		Lastname:       "Chapelet",
		Password:       "adchapwd",
		Email:          "admin@iothings.fr",
		OrganizationId: organization.Id,
		Admin:          true,
	}

	store.CreateUser(user)
	store.ActivateUser(user.ActivationKey, user.Id)
	return nil
}
