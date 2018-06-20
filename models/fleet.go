package models

type Fleet struct {
	Id        string   `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	DeviceIds []string `json:"device_ids" bson:"device_ids"`
	UserId    string   `json:"user_id" bson:"user_id"`
}

const FleetsCollection = "fleets"
