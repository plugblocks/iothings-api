package mongodb

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func (db *mongo) CreateDevice(user *models.User, device *models.Device) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	device.BeforeCreate(user)

	err := devices.Insert(device)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_creation_failed", "Failed to create the device", err)
	}

	return nil
}

func (db *mongo) GetDevices(user *models.User, customerId string) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.Device{}
	err := devices.Find(params.M{"customer_id": customerId, "organization_id": user.OrganizationId}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "devices_not_found", "Devices not found", err)
	}

	return list, nil
}

func (db *mongo) UpdateDevice(user *models.User, id string, m params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	change := mgo.Change{
		Update:    m,
		Upsert:    false,
		Remove:    false,
		ReturnNew: false,
	}
	_, err := devices.Find(bson.M{"_id": id, "organization_id": user.OrganizationId}).Apply(change, nil)

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device", err)
	}

	return nil
}

func (db *mongo) DeleteDevice(user *models.User, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	err := devices.Remove(bson.M{"_id": id, "organization_id": user.OrganizationId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_delete_failed", "Failed to delete the device", err)
	}

	return nil
}

func (db *mongo) GetDevice(user *models.User, id string) (*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"_id": id, "organization_id": user.OrganizationId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	return device, nil
}

func (db *mongo) GetDeviceFromSigfoxId(sigfoxId string) (*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"sigfox_id": sigfoxId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	return device, nil
}
