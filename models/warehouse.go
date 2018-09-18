package models

type Warehouse struct {
	Id          string         `json:"id" bson:"_id,omitempty"`
	Name        string         `json:"name" bson:"name"`
	Coordinates GeoCoordinates `json:"coordinates" bson:"coordinates"`
	Address     string         `json:"address" bson:"address"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
}

const WarehousesCollection = "warehouses"

