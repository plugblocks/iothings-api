package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateFleet(c context.Context, record *models.Fleet) error {
	return FromContext(c).CreateFleet(Current(c), record)
}

func GetAllFleets(c context.Context) ([]models.Fleet, error) {
	return FromContext(c).GetAllFleets(Current(c))
}

func GetFleetById(c context.Context, id string) (*models.Fleet, error) {
	return FromContext(c).GetFleetById(Current(c), id)
}

func UpdateFleet(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateFleet(Current(c), id, params)
}

func DeleteFleet(c context.Context, id string) error {
	return FromContext(c).DeleteFleet(Current(c), id)
}
