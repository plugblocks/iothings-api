package store

import (
	"context"

	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateDevice(c context.Context, record *models.Device) error {
	return FromContext(c).CreateDevice(record)
}

func GetDevices(c context.Context, customerId string) ([]*models.Device, error) {
	return FromContext(c).GetDevices(customerId)
}

func UpdateDevice(c context.Context, id string, m params.M) error {
	return FromContext(c).UpdateDevice(id, m)
}

func DeleteDevice(c context.Context, id string) error {
	return FromContext(c).DeleteDevice(id)
}

func GetDevice(c context.Context, id string) (*models.Device, error) {
	return FromContext(c).GetDevice(id)
}
