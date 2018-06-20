package models

type Observation struct {
	Id         string     `json:"id" bson:"_id"`
	DeviceId   string     `json:"device_id" bson:"device_id"`
	Type       string     `json:"type" bson:"type"`
	Properties []Property `json:"properties" bson:"properties"`
}

const ObservationsCollection = "observations"
