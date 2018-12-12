package models

import (
	"github.com/globalsign/mgo/bson"
)

type Alert struct {
	Id       string  `json:"id" bson:"_id,omitempty" valid:"-"`
	UserId   string  `json:"user_id" bson:"user_id"`
	FleetId  string  `json:"fleet_id" bson:"fleet_id"`
	DeviceId string  `json:"device_id" bson:"device_id"`
	Name     string  `json:"name" bson:"name" valid:"-"`

	Type 	 string  `json:"type" bson:"type" valid:"-"`

	Property string  `json:"property" bson:"property" valid:"-"`
	Trigger  string  `json:"trigger" bson:"trigger" valid:"-"`
	Value    float64 `json:"value" bson:"value" valid:"-"`
	Channel  string  `json:"channel" bson:"channel" valid:"-"`

	Active   bool    `json:"active" bson:"active" valid:"-"`
}

func (a *Alert) BeforeCreate(user *User) {
	a.Id = bson.NewObjectId().Hex()
	a.UserId = user.Id
}

const AlertsCollection = "alerts"
