package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) GetDeviceObservations(customerId string, deviceId string, typ string) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	//Checking that user request one of its device
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}
	//TODO: Refactor code to use only one MongoDB request
	err := devices.Find(bson.M{"_id": deviceId, "customer_id": customerId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customer_device_not_found", "Failed to find customer device", err)
	}

	observations := db.C(models.ObservationsCollection).With(session)
	list := []models.Observation{}
	err = observations.Find(params.M{"device_id": deviceId}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
	}

	return list, nil
}

func (db *mongo) GetDeviceLatestObservation(customerId string, deviceId string, typ string) (*models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	//Checking that customer request one of its device
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	observations := db.C(models.ObservationsCollection).With(session)
	observation := &models.Observation{}

	//TODO: Refactor code to use only one MongoDB request
	err := devices.Find(bson.M{"_id": deviceId, "customer_id": customerId}).One(device)
	if err != nil {
		return observation, helpers.NewError(http.StatusNotFound, "customer_device_not_found", "Failed to find customer device", err)
	}

	err = observations.Find(params.M{"device_id": deviceId}).One(observation)
	if err != nil {
		return observation, helpers.NewError(http.StatusNotFound, "observation_device_not_found", "Failed to find observation for device", err)
	}

	return observation, nil
}

func (db *mongo) GetFleetObservations(user *models.User, fleetId string, typ string) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	retObservationsList := []models.Observation{}

	//Checking that user request one of its fleets
	fleets := db.C(models.FleetsCollection).With(session)
	fleet := &models.Fleet{}
	err := fleets.Find(bson.M{"_id": fleetId, "user_id": user.Id}).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_fleet_not_found", "Failed to find user fleet", err)
	}

	//Finding all observations of devices of user's fleet
	observations := db.C(models.ObservationsCollection).With(session)
	for _, deviceId := range fleet.DeviceIds {
		tempObservationsList := []models.Observation{}
		err = observations.Find(params.M{"device_id": deviceId}).All(tempObservationsList)
		if err != nil {
			return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
		}
		retObservationsList = append(retObservationsList, tempObservationsList...)
	}

	return retObservationsList, nil
}

func (db *mongo) GetFleetLatestObservation(user *models.User, fleetId string, typ string) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	retObservationsList := []models.Observation{}

	//Checking that user request one of its fleets
	//TODO: Refactor code to use only one MongoDB request
	fleets := db.C(models.FleetsCollection).With(session)
	fleet := &models.Fleet{}
	err := fleets.Find(bson.M{"_id": fleetId, "user_id": user.Id}).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_fleet_not_found", "Failed to find user fleet", err)
	}

	//Finding latest observations of devices of user's fleet
	observations := db.C(models.ObservationsCollection).With(session)
	for _, deviceId := range fleet.DeviceIds {
		tempObservation := &models.Observation{}
		err = observations.Find(params.M{"device_id": deviceId}).One(tempObservation)
		if err != nil {
			return nil, helpers.NewError(http.StatusNotFound, "observation_device_not_found", "Failed to find observation for device", err)
		}
		retObservationsList = append(retObservationsList, *tempObservation)
	}

	return retObservationsList, nil
}

func (db *mongo) CreateObservation(record *models.Observation) error {
	session := db.Session.Copy()
	defer session.Close()
	observations := db.C(models.ObservationsCollection).With(session)
	devices := db.C(models.DevicesCollection).With(session)

	device := &models.Device{}

	err := devices.Find(bson.M{"_id": record.DeviceId}).One(device)
	if err != nil {
		return helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	record.BeforeCreate(device)

	err = observations.Insert(record)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "observation_creation_failed", "Failed to create the observation", err)
	}

	return nil
}
