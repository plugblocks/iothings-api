package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateGeolocation(c context.Context, location *models.Geolocation) error {
	return FromContext(c).CreateGeolocation(location)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetDeviceGeoJSON(c context.Context, id string) (*models.GeoJSON, error) {
	return FromContext(c).GetDeviceGeoJSON( /*Current(c), */ id)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetFleetGeoJSON(c context.Context, id string) (*models.GeoJSON, error) {
	return FromContext(c).GetFleetGeoJSON( /*Current(c), */ id)
}

//TODO: DANGER: Protect by auth device GeoJSON
func GetFleetsGeoJSON(c context.Context, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
	return FromContext(c).GetFleetsGeoJSON(source, limit, startTime, endTime)
}

func GetUserFleetsGeoJSON(c context.Context) (*models.GeoJSON, error) {
	return FromContext(c).GetUserFleetsGeoJSON(Current(c))
}

func CountGeolocations(c context.Context) (int, error) {
	return FromContext(c).CountGeolocations()
}
