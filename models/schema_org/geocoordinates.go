package schema_org

import "gitlab.com/plugblocks/iothings-api/models"

type GeoCoordinates struct {
	*models.DefaultProperty

	Geo  Geo    `json:"geo" bson:"geo"`
	Name string `json:"name" bson:"name"`
}

type Geo struct {
	Type      string `json:"@type" bson:"@type"`
	Latitude  string `json:"latitude" bson:"latitude"`
	Longitude string `json:"longitude" bson:"longitude"`
}
