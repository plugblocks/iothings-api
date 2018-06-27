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

func AssignOrganization(c context.Context, userId string, organizationId string) error {
	return FromContext(c).AssignOrganization(userId, organizationId)
}

func GetUserOrganization(c context.Context, user *models.User) (*models.Organization, error) {
	return FromContext(c).GetUserOrganization(user)
}
