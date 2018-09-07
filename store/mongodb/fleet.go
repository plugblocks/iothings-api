package mongodb

import (
	"fmt"
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

	return fleet, nil
}

func (db *mongo) GetFleets(user *models.User) ([]models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)
	retFleetsList := []models.Fleet{}

	//Get all users ids from this organization
	usersCollection := db.C(models.UsersCollection).With(session)
	users := []models.User{}

	err := usersCollection.Find(bson.M{"organization_id": user.OrganizationId}).All(&users)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "organization_user_retrieval_failed", "Failed to retrieve the users of the organization", err)
	}

	//Get fleets from all these users id
	for _, user := range users {
		fmt.Println("User+", user)
		tempFleet := []models.Fleet{}
		err := fleetCollection.Find(bson.M{"user_id": user.Id}).All(tempFleet)
		if err != nil {
			return nil, helpers.NewError(http.StatusInternalServerError, "query_fleets_failed", "Failed to get the user fleets: "+err.Error(), err)
		}
		fmt.Println(tempFleet)
		retFleetsList = append(retFleetsList, tempFleet...)
	}

	return retFleetsList, nil
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

func (db *mongo) GetDevicesFromFleet(user *models.User, id string) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)
	deviceCollection := db.C(models.DevicesCollection).With(session)

	fleet := &models.Fleet{}
	err := fleetCollection.Find(bson.M{"_id": id, "user_id": user.Id}).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Could not find the fleet", err)
	}

	retDevicesList := []*models.Device{}
	for _, deviceId := range fleet.DeviceIds {
		fmt.Println("Device+", deviceId)
		tempDevice := &models.Device{}
		err := deviceCollection.Find(bson.M{"_id": deviceId}).One(tempDevice)
		if err != nil {
			fmt.Println(err)
			return nil, helpers.NewError(http.StatusNotFound, "fleet_device_not_found", "Failed to find device from fleet", err)
		}
		fmt.Println(tempDevice)
		retDevicesList = append(retDevicesList, tempDevice)
	}

	return retDevicesList, nil
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

func (db *mongo) CountFleets() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	fleets := db.C(models.FleetsCollection).With(session)

	nbr, err := fleets.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "fleets_not_found", "Fleets not found", err)
	}
	return nbr, nil
}
