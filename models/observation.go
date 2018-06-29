package models

import (
	"time"
	"gitlab.com/plugblocks/iothings-api/models/schema_org"
)

type Observation struct {
	Id         string     `json:"id" bson:"_id"`
	Timestamp  int64      `json:"timestamp" bson:"timestamp" valid:"-"`
	DeviceId   string     `json:"device_id" bson:"device_id"`
	Type       string     `json:"type" bson:"type"`
	Values []schema_org.QuantitativeValue `json:"values" bson:"values"`
}

func (o *Observation) BeforeCreate(device *Device) {
	device.LastAccess = time.Now().Unix()
}

const ObservationsCollection = "observations"
