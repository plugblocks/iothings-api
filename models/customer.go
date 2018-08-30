package models

type Customer struct {
	Id             string `json:"id" bson:"_id,omitempty" valid:"-"`
	Firstname      string `json:"first_name" bson:"first_name"`
	Lastname       string `json:"last_name" bson:"last_name"`
	Password       string `json:"password" bson:"password" valid:"required"`
	Email          string `json:"email" bson:"email" valid:"email,required"`
	Phone          string `json:"phone" bson:"phone"`
	Active         bool   `json:"active" bson:"active"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
	ActivationKey  string `json:"activationKey" bson:"activationKey"`
	ResetKey       string `json:"resetKey" bson:"resetKey"`
}

const CustomersCollection = "customers"
