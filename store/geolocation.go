package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateGeolocation(c context.Context, location *models.Geolocation) error {
	return FromContext(c).CreateGeolocation(location)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetFleetsGeoJSON(c context.Context, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
	return FromContext(c).GetFleetsGeoJSON(source, limit, startTime, endTime)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetFleetGeoJSON(c context.Context, id string, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
	return FromContext(c).GetFleetGeoJSON( /*Current(c), */ id, source, limit, startTime, endTime)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetDeviceGeoJSON(c context.Context, id string, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
	return FromContext(c).GetDeviceGeoJSON( /*Current(c), */ id, source, limit, startTime, endTime)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetDeviceGeolocation(c context.Context, deviceId string, source string) (*models.Geolocation, error) {
	return FromContext(c).GetDeviceGeolocation(Current(c), deviceId, source)
}

func GetUserFleetsGeoJSON(c context.Context) (*models.GeoJSON, error) {
	return FromContext(c).GetUserFleetsGeoJSON(Current(c))
}

func GetOrderGeolocations(c context.Context, orderId string) ([]*models.Geolocation, error) {
	return FromContext(c).GetOrderGeolocations(Current(c).OrganizationId, orderId)
}

func CountGeolocations(c context.Context) (int, error) {
	return FromContext(c).CountGeolocations()
}
