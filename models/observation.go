package models

import (
	"time"
)

type Observation struct {
	Id        string              `json:"id" bson:"_id,omitempty" valid:"-"`
	Timestamp int64               `json:"timestamp" bson:"timestamp" valid:"-"`
	DeviceId  string              `json:"device_id" bson:"device_id"`
	Resolver  string              `json:"resolver" bson:"resolver" valid:"-"`
	Values    []QuantitativeValue `json:"values" bson:"values"`
}

func (o *Observation) BeforeCreate(device *Device) {
	device.LastAccess = time.Now().Unix()
}

const ObservationsCollection = "observations"
