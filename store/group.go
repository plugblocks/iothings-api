package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateGroup(c context.Context, record *models.Group) error {
	return FromContext(c).CreateGroup(Current(c), record)
}

func GetAllGroups(c context.Context) ([]models.Group, error) {
	return FromContext(c).GetAllGroups(Current(c))
}

func GetGroupById(c context.Context, id string) (*models.Group, error) {
	return FromContext(c).GetGroupById(Current(c), id)
}

func UpdateGroup(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateGroup(Current(c), id, params)
}

func DeleteGroup(c context.Context, id string) error {
	return FromContext(c).DeleteGroup(Current(c), id)
}

func CountGroups(c context.Context) (int, error) {
	return FromContext(c).CountGroups()
}
