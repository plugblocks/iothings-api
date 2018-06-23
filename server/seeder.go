package server

import (
	"gitlab.com/plugblocks/iothings-api/store/mongodb"
	"gitlab.com/plugblocks/iothings-api/models"
)

func (a *API) SetupSeeds() error {
	store := mongodb.New(a.Database)
	user := &models.User{
		Firstname: "admin",
		Lastname: "admin",
		Password: "admin",
		Email: "admin@iothings.fr",
		Admin: true,
	}

	store.CreateUser(user)
	store.ActivateUser(user.ActivationKey, user.Id)
	return nil
}