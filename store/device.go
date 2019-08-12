package store

import (
	"context"
	"fmt"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"

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

func GetAvailableDevices(c context.Context) ([]*models.Device, error) {
	fmt.Println("Context: ", c)
	return FromContext(c).GetAvailableDevices(Current(c).OrganizationId)
}

func UpdateDevice(c context.Context, id string, m params.M) error {
	return FromContext(c).UpdateDevice(Current(c).Id, id, m)
}

func UpdateDeviceActivity(c context.Context, deviceId string, actDiff int) error {
	return FromContext(c).UpdateDeviceActivity(deviceId, actDiff)
}

func DeleteDevice(c context.Context, id string) error {
	return FromContext(c).DeleteDevice(Current(c), id)
}

func GetDevice(c context.Context, id string) (*models.Device, error) {
	return FromContext(c).GetDevice(Current(c).OrganizationId, id)
}

func GetDeviceFromSigfoxId(c context.Context, sigfoxId string) (*models.Device, error) {
	return FromContext(c).GetDeviceFromSigfoxId(sigfoxId)
}

func GetDeviceMessages(c context.Context, sigfoxId string) ([]*sigfox.Message, error) {
	return FromContext(c).GetDeviceMessages(sigfoxId)
}

func CountDevices(c context.Context) (int, error) {
	return FromContext(c).CountDevices()
}

func GetDeviceGeolocations(c context.Context, id string, source string, limit int, startTime int, endTime int) ([]*models.Geolocation, error) {
	return FromContext(c).GetDeviceGeolocations(id, source, limit, startTime, endTime)
}

func GetDevicePreciseGeolocations(c context.Context, id string, limit int, startTime int, endTime int) ([]*models.Geolocation, error) {
	return FromContext(c).GetDevicePreciseGeolocations(id, limit, startTime, endTime)
}

func DeleteDeviceObservations(c context.Context, deviceId string) error {
	return FromContext(c).DeleteDeviceObservations(deviceId)
}
func DeleteDeviceGeolocations(c context.Context, deviceId string) error {
	return FromContext(c).DeleteDeviceGeolocations(deviceId)
}
