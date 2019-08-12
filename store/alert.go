package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func GetAlerts(c context.Context) ([]*models.Alert, error) {
	return FromContext(c).GetAlerts(Current(c))
}

func CreateAlert(c context.Context, alert *models.Alert) error {
	return FromContext(c).CreateAlert(Current(c), alert)
}

func GetAlert(c context.Context, id string) (*models.Alert, error) {
	return FromContext(c).GetAlert(Current(c), id)
}

func GetFleetAlerts(c context.Context, fleetId string) ([]*models.Alert, error) {
	return FromContext(c).GetFleetAlerts(Current(c), fleetId)
}

func GetDeviceAlerts(c context.Context, deviceId string) ([]*models.Alert, error) {
	return FromContext(c).GetDeviceAlerts(Current(c), deviceId)
}

func UpdateAlert(c context.Context, id string, params params.M) error {
	return FromContext(c).UpdateAlert(Current(c), id, params)
}

func DeleteAlert(c context.Context, id string) error {
	return FromContext(c).DeleteAlert(Current(c), id)
}

func CountAlerts(c context.Context) (int, error) {
	return FromContext(c).CountAlerts()
}
