package mongodb

import (
	"net/http"

	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func (db *mongo) CreateDevice(organizationId string/*user *models.User*/, device *models.Device) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	//fmt.Println("Mongo User: " + user.Email + "OrgaId" + user.OrganizationId)
	device.BeforeCreate()
	device.OrganizationId = organizationId

	if device.SigfoxId != "" {
		count, _ := devices.Find(bson.M{"sigfox_id": device.SigfoxId}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "Sigfox ID already registered", errors.New("Sigfox ID already registered"))
		}

		/*if count == 0 {
			fmt.Println("Sigfox Messages for this New device: ", count)
			sigfoxMessages := db.C(sigfox.SigfoxMessagesCollection).With(session)
			nbr, err := sigfoxMessages.Find(params.M{}).Count()
			if err != nil {
				return helpers.NewError(http.StatusNotFound, "sigfox_messages_not_found", "Sigfox Messages not found", err)
			}
			fmt.Println("Sigfox Messages", nbr)

		}*/
	}

	if device.BleMac != "" {
		count, _ := devices.Find(bson.M{"ble_mac": device.BleMac}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "BLE MAC already registered", errors.New("BLE MAC already registered"))
		}
	}

	if device.WifiMac != "" {
		count, _ := devices.Find(bson.M{"wifi_mac": device.WifiMac}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "WiFi MAC already registered", errors.New("WiFi MAC already registered"))
		}
	}

	err := devices.Insert(&device)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_creation_failed", "Failed to create the device", err)
	}

	return nil
}

func (db *mongo) GetDevices(user *models.User) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.Device{}
	err := devices.Find(params.M{"organization_id": user.OrganizationId}).All(&list)
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

func (db *mongo) CountDevices() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	nbr, err := devices.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "devices_not_found", "Devices not found", err)
	}
	return nbr, nil
}

func (db *mongo) DeleteDeviceObservations(deviceId string) error {
	session := db.Session.Copy()
	defer session.Close()

	observations := db.C(models.ObservationsCollection).With(session)

	err := observations.Remove(bson.M{"device_id": deviceId})

	if err != nil {
		return helpers.NewError(http.StatusNotFound, "observations_not_found", "Observations not found", err)
	}
	return nil
}

func (db *mongo) DeleteDeviceGeolocations(deviceId string) error {
	session := db.Session.Copy()
	defer session.Close()

	geolocations := db.C(models.GeolocationsCollection).With(session)

	err := geolocations.Remove(bson.M{"device_id": deviceId})

	if err != nil {
		return helpers.NewError(http.StatusNotFound, "geolocations_not_found", "Geolocations not found", err)
	}
	return nil
}
