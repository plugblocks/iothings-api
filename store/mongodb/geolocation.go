package mongodb

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateGeolocation(location *models.Geolocation) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(models.GeolocationsCollection).With(session)

	fmt.Println("CreateGeolocation", location)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "geolocation_creation_failed", "Failed to insert the geolocation", err)
	}

	return nil
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

	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	locations := []models.Geolocation{}
	err = geolocationCollection.Find(bson.M{"device_id": deviceId}).Sort("-timestamp").All(&locations)
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

//TODO: DANGER: Protect by auth device GeoJSON
func (db *mongo) GetFleetGeoJSON( /*user *models.User, */ fleetId string) (*models.GeoJSON, error) {
	session := db.Session.Copy()
	defer session.Close()

	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleet := &models.Fleet{}
	err := fleetCollection.Find(bson.M{"_id": fleetId}).One(fleet) //TODO: add security w/user organization_id , "organization_id": user.OrganizationId
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Could not find the fleet", err)
	}

	//////////////////////////////////////
	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	locations := []models.Geolocation{}
	features := []models.Feature{}

	for _, deviceId := range fleet.DeviceIds {
		err = geolocationCollection.Find(bson.M{"device_id": deviceId}).Sort("-timestamp").All(&locations)
		if err != nil {
			return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
		}

		for _, location := range locations {
			coords := []float64{}
			coords = append(coords, location.Longitude, location.Latitude)

			geometry := models.Geometry{"Point", coords}
			feature := models.Feature{"Feature", geometry}

			features = append(features, feature)
		}
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

	geojson := &models.GeoJSON{"FeatureCollection", features}

	//fmt.Println("GetFleetGeoJSON: fleet: ", fleet, "\t locations:", locations)

	return geojson, nil
}
