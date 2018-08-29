package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateObservation(c context.Context, record *models.Observation) error {
	return FromContext(c).CreateObservation(record)
}

func GetDeviceObservations(c context.Context, deviceId string, typ string, ord string, lim int) ([]*models.Observation, error) {
	return FromContext(c).GetDeviceObservations(deviceId, typ, ord, lim)
}

func GetFleetObservations(c context.Context, fleetId string, typ string, ord string, lim int) ([]*models.Observation, error) {
	return FromContext(c).GetFleetObservations(Current(c), fleetId, typ, ord, lim)
}

func CountObservations(c context.Context) (int, error) {
	return FromContext(c).CountObservations()
}
