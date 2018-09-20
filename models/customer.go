package models

type Customer struct {
	Id             string         `json:"id" bson:"_id"`
	FirstName      string         `json:"first_name" bson:"first_name"`
	LastName       string         `json:"last_name" bson:"last_name"`
	Address        GeoCoordinates `json:"address" bson:"address"`
	Email          string         `json:"email" bson:"email"`
	PhoneNumber    string         `json:"phone_number" bson:"phone_number"`
	OrganizationId string         `json:"organization_id" bson:"organization_id"`
}

const CustomersCollection = "customers"
