package mongodb

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func (db *mongo) CreateDevice(device *models.Device) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	device.BeforeCreate()
	err := devices.Insert(device)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_creation_failed", "Failed to create the device", err)
	}

	return nil
}

func (db *mongo) GetDevices(customerId string) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.Device{}
	err := devices.Find(params.M{"customer_id": customerId}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "devices_not_found", "Devices not found", err)
	}

	return list, nil
}

func (db *mongo) UpdateDevice(id string, m params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	change := mgo.Change{
		Update:    m,
		Upsert:    false,
		Remove:    false,
		ReturnNew: false,
	}
	_, err := devices.Find(bson.M{"_id": id}).Apply(change, nil)

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device", err)
	}

	return nil
}

func (db *mongo) DeleteDevice(id string) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	err := devices.RemoveId(id)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_delete_failed", "Failed to delete the device", err)
	}

	return nil
}

func (db *mongo) GetDevice(id string) (*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	return device, nil
}
