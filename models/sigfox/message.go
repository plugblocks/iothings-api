package sigfox

type Message struct {
	Id          string `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string `json:"sigfox_id" bson:"sigfox_id" valid:"-"`
	FrameNumber uint   `json:"frameNumber" bson:"frameNumber" valid:"-"` //Sigfox: (daily frames under 140)
	Timestamp   int64  `json:"timestamp" bson:"timestamp" valid:"-"`     //Sigfox: time
	/* These values will be deprecated on June 1st 2019 */
	Snr     float64 `json:"snr" bson:"snr" valid:"-"`         //Sigfox: snr
	Rssi    float64 `json:"rssi" bson:"rssi" valid:"-"`       //Sigfox: rssi
	AvgSnr  float64 `json:"avgSnr" bson:"avgSnr" valid:"-"`   //Sigfox: avgSnr
	Station string  `json:"station" bson:"station" valid:"-"` //Sigfox: station
	Lat     int8    `json:"lat" bson:"lat" valid:"-"`         //
	Lng     int8    `json:"lng" bson:"lng" valid:"-"`
	/* End of deprecate */
	Resolver string `json:"resolver" bson:"resolver" valid:"-"` //Custom: message type to dispatch cases
	Data     string `json:"data" bson:"data" valid:"-"`         //Sigfox: data
	//Ack for downlink
}

type MessageDataAdvanced struct {
	Id               string           `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId         string           `json:"sigfox_id" bson:"sigfox_id" valid:"-"`
	Timestamp        int64            `json:"timestamp" bson:"timestamp" valid:"-"`
	Data             string           `json:"data" bson:"data" valid:"-"`
	SeqNumber        int              `json:"seqNumber" bson:"seqNumber" valid:"-"`
	Lqi              string           `json:"lqi" bson:"lqi" valid:"-"`
	OperatorName     string           `json:"operatorName" bson:"operatorName" valid:"-"`
	CountryCode      string           `json:"countryCode" bson:"countryCode" valid:"-"`
	ComputedLocation ComputedLocation `json:"computedLocation" bson:"computedLocation" valid:"-"`
	Resolver         string           `json:"resolver,omitempty" bson:"resolver" valid:"-"` //Custom: message type to dispatch cases
}

type ComputedLocation struct {
	Lat    float64 `json:"lat" bson:"lat" valid:"-"`
	Lng    float64 `json:"lng" bson:"lng" valid:"-"`
	Radius int     `json:"radius" bson:"radius" valid:"-"`
	Source int     `json:"source" bson:"source" valid:"-"`
	Status int     `json:"status" bson:"status" valid:"-"`
}

const SigfoxMessagesCollection = "sigfoxMessages"
