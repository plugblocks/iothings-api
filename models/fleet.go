package models

type Fleet struct {
	Id        string   `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	DeviceIds []string `json:"device_ids" bson:"device_ids"`
}

const FleetsCollection = "fleets"
