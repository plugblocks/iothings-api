package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateGeolocation(location *models.Geolocation) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(models.GeolocationsCollection).With(session)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "geolocation_creation_failed", "Failed to insert the geolocation", err)
	}

	return nil
}

//TODO: DANGER: Protect by auth device GeoJSON
func (db *mongo) GetDeviceGeolocation(user *models.User, deviceId string, source string) (*models.Geolocation, error) {
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
	location := &models.Geolocation{}

	sources := []string{"gps", "wifi"}
	err = geolocationCollection.Find(bson.M{"device_id": deviceId, "source": bson.M{"$in": sources}}).Sort("-timestamp").One(location)

	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	return location, nil
}

//TODO: DANGER: Protect by auth device GeoJSON
func (db *mongo) GetDeviceGeoJSON( /*user *models.User, */ deviceId string, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
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

	err = geolocationCollection.Find(bson.M{"device_id": deviceId, "source": source, "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&locations)

	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_locations_failed", "Failed to get the locations: "+err.Error(), err)
	}

	//TODO: use observations with trick to find values in observation
	// if(locationObs.Values[0].SemanticProperty.Type) == "location"
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
func (db *mongo) GetFleetGeoJSON( /*user *models.User, */ fleetId string, source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
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
		err = geolocationCollection.Find(bson.M{"device_id": deviceId, "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&locations)

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
	// if(locationObs.Values[0].SemanticProperty.Type) == "location"
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

//TODO: DANGER: Protect by auth device GeoJSON
func (db *mongo) GetFleetsGeoJSON(source string, limit int, startTime int, endTime int) (*models.GeoJSON, error) {
	session := db.Session.Copy()
	defer session.Close()

	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleets := []models.Fleet{}
	err := fleetCollection.Find(bson.M{}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_fleets_failed", "Failed to get the fleets: "+err.Error(), err)
	}

	//////////////////////////////////////
	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	locations := []models.Geolocation{}
	features := []models.Feature{}

	for _, fleet := range fleets {
		for _, deviceId := range fleet.DeviceIds {
			err = geolocationCollection.Find(bson.M{"device_id": deviceId, "timestamp": bson.M{"$gt": startTime, "$lt": endTime}}).Sort("-timestamp").Limit(limit).All(&locations)
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
	}

	//TODO: use observations with trick to find values in observation
	// if(locationObs.Values[0].SemanticProperty.Type) == "location"
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

func (db *mongo) GetUserFleetsGeoJSON(user *models.User) (*models.GeoJSON, error) {
	session := db.Session.Copy()
	defer session.Close()

	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleets := []models.Fleet{}
	err := fleetCollection.Find(bson.M{"user_id": user.Id}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_fleets_failed", "Failed to get the fleets: "+err.Error(), err)
	}

	//////////////////////////////////////
	geolocationCollection := db.C(models.GeolocationsCollection).With(session)
	locations := []models.Geolocation{}
	features := []models.Feature{}

	for _, fleet := range fleets {
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
	}

	//TODO: use observations with trick to find values in observation
	// if(locationObs.Values[0].SemanticProperty.Type) == "location"
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

func (db *mongo) CountGeolocations() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	geolocations := db.C(models.GeolocationsCollection).With(session)

	nbr, err := geolocations.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "geolocations_not_found", "Geolocations not found", err)
	}
	return nbr, nil
}
