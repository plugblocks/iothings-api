package store

import (
	"context"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
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

func CountSigfoxMessages(c context.Context) (int, error) {
	return FromContext(c).CountSigfoxMessages()
}
