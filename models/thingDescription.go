package models

type TD struct {
	Context             string       `json:"@context" bson:"@context"`
	Type                string       `json:"@type" bson:"@type"`
	Id                  string       `json:"id" bson:"id"`
	Name                string       `json:"name" bson:"name"`
	Description         string       `json:"description" bson:"description"`
	SecurityDefinitions string       `json:"securityDefinitions" bson:"securityDefinitions"`
	Security            string       `json:"security" bson:"security"`
	Properties          []TDProperty `json:"properties" bson:"properties"`
	Actions             []TDAction   `json:"actions" bson:"actions"`
	Events              []TDEvent    `json:"events" bson:"events"`
}

type TDProperty struct {
	Name               string             `json:"name" bson:"name"`
	TDPropertyElements TDPropertyElements `json:"TDPropertyElements" bson:"TDPropertyElements"`
}

type TDPropertyElements struct {
	AType       string `json:"@type" bson:"@type"`
	Description string `json:"description" bson:"description"`
	ReadOnly    bool   `json:"readOnly" bson:"readOnly"`
	Observable  bool   `json:"observable" bson:"observable"`
	Type        string `json:"type" bson:"type"`
	Forms       TDForm `json:"forms" bson:"forms"`
}

type TDForm struct {
	Href        string `json:"href" bson:"href"`
	ContentType string `json:"contentType" bson:"contentType"`
}

type TDAction struct {
	Name             string           `json:"name" bson:"name"`
	TDActionElements TDActionElements `json:"TDActionElements" bson:"TDActionElements"`
}

type TDActionElements struct {
	AType       string `json:"@type" bson:"@type"`
	Description string `json:"description" bson:"description"`
	Forms       TDForm `json:"forms" bson:"forms"`
}

type TDEvent struct {
	Name            string          `json:"name" bson:"name"`
	TDEventElements TDEventElements `json:"TDEventElements" bson:"TDEventElements"`
}

type TDData struct {
	Type string `json:"type" bson:"type"`
}

type TDEventElements struct {
	AType       string `json:"@type" bson:"@type"`
	Description string `json:"description" bson:"description"`
	Data        TDData `json:"data" bson:"data"`
	Forms       TDForm `json:"forms" bson:"forms"`
}

/* Example:
{
    "@context": ["http://www.w3.org/ns/td",
    		{"iot": "http://iotschema.org/"}],
    "@type" : "Thing",
    "id": "urn:dev:wot:com:example:servient:lamp",
    "name": "MyLampThing",
    "description" : "MyLampThing uses JSON-LD 1.1 serialization",
    "securityDefinitions": {"psk_sc":{"scheme": "psk"}},
    "security": ["psk_sc"],
    "properties": {
        "status": {
            "@type" : "iot:SwitchStatus",
            "description" : "Shows the current status of the lamp",
            "readOnly": true,
            "observable": false,
            "type": "string",
            "forms": [{
                "href": "coaps://mylamp.example.com/status",
                "contentType": "application/json"
            }]
        }
    },
    "actions": {
        "toggle": {
            "@type" : "iot:SwitchStatus",
            "description" : "Turn on or off the lamp",
            "forms": [{
                "href": "coaps://mylamp.example.com/toggle",
                "contentType": "application/json"
            }]
        }
    },
    "events": {
        "overheating": {
            "@type" : "iot:TemperatureAlarm",
            "description" : "Lamp reaches a critical temperature (overheating)",
            "data": {"type": "string"},
            "forms": [{
                "href": "coaps://mylamp.example.com/oh",
                "contentType": "application/json"
            }]
        }
    }
}
*/
