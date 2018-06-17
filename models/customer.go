package models

type Customer struct {
	Id       string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName string `json:"last_name" bson:"last_name"`
}
