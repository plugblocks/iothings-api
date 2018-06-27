package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Device struct {
	Id             string   `json:"id" bson:"_id,omitempty" valid:"-"`
	OrganizationId string   `json:"organization_id" bson:"organization_id"`
	CustomerId     string   `json:"customer_id" bson:"customer_id"`
	Name           string   `json:"name" bson:"name" valid:"-"`
	Type           string   `json:"type" bson:"type"`
	Metadata       Metadata `json:"metadata" bson:"metadata"`
	LastAccess     int64    `json:"last_access" bson:"last_access" valid:"-"`
	Active         bool     `json:"active" bson:"active" valid:"-"`
}

type Metadata struct {
	BleMac   string `json:"ble_mac" bson:"ble_mac" valid:"-"`
	WifiMac  string `json:"wifi_mac" bson:"wifi_mac" valid:"-"`
	SigfoxId string `json:"sigfox_id" bson:"sigfox_id" valid:"-"`
}

func (d *Device) BeforeCreate(user *User) {
	d.Id = bson.NewObjectId().Hex()
	d.LastAccess = time.Now().Unix()
	d.OrganizationId = user.OrganizationId
}

const DevicesCollection = "devices"
