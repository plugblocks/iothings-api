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
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleet.Id = bson.NewObjectId().Hex()
	fleet.UserId = user.Id

	err := fleetCollection.Insert(fleet)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_creation_failed", "Failed to insert the fleet in the database", err)
	}

	return nil
}

func (db *mongo) AddDeviceToFleet(user *models.User, fleetId string, deviceId string) (*models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)
	deviceCollection := db.C(models.DevicesCollection).With(session)

	fleet := &models.Fleet{}
	err := fleetCollection.Find(bson.M{"_id": fleetId, "user_id": user.Id}).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Could not find the fleet", err)
	}

	device := &models.Device{}
	err = deviceCollection.Find(bson.M{"_id": deviceId, "organization_id": user.OrganizationId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	fleet.DeviceIds = append(fleet.DeviceIds, device.Id)
	err = fleetCollection.Update(bson.M{"_id": fleet.Id, "user_id": user.Id}, fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "fleet_update_failed", "Failed to update the fleets "+err.Error(), err)
	}

	return nil, helpers.NewError(http.StatusNotFound, "fleet_device_add_error", "Could not add device to fleet", err)
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
	fleetCollection := db.C(models.FleetsCollection).With(session)

	err := fleetCollection.Update(bson.M{"_id": id, "user_id": user.Id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_update_failed", "Failed to update the fleet: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteFleet(user *models.User, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	err := fleetCollection.Remove(bson.M{"_id": id, "user_id": user.Id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_delete_failed", "Failed to delete the fleet", err)
	}

	return nil
}
