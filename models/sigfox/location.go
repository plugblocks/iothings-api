package sigfox

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Location struct {
	Id          string  `json:"id" bson:"_id,omitempty" valid:"-"`
	DeviceId    string  `json:"device_id" bson:"device_id"`
	SigfoxId    string  `json:"sigfox_id" bson:"sigfox_id" valid:"-"`
	FrameNumber uint    `json:"frameNumber" bson:"frameNumber" valid:"-"` //Device : (daily frames under 140)
	Timestamp   int64   `json:"timestamp" bson:"timestamp" valid:"-"`
	Latitude    float64 `json:"latitude" bson:"latitude" valid:"-"`
	Longitude   float64 `json:"longitude" bson:"longitude" valid:"-"`
	Radius      float64 `json:"radius" bson:"radius" valid:"-"`
	SpotIt      bool    `json:"spotIt" bson:"spotIt" valid:"-"`
	GPS         bool    `json:"gps" bson:"gps" valid:"-"`
	WiFi        bool    `json:"wifi" bson:"wifi" valid:"-"`
}

type LastLocation struct {
	DeviceId   string   `json:"id" bson:"_id,omitempty" valid:"-"`
	DeviceName string   `json:"name" bson:"name"`
	Location   Location `json:"location" bson:"location"`
}

func (l *Location) BeforeCreate() {
	l.Id = bson.NewObjectId().Hex()
	l.Timestamp = time.Now().Unix()
}

const SigfoxLocationsCollection = "sigfoxLocations"
