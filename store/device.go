package store

import (
	"context"
	"fmt"

	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateDevice(c context.Context, record *models.Device) error {
	if c.Value(CurrentKey) == nil {
		fmt.Println("Nil context")
		return FromContext(nil).CreateDevice("", record)
	}

	fmt.Println("Context: ", c)
	return FromContext(c).CreateDevice(Current(c).OrganizationId, record)
}

func GetDevices(c context.Context) ([]*models.Device, error) {
	fmt.Println("Context: ", c)
	return FromContext(c).GetDevices(Current(c))
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

func GetDeviceFromSigfoxId(c context.Context, sigfoxId string) (*models.Device, error) {
	return FromContext(c).GetDeviceFromSigfoxId(sigfoxId)
}

func CountDevices(c context.Context) (int, error) {
	return FromContext(c).CountDevices()
}

func DeleteDeviceObservations(c context.Context, deviceId string) error {
	return FromContext(c).DeleteDeviceObservations(deviceId)
}
func DeleteDeviceGeolocations(c context.Context, deviceId string) error {
	return FromContext(c).DeleteDeviceGeolocations(deviceId)
}
