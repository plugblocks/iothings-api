package models

type Property interface {
	GetContext() string
	GetType() string
}

type DefaultProperty struct {
	Context string `json:"@context" bson:"@context"`
	Type    string `json:"@type" bson:"@type"`
}

func (dp *DefaultProperty) GetType() string {
	return dp.Type
}

func (dp *DefaultProperty) GetContext() string {
	return dp.Context
}

func (dp *DefaultProperty) SetContext(ctxt string) {
	dp.Context = ctxt
}

func (dp *DefaultProperty) SetType(typ string) {
	dp.Type = typ
}
