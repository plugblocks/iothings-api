package models

type GeoJSON struct {
	Type     string    `json:"type" bson:"type"`
	Features []Feature `json:"features" bson:"features"`
}

type Feature struct {
	Type string `json:"type" bson:"type"`
	//Properties Property
	Geometry Geometry `json:"geometry" bson:"geometry"`
}

type Geometry struct {
	Type        string       `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}


/*{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.8044586181640625,
          45.766564985445
        ]
      }
    },
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.884796142578125,
          45.749079020680476
        ]
      }
    }
  ]
}*/
