package schema_org

import "gitlab.com/plugblocks/iothings-api/models"

type QuantitativeValue struct {
	*models.DefaultProperty

	Identifier string `json:"identifier" bson:"identifier"`
	UnitText string `json:"unitText" bson:"unitText"`
	Value interface{} `json:"value" bson:"value"`
}