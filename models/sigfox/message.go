package sigfox

type Message struct {
	Id          string `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
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

const SigfoxMessagesCollection = "sigfoxMessages"
