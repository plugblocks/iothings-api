package models

type QuantitativeValue struct {
	*SemanticProperty

	Identifier string      `json:"identifier" bson:"identifier"`
	UnitText   string      `json:"unitText" bson:"unitText"`
	Value      interface{} `json:"value" bson:"value"`
}
