package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Device struct {
	Id             string `json:"id" bson:"_id,omitempty" valid:"-"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
	CustomerId     string `json:"customer_id" bson:"customer_id"`
	Name           string `json:"name" bson:"name" valid:"-"`
	BleMac         string `json:"ble_mac" bson:"ble_mac" valid:"-"`
	WifiMac        string `json:"wifi_mac" bson:"wifi_mac" valid:"-"`
	SigfoxId       string `json:"sigfox_id" bson:"sigfox_id" valid:"-"`
	LastAccess     int64  `json:"last_access" bson:"last_access" valid:"-"`
	Activity       int64  `json:"activity" bson:"activity" valid:"-"`
	Active         bool   `json:"active" bson:"active" valid:"-"`
	Available 	   bool   `json:"available" bson:"available"`
}

func (d *Device) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAccess = time.Now().Unix()
}

const DevicesCollection = "devices"
