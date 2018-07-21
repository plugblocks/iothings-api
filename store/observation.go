package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateObservation(c context.Context, record *models.Observation) error {
	return FromContext(c).CreateObservation(record)
}

func GetDeviceObservations(c context.Context, customerId string, deviceId string, typ string, lim int) ([]*models.Observation, error) {
	return FromContext(c).GetDeviceObservations(customerId, deviceId, typ, lim)
}

func GetDeviceLatestObservation(c context.Context, customerId string, deviceId string, typ string) (*models.Observation, error) {
	return FromContext(c).GetDeviceLatestObservation(customerId, deviceId, typ)
}

func GetFleetObservations(c context.Context, fleetId string, typ string, lim int) ([]*models.Observation, error) {
	return FromContext(c).GetFleetObservations(Current(c), fleetId, typ, lim)
}

func GetFleetLatestObservation(c context.Context, fleetId string, typ string) ([]models.Observation, error) {
	return FromContext(c).GetFleetLatestObservation(Current(c), fleetId, typ)
}
