package models

type Customer struct {
	Id string `json:"id" bson:"_id,omitempty" valid:"-"`
}

const CustomersCollection = "customers"
