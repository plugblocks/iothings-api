package mongodb

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

//TODO: CHECK USER/CUSTOMER RIGHTS
func (db *mongo) GetDeviceObservations(deviceId string, resolver string, order string, lim int) ([]*models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	//Checking that user request one of its device
	/*devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}
	//TODO: Refactor code to use only one MongoDB request, restore customer ownership security
	err := devices.Find(bson.M{"_id": deviceId, "customer_id": customerId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "customer_device_not_found", "Failed to find customer device", err)
	}*/

	observations := db.C(models.ObservationsCollection).With(session)
	list := []*models.Observation{}
	if resolver == "" {
		err := observations.Find(bson.M{"device_id": deviceId}).Sort(order).Limit(lim).All(&list)
		if err != nil {
			fmt.Println("device get obs id:", deviceId, " err:", err)
			return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
		}
	} else {
		err := observations.Find(bson.M{"device_id": deviceId, "resolver": resolver}).Sort(order).Limit(lim).All(&list)
		if err != nil {
			fmt.Println("device get obs id:", deviceId, " err:", err)
			return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
		}
	}

	return list, nil
}

func (db *mongo) GetFleetObservations(user *models.User, fleetId string, resolver string, order string, lim int) ([]*models.Observation, error) {
	session := db.Session.Copy()
	defer session.Close()

	retObservationsList := []*models.Observation{}

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
		//fmt.Println("Device+", deviceId)
		tempObservationsList := []*models.Observation{}
		// TODO: Sort by timestamp decreasing
		err = observations.Find(params.M{"device_id": deviceId}).Sort(order).Limit(lim).All(&tempObservationsList)
		if err != nil {
			fmt.Println(err)
			return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
		}
		//fmt.Println(tempObservationsList)
		retObservationsList = append(retObservationsList, tempObservationsList...)
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

	record.Id = bson.NewObjectId().Hex()

	err = observations.Insert(record)
	if err != nil {
		fmt.Println(err)
		//return helpers.NewError(http.StatusInternalServerError, "observation_creation_failed", "Failed to create the observation", err)
	}

	return nil
}

func (db *mongo) CountObservations() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	observations := db.C(models.ObservationsCollection).With(session)

	nbr, err := observations.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "observations_not_found", "Observations not found", err)
	}
	return nbr, nil
}

func (db *mongo) DeleteObservation(id string) error {
	session := db.Session.Copy()
	defer session.Close()

	observations := db.C(models.ObservationsCollection).With(session)

	err := observations.Remove(bson.M{"_id": id})

	if err != nil {
		return helpers.NewError(http.StatusNotFound, "observation_not_found", "Observation not found", err)
	}
	return nil
}
