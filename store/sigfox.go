package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/models"
)

func CreateSigfoxMessage(c context.Context, message *sigfox.Message) error {
	return FromContext(c).CreateSigfoxMessage(message)
}

func CreateSigfoxLocation(c context.Context, location *sigfox.Location) error {
	return FromContext(c).CreateSigfoxLocation(location)
}

func GetSigfoxLocations(c context.Context) ([]sigfox.Location, error) {
	return FromContext(c).GetSigfoxLocations()
}
func GetGeoJSON(c context.Context) (*models.GeoJSON, error) {
	return FromContext(c).GetGeoJSON()
}