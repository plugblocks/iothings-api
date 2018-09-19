package models

type GeoCoordinates struct {
	*SemanticProperty

	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Name      *string `json:"name,omitempty" bson:"name,omitempty"`
}
