package mongodb

import (
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"net/http"
)

func (db *mongo) CreateSigfoxLocation(location *sigfox.Location) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(sigfox.SigfoxLocationsCollection).With(session)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "location_creation_failed", "Failed to insert the location", err)
	}

	return nil
}
