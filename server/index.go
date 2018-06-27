package server

import (
	"gitlab.com/plugblocks/iothings-api/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (a *API) SetupIndexes() error {
	database := a.Database

	// Creates a list of indexes to ensure
	collectionIndexes := make(map[*mgo.Collection][]mgo.Index)

	// User indexes
	users := database.C(models.UsersCollection)
	collectionIndexes[users] = []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	}

	// Devices indexes & validators
	devices := database.C(models.DevicesCollection)
	CreateValidator(devices, bson.M{"organization_id": bson.M{"$exists": true}})
	collectionIndexes[devices] = []mgo.Index{
		{
			Key: []string{"organization_id"},
		},
		{
			Key:    []string{"metadata.sigfox_id"},
			Unique: true,
		},
		{
			Key:    []string{"metadata.ble_mac"},
			Unique: true,
		},
		{
			Key:    []string{"metadata.wifi_mac"},
			Unique: true,
		},
	}

	for collection, indexes := range collectionIndexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateValidator(collection *mgo.Collection, validator bson.M) {
	info := &mgo.CollectionInfo{
		Validator:       validator,
		ValidationLevel: "strict",
	}
	collection.Create(info)
}
