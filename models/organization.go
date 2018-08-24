package models

type Organization struct {
	Id     string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name   string `json:"name" bson:"name"`
	Active bool   `json:"active" bson:"active"`
	Image  string `json:"image" bson:"image"`
	Admin  bool   `json:"admin" bson:"admin"`
}

const OrganizationsCollection = "organizations"
