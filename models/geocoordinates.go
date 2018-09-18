package models

type GeoCoordinates struct {
	*SemanticProperty

	Latitude  string `json:"latitude" bson:"latitude"`
	Longitude string `json:"longitude" bson:"longitude"`
	Name      *string `json:"name,omitempty" bson:"name,omitempty"`
}
