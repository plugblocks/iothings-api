package mongodb

import (
	"net/http"

	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateMessage(message *models.SigfoxMessage) error {
	session := db.Session.Copy()
	defer session.Close()
	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	message.BeforeCreate()
	err := sigfoxMessages.Insert(message)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "message_creation_failed", "Failed to insert the sigfox message")
	}

	devices := db.C(models.DevicesCollection).With(session)

	devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"lastAcc": message.Timestamp}})

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device")
	}

	return nil
}

func (db *mongo) CreateLocation(location *models.Location) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(models.LocationsCollection).With(session)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "location_creation_failed", "Failed to insert the location")
	}

	return nil
}
