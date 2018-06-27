package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func GetDeviceObservations(c context.Context, deviceId string, typ string) ([]models.Observation, error) {
	return FromContext(c).GetDeviceObservations(CurrentCustomer(c), deviceId, typ)
}

func GetDeviceLatestObservation(c context.Context, deviceId string, typ string) (*models.Observation, error) {
	return FromContext(c).GetDeviceLatestObservation(CurrentCustomer(c), deviceId, typ)
}

func GetFleetObservations(c context.Context, fleetId string, typ string) ([]models.Observation, error) {
	return FromContext(c).GetFleetObservations(Current(c), fleetId, typ)
}

func GetFleetLatestObservation(c context.Context, fleetId string, typ string) ([]models.Observation, error) {
	return FromContext(c).GetFleetLatestObservation(Current(c), fleetId, typ)
}

func GetAllFleetsObservations(c context.Context, typ string) ([]models.Observation, error) {
	return FromContext(c).GetAllFleetsObservations(Current(c), typ)
}

func GetAllFleetsLatestObservation(c context.Context, typ string) ([]models.Observation, error) {
	return FromContext(c).GetAllFleetsLatestObservation(Current(c), typ)
}

func CreateObservation(c context.Context, record *models.Observation) error {
	return FromContext(c).CreateObservation(record)
}
