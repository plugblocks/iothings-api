package mongodb

import (
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
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
	/*err := message.BeforeCreate()
	if err != nil {
		return err
	}*/

	err := sigfoxMessages.Insert(message)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "message_creation_failed", "Failed to insert the sigfox message", err)
	}
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err = devices.Find(params.M{"sigfoxId": message.SigfoxId}).One(device)
	if err != nil {
		return helpers.NewError(http.StatusPartialContent, "sigfox_device_id_not_found", "Device Sigfox ID not found", err)
	} else {
		err = devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"lastAcc": message.Timestamp}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device last activity", err)
		}

		err = devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device status", err)
		}
	}

	return nil
}
