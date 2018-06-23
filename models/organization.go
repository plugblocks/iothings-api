package models

type Organization struct {
	Id     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Active bool   `json:"active" bson:"active"`
}

const OrganizationsCollection = "organizations"
