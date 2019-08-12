package sigfox

type SfxApiLocation struct {
	//Id string `json:"id" bson:"_id,omitempty" valid:"-"`
	//SigfoxId    string  `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	Time      int64   `json:"time" bson:"time" valid:"-"`
	Valid     bool    `json:"valid" bson:"valid" valid:"-"`
	Latitude  float64 `json:"lat" bson:"lat" valid:"-"`
	Longitude float64 `json:"lng" bson:"lng" valid:"-"`
	Radius    int32   `json:"radius" bson:"radius" valid:"-"`
	//TODO: Handle double source from API
	/*SourceInt    int8  `json:"source" bson:"source" valid:"-"`
	SourceStr    string  `json:"source" bson:"source" valid:"-"`*/
}

type SfxApiLocations struct {
	Messages []SfxApiLocation `json:"data" bson:"data"`
}
