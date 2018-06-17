package models

type Observation struct {
	Id string `json:"id" bson:"_id"`
	CustomerId string `json:"customer_id" bson:"customer_id"`
	Properties []Property `json:"properties" bson:"properties"`
	DeviceId string `json:"device_id" bson:"device_id"`
}

const ObservationsCollection = "observations"