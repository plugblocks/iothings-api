package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) GetDeviceObservations(customer *models.Customer, deviceId string) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	//Checking that user request one of its device
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}
	err := devices.Find(bson.M{"_id": deviceId, "customer_id": customer.Id}).One(device)
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

func (db *mongo) GetDeviceLatestObservation(customer *models.Customer, deviceId string) (*models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	//Checking that customer request one of its device
	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	observations := db.C(models.ObservationsCollection).With(session)
	observation := &models.Observation{}

	err := devices.Find(bson.M{"_id": deviceId, "customer_id": customer.Id}).One(device)
	if err != nil {
		return observation, helpers.NewError(http.StatusNotFound, "customer_device_not_found", "Failed to find customer device", err)
	}

	err = observations.Find(params.M{"device_id": deviceId}).One(observation)
	if err != nil {
		return observation, helpers.NewError(http.StatusNotFound, "observation_device_not_found", "Failed to find observation for device", err)
	}

	return observation, nil
}

func (db *mongo) GetFleetObservations(user *models.User, fleetId string) ([]models.Observation, error) {
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

func (db *mongo) GetFleetLatestObservation(user *models.User, fleetId string) ([]models.Observation, error) {
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

func (db *mongo) GetAllFleetsObservations(user *models.User) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	retObservationsList := []models.Observation{}

	//Getting user fleets
	fleetsCol := db.C(models.FleetsCollection).With(session)
	fleets := []*models.Fleet{}
	err := fleetsCol.Find(bson.M{"user_id": user.Id}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_fleets_not_found", "Failed to find user fleets", err)
	}

	//Finding all observations of devices of user's fleet
	observations := db.C(models.ObservationsCollection).With(session)
	for _, fleet := range fleets {
		for _, deviceId := range fleet.DeviceIds {
			tempObservationsList := []models.Observation{}
			err = observations.Find(params.M{"device_id": deviceId}).All(tempObservationsList)
			if err != nil {
				return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
			}
			retObservationsList = append(retObservationsList, tempObservationsList...)
		}
	}
	return retObservationsList, nil
}
func (db *mongo) GetAllFleetsLatestObservation(user *models.User) ([]models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	retObservationsList := []models.Observation{}

	//Getting user fleets
	fleetsCol := db.C(models.FleetsCollection).With(session)
	fleets := []*models.Fleet{}
	err := fleetsCol.Find(bson.M{"user_id": user.Id}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_fleets_not_found", "Failed to find user fleets", err)
	}

	//Finding latest observations of devices of user's fleet
	observations := db.C(models.ObservationsCollection).With(session)
	for _, fleet := range fleets {
		for _, deviceId := range fleet.DeviceIds {
			tempObservation := &models.Observation{}
			err = observations.Find(params.M{"device_id": deviceId}).One(tempObservation)
			if err != nil {
				return nil, helpers.NewError(http.StatusNotFound, "observation_device_not_found", "Failed to find observation for device", err)
			}
			retObservationsList = append(retObservationsList, *tempObservation)
		}
	}

	return retObservationsList, nil
}
