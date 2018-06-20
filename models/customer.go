package models

type Customer struct {
	Id        string   `json:"id" bson:"_id,omitempty" valid:"-"`
	DeviceIds []string `json:"device_ids" bson:"device_ids"`
}

const CustomersCollection = "customers"
