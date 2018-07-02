package mongodb

import (
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (db *mongo) CreateSigfoxMessage(message *sigfox.Message) error {
	session := db.Session.Copy()
	defer session.Close()
	sigfoxMessages := db.C(sigfox.SigfoxMessagesCollection).With(session)

	message.Id = bson.NewObjectId().Hex()

	err := sigfoxMessages.Insert(message)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "message_creation_failed", "Failed to insert the sigfox message", err)
	}
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	//err = devices.Find(params.M{"sigfox_id": message.SigfoxId}).One(device)
	err = devices.Find(bson.M{"sigfox_id": message.SigfoxId}).One(device)
	if err != nil {
		return helpers.NewError(http.StatusPartialContent, "sigfox_device_id_not_found", "Device Sigfox ID not found", err)
	} else {
		err = devices.Update(bson.M{"sigfox_id": message.SigfoxId}, bson.M{"$set": bson.M{"last_access": message.Timestamp}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device last activity", err)
		}

		err = devices.Update(bson.M{"sigfox_id": message.SigfoxId}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device status", err)
		}
	}

	return nil
}

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

func (db *mongo) GetSigfoxLocations(user *models.User) ([]sigfox.Location, error) {
	session := db.Session.Copy()
	defer session.Close()
	locationCollection := db.C(sigfox.SigfoxLocationsCollection).With(session)

	locations := []sigfox.Location{}
	err := locationCollection.Find(bson.M{"wifi": true}).All(&locations)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	return locations, nil
}
