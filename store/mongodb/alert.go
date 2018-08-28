package mongodb

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
	"gitlab.com/plugblocks/iothings-api/models"
	"net/http"
)

func (db *mongo) CreateAlert(user *models.User, alert *models.Alert) error {
	session := db.Session.Copy()
	defer session.Close()
	alerts := db.C(models.AlertsCollection).With(session)

	alert.BeforeCreate(user)

	err := alerts.Insert(alert)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "alert_creation_failed", "Failed to insert the alert in the database", err)
	}

	return nil
}

func (db *mongo) GetAlert(user *models.User, id string) (*models.Alert, error) {
	session := db.Session.Copy()
	defer session.Close()
	alertCollection := db.C(models.AlertsCollection).With(session)

	alert := &models.Alert{}
	err := alertCollection.Find(bson.M{"_id": id, "user_id": user.Id}).One(alert)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "alert_not_found", "Could not find the alert", err)
	}

	return alert, nil
}

func (db *mongo) GetFleetAlerts(user *models.User, fleetId string) ([]*models.Alert, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetsCollection := db.C(models.FleetsCollection).With(session)
	alertsCollection := db.C(models.AlertsCollection).With(session)

	fleet := &models.Fleet{}
	alerts := []*models.Alert{}

	err := fleetsCollection.Find(bson.M{"_id": fleetId, "user_id": user.Id}).Sort("-timestamp").One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_fleet_not_found", "Failed to find user fleet", err)
	}

	err = alertsCollection.Find(bson.M{"fleet_id": fleet.Id}).All(&alerts)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_alerts_failed", "Failed to get the alerts: "+err.Error(), err)
	}

	return alerts, nil
}

func (db *mongo) GetDeviceAlerts(user *models.User, deviceId string) ([]*models.Alert, error) {
	session := db.Session.Copy()
	defer session.Close()
	devicesCollection := db.C(models.DevicesCollection).With(session)
	alertsCollection := db.C(models.AlertsCollection).With(session)

	device := &models.Device{}
	alerts := []*models.Alert{}

	err := devicesCollection.Find(bson.M{"_id": deviceId, "user_id": user.Id}).Sort("-timestamp").One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_deviuce_not_found", "Failed to find user device", err)
	}

	err = alertsCollection.Find(bson.M{"device_id": device.Id}).All(&alerts)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_alerts_failed", "Failed to get the alerts: "+err.Error(), err)
	}

	return alerts, nil
}

func (db *mongo) UpdateAlert(user *models.User, id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	alertCollection := db.C(models.AlertsCollection).With(session)

	err := alertCollection.Update(bson.M{"_id": id, "user_id": user.Id}, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "alert_update_failed", "Failed to update the alert: "+err.Error(), err)
	}

	return nil
}

func (db *mongo) DeleteAlert(user *models.User, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	alertCollection := db.C(models.AlertsCollection).With(session)

	err := alertCollection.Remove(bson.M{"_id": id, "user_id": user.Id})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "alert_delete_failed", "Failed to delete the alert", err)
	}

	return nil
}

func (db *mongo) CountAlerts() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	alerts := db.C(models.AlertsCollection).With(session)

	nbr, err := alerts.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "alerts_not_found", "Alerts not found", err)
	}
	return nbr, nil
}
