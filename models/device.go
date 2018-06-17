package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Device struct {
	Id       string `json:"id" bson:"_id,omitempty" valid:"-"`
	CustomerId   string `json:"customer_id" bson:"customer_id" valid:"-"`
	Name     string `json:"name" bson:"name" valid:"-"`
	Type 	 string `json:"type" bson:"type"`
	LastAccess  int64  `json:"last_access" bson:"last_access" valid:"-"`
	Active   bool   `json:"active" bson:"active" valid:"-"`
}

func (d *Device) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAccess = time.Now().Unix()
}

const DevicesCollection = "devices"
