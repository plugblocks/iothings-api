package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateFleet(user *models.User, fleet *models.Fleet) error {
	session := db.Session.Copy()
	defer session.Close()
	fleets := db.C(models.FleetsCollection).With(session)

	fleet.Id = bson.NewObjectId().Hex()
	fleet.UserId = user.Id

	err := fleets.Insert(fleet)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_creation_failed", "Failed to insert the fleet in the database", err)
	}

	return nil
}

func (db *mongo) GetAllFleets(user *models.User) ([]models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleets := []models.Fleet{}
	err := fleetCollection.Find(bson.M{"user_id": user.Id}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_fleets_failed", "Failed to get the fleets: "+err.Error(), err)
	}

	return fleets, nil
}

func (db *mongo) GetFleetById(user *models.User, id string) (*models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleet := &models.Fleet{}
	err := fleetCollection.Find(bson.M{"_id": id, "user_id": user.Id}).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Could not find the fleet", err)
	}

	return fleet, nil
}

func (db *mongo) UpdateFleet(user *models.User, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	fleets := db.C(models.FleetsCollection).With(session)

	err := fleets.Update(bson.M{"_id": id, "user_id": user.Id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_update_failed", "Failed to update the fleets: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteFleet(user *models.User, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	fleets := db.C(models.FleetsCollection).With(session)

	err := fleets.Remove(bson.M{"_id": id, "user_id": user.Id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_delete_failed", "Failed to delete the fleet", err)
	}

	return nil
}
