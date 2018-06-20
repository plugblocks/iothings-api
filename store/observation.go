package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func GetDeviceObservations(c context.Context, deviceId string) ([]models.Observation, error) {
	return FromContext(c).GetDeviceObservations(CurrentCustomer(c), deviceId)
}

func GetDeviceLatestObservation(c context.Context, deviceId string) (*models.Observation, error) {
	return FromContext(c).GetDeviceLatestObservation(CurrentCustomer(c), deviceId)
}

func GetFleetObservations(c context.Context, fleetId string) ([]models.Observation, error) {
	return FromContext(c).GetFleetObservations(Current(c), fleetId)
}

func GetFleetLatestObservation(c context.Context, fleetId string) ([]models.Observation, error) {
	return FromContext(c).GetFleetLatestObservation(Current(c), fleetId)
}

func GetAllFleetsObservations(c context.Context) ([]models.Observation, error) {
	return FromContext(c).GetAllFleetsObservations(Current(c))
}

func GetAllFleetsLatestObservation(c context.Context) ([]models.Observation, error) {
	return FromContext(c).GetAllFleetsLatestObservation(Current(c))
}
