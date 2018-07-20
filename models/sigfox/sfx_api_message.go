package sigfox

type SfxApiReceptionInfos struct {
	Tap       string  `json:"tap" bson:"tap" valid:"-"`
	Delay     float64 `json:"delay" bson:"delay" valid:"-"`
	Latitude  string  `json:"lat" bson:"lat" valid:"-"` //TODO: conversion from string to float64
	Longitude string  `json:"lng" bson:"lng" valid:"-"` //TODO: conversion from string to float64
}

type SfxApiComputedLocation struct {
	Latitude  float64 `json:"lat" bson:"lat" valid:"-"`
	Longitude float64 `json:"lng" bson:"lng" valid:"-"`
	Radius    int64   `json:"radius" bson:"radius" valid:"-"`
	Source    int64   `json:"source" bson:"source" valid:"-"`
}

type SfxApiMessage struct {
	//Nested as a data array
	//Id          string         `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string                 `json:"device" bson:"device" valid:"-"`
	Time        int64                  `json:"time" bson:"time" valid:"-"`
	Data        string                 `json:"data" bson:"data" valid:"-"`
	SequenceNbr uint                   `json:"seqNumber" bson:"seqNumber" valid:"-"`
	RInfos      []SfxApiReceptionInfos `json:"rinfos" bson:"rinfos"`
	CompLoc     SfxApiComputedLocation `json:"computedLocation" bson:"computedLocation"`
	FramesNbr   int64                  `json:"nbFrames" bson:"nbFrames" valid:"-"`
	Operator    string                 `json:"operator" bson:"operator" valid:"-"`
	Country     string                 `json:"country" bson:"country" valid:"-"`
	Snr         string                 `json:"snr" bson:"snr" valid:"-"` //TODO: conversion from string to float64
	LinkQuality string                 `json:"linkQuality" bson:"linkQuality" valid:"-"`
	GroupId     string                 `json:"groupId" bson:"groupId" valid:"-"`
}

type SfxApiNextURL struct {
	SigfoxNextURL string `json:"next" bson:"next"`
}

type SfxApiMessages struct {
	Messages []SfxApiMessage `json:"data" bson:"data"`
	Paging   SfxApiNextURL   `json:"paging" bson:"paging"`
}
