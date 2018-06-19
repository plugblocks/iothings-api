package store

import (
	"context"

	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

func FindUserById(c context.Context, id string) (*models.User, error) {
	return FromContext(c).FindUserById(id)
}

func FindUser(c context.Context, params params.M) (*models.User, error) {
	return FromContext(c).FindUser(params)
}

func GetUsers(c context.Context) ([]*models.User, error) {
	return FromContext(c).GetUsers()
}

func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}

func UpdateUser(c context.Context, params params.M) error {
	return FromContext(c).UpdateUser(Current(c), params)
}


func UserIsAdmin(c context.Context, userId string) (bool, error) {
	return FromContext(c).UserIsAdmin(userId)
}

func UserAttachFleet(c context.Context, userId string, fleetId string) (*models.User, error) {
	return FromContext(c).UserAttachFleet(userId, fleetId)
}

func UserDetachFleet(c context.Context, userId string, fleetId string) (*models.User, error) {
	return FromContext(c).UserDetachFleet(userId, fleetId)
}

func UserGetFleet(c context.Context, userId string, fleetId string) (*models.User, error) {
	return FromContext(c).UserGetFleet(userId, fleetId)
}

func UserGetFleets(c context.Context, userId string) ([]*models.User, error) {
	return FromContext(c).UserGetFleets(userId)
}