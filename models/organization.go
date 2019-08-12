package models

type Organization struct {
	Id     string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name   string `json:"name" bson:"name"`
	Siret  uint64 `json:"siret" bson:"siret"`
	Image  string `json:"image" bson:"image"`
	Admin  bool   `json:"admin" bson:"admin"`
	Active bool   `json:"active" bson:"active"`
}

const OrganizationsCollection = "organizations"
