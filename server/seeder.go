package server

import (
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/store/mongodb"
)

func (a *API) SetupSeeds() error {
	store := mongodb.New(a.Database)

	//Mails: 0.10$/1000         Texts: 0.05-0.10$/1       WiFi: 5$/1000
	organization := &models.Organization{
		Name:   "PlugBlocks",
		Active: true,
		Admin:  true,
	}
	store.CreateOrganization(organization)

	user := &models.User{
		Firstname:      "Adrien",
		Lastname:       "Chapelet",
		Password:       "adchapwd",
		Email:          "admin@plugblocks.fr",
		OrganizationId: organization.Id,
		Admin:          true,
	}

	store.CreateUser(user)
	store.ActivateUser(user.ActivationKey, user.Id)
	return nil
}
