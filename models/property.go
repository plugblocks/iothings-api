package models

type Property interface {
	GetContext() string
	GetType() string
	SetContext(string)
	SetType(string)
}

type SemanticProperty struct {
	Context string `json:"context" bson:"context"`
	Type    string `json:"type" bson:"type"`
}

func (dp *SemanticProperty) GetType() string {
	return dp.Type
}

func (dp *SemanticProperty) GetContext() string {
	return dp.Context
}

func (dp *SemanticProperty) SetContext(ctxt string) {
	dp.Context = ctxt
}

func (dp *SemanticProperty) SetType(typ string) {
	dp.Type = typ
}
