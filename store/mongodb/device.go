package mongodb

import (
	"net/http"

	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
)

func (db *mongo) CreateDevice(user *models.User, device *models.Device) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	fmt.Println("Mongo User: " + user.Email + "OrgaId" + user.OrganizationId)
	device.BeforeCreate(user)

	if device.SigfoxId != "" {
		count, _ := devices.Find(bson.M{"sigfox_id": device.SigfoxId}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "Failed to create the device Sigfox", errors.New("Sigfox Device already exists"))
		}
	}

	if device.BleMac != "" {
		count, _ := devices.Find(bson.M{"ble_mac": device.BleMac}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "Failed to create the device BLE", errors.New("BLE Device already exists"))
		}
	}

	if device.WifiMac != "" {
		count, _ := devices.Find(bson.M{"wifi_mac": device.WifiMac}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_creation_failed", "Failed to create the device WiFi", errors.New("WiFi Device already exists"))
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

//TODO: DANGER: Protect by auth device GeoJSON
func (db *mongo) GetDeviceGeoJSON( /*user *models.User, */ deviceId string) (*models.GeoJSON, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.Find(bson.M{"_id": deviceId}).One(device) //TODO: add security w/user organization_id , "organization_id": user.OrganizationId
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found", err)
	}

	//////////////////////////////////////

	locationCollection := db.C(sigfox.SigfoxLocationsCollection).With(session)
	locations := []sigfox.Location{}
	err = locationCollection.Find(bson.M{"sigfox_id": device.SigfoxId}).Sort("-timestamp").All(&locations)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	//TODO: use observations with trick to find values in observation
	// if(locationObs.Values[0].DefaultProperty.Type) == "location"
	/*
		observations := db.C(models.ObservationsCollection).With(session)
		locationsObservations := []*models.Observation{}
		//TODO: customize source in URL: spotit, wifi, gps ...
		err = observations.Find(bson.M{"device_id": deviceId, "resolver": "spotit"}).All(&locationsObservations)
		if err != nil {
			fmt.Println("device get obs id:", deviceId," err:", err)
			return nil, helpers.NewError(http.StatusNotFound, "observations_device_not_found", "Failed to find observations for device", err)
		}*/
	features := []models.Feature{}

	for _, location := range locations {
		coords := []float64{}
		coords = append(coords, location.Longitude, location.Latitude)

		geometry := models.Geometry{"Point", coords}
		feature := models.Feature{"Feature", geometry}

		features = append(features, feature)
	}
	geojson := &models.GeoJSON{"FeatureCollection", features}

	//fmt.Println("GetDeviceGeoJSON: device: ", device, "\t locations:", locations)

	return geojson, nil
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
