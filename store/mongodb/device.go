package mongodb

import (
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"net/http"
	"sort"

	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
)

func (db *mongo) CreateDevice(organizationId string /*user *models.User*/, device *models.Device) error {
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

func (db *mongo) GetAvailableDevices(organizationId string) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.Device{}
	err := devices.Find(params.M{"organization_id": organizationId, "available": true}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "devices_not_found", "Devices not found", err)
	}

	return list, nil
}

func (db *mongo) UpdateDevice(organizationId, id string, m params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	change := mgo.Change{
		Update:    m,
		Upsert:    false,
		Remove:    false,
		ReturnNew: false,
	}
	_, err := devices.Find(bson.M{"_id": id, "organization_id": organizationId}).Apply(change, nil)

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device", err)
	}

	return nil
}

func (db *mongo) UpdateDeviceActivity(id string, act int) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	err := devices.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"activity": act}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device status", err)
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

func (db *mongo) GetDevice(organizationId string, id string) (*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"_id": id, "organization_id": organizationId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	return device, nil
}

func (db *mongo) GetDeviceGeolocations(deviceId string, source string, limit int, startTime int, endTime int) ([]*models.Geolocation, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"_id": deviceId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	deviceGeolocations := []*models.Geolocation{}

	if source == "" {
		err = geolocationCollection.Find(bson.M{"device_id": deviceId, "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&deviceGeolocations)
	} else {
		err = geolocationCollection.Find(bson.M{"device_id": deviceId, "source": source, "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&deviceGeolocations)
	}

	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	return deviceGeolocations, nil
}

func (db *mongo) GetDevicePreciseGeolocations(deviceId string, limit int, startTime int, endTime int) ([]*models.Geolocation, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"_id": deviceId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	if limit == 0 {
		limit = 50
	}

	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	deviceGeolocations := []*models.Geolocation{}
	wifiGeolocations := []*models.Geolocation{}
	gpsGeolocations := []*models.Geolocation{}

	err = geolocationCollection.Find(bson.M{"device_id": deviceId, "source": "gps", "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&gpsGeolocations)
	err = geolocationCollection.Find(bson.M{"device_id": deviceId, "source": "wifi", "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&wifiGeolocations)

	deviceGeolocations = append(deviceGeolocations, gpsGeolocations...)
	deviceGeolocations = append(deviceGeolocations, wifiGeolocations...)

	sort.Slice(deviceGeolocations, func(i, j int) bool { return deviceGeolocations[i].Timestamp < deviceGeolocations[i].Timestamp })

	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	return deviceGeolocations, nil
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

func (db *mongo) GetDeviceMessages(sigfoxId string) ([]*sigfox.Message, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	messages := db.C(sigfox.SigfoxMessagesCollection).With(session)
	device := &models.Device{}
	deviceMessages := []*sigfox.Message{}

	err := devices.Find(bson.M{"sigfox_id": sigfoxId}).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	err = messages.Find(bson.M{"sigfox_id": sigfoxId}).All(&deviceMessages)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_messages_not_found", "Device messages not found", err)
	}

	return deviceMessages, nil
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
