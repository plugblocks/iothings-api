package store

import (
	"context"

	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateDevice(c context.Context, record *models.Device) error {
	return FromContext(c).CreateDevice(Current(c), record)
}

func GetDevices(c context.Context, customerId string) ([]*models.Device, error) {
	return FromContext(c).GetDevices(Current(c), customerId)
}

func UpdateDevice(c context.Context, id string, m params.M) error {
	return FromContext(c).UpdateDevice(Current(c), id, m)
}

func DeleteDevice(c context.Context, id string) error {
	return FromContext(c).DeleteDevice(Current(c), id)
}

func GetDevice(c context.Context, id string) (*models.Device, error) {
	return FromContext(c).GetDevice(Current(c), id)
}

func GetDeviceIdFromSigfoxId(c context.Context, sigfoxId string) (*models.Device, error) {
	return FromContext(c).GetDevice(Current(c), sigfoxId)
}
