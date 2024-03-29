package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type QuantitativeValue struct {
	*SemanticProperty

	Identifier string      `json:"identifier" bson:"identifier"`
	UnitText   string      `json:"unitText" bson:"unitText"`
	Value      interface{} `json:"value" bson:"value"`
}

type ObservationQueryParams struct {
	Order     string `form:"order" json:"order"`
	Limit     int    `form:"limit" json:"limit"`
	Resolver  string `form:"resolver" json:"resolver"`
	StartTime int    `form:"starttime" json:"starttime"`
	EndTime   int    `form:"endtime" json:"endtime"`
}

type Observation struct {
	Id        string              `json:"id" bson:"_id,omitempty" valid:"-"`
	Timestamp int64               `json:"timestamp" bson:"timestamp" valid:"-"`
	DeviceId  string              `json:"device_id" bson:"device_id"`
	Resolver  string              `json:"resolver" bson:"resolver" valid:"-"`
	Values    []QuantitativeValue `json:"values" bson:"values"`
}

func (o *Observation) BeforeCreate(device *Device) {
	o.Id = bson.NewObjectId().Hex()
	if o.Timestamp == 0 {
		o.Timestamp = time.Now().Unix()
	}
	device.LastAccess = time.Now().Unix()
	device.Active = true
}

const ObservationsCollection = "observations"
