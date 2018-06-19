package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateGroup(c context.Context, record *models.Group) error {
	return FromContext(c).CreateGroup(record)
}