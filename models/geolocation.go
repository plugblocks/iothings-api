package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type GeoJSON struct {
	Type     string    `json:"type" bson:"type"`
	Features []Feature `json:"features" bson:"features"`
}

type Feature struct {
	Type string `json:"type" bson:"type"`
	//Properties Property
	Geometry Geometry `json:"geometry" bson:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

/*{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.8044586181640625,
          45.766564985445
        ]
      }
    },
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.884796142578125,
          45.749079020680476
        ]
      }
    }
  ]
}*/

type Geolocation struct {
	Id        string  `json:"id" bson:"_id,omitempty" valid:"-"`
	DeviceId  string  `json:"device_id" bson:"device_id"`
	Timestamp int64   `json:"timestamp" bson:"timestamp" valid:"-"`
	Latitude  float64 `json:"latitude" bson:"latitude" valid:"-"`
	Longitude float64 `json:"longitude" bson:"longitude" valid:"-"`
	Radius    float64 `json:"radius" bson:"radius" valid:"-"`
	Source    string  `json:"source" bson:"source" valid:"-"`
}

type GeolocationQueryParams struct {
	Order     bool   `form:"order" json:"order"`
	Limit     int    `form:"limit" json:"limit"`
	Source    string `form:"source" json:"source"`
	StartTime int    `form:"starttime" json:"starttime"`
	EndTime   int    `form:"endtime" json:"endtime"`
}

func (l *Geolocation) BeforeCreate() {
	l.Id = bson.NewObjectId().Hex()
	if l.Timestamp == 0 {
		l.Timestamp = time.Now().Unix()
	}
}

const GeolocationsCollection = "geolocations"
