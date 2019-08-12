package models

import "github.com/globalsign/mgo/bson"

type Order struct {
	Id                  string         `json:"id" bson:"_id,omitempty" valid:"-"`
	Reference           string         `json:"reference" bson:"reference"`
	Status              string         `json:"status" bson:"status"`
	Destination         GeoCoordinates `json:"destination" bson:"destination"`
	DepartureTime       int            `json:"departure_time" bson:"departure_time"`
	ExpectedArrivalTime int            `json:"expected_arrival_time" bson:"expected_arrival_time"`
	CustomerId          *string        `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	ContactPhoneNumber  string         `json:"contact_phone_number" bson:"contact_phone_number"`
	DeviceId            string         `json:"device_id" bson:"device_id"`
	OrganizationId      string         `json:"organization_id" bson:"organization_id"`
	HasNotifiedDelay    bool           `json:"has_notified_delay" bson:"has_notified_delay"`
	LiveETA             int64          `json:"live_eta" bson:"live_eta"`
}

func (o *Order) BeforeCreate() {
	o.Id = bson.NewObjectId().Hex()
	o.HasNotifiedDelay = false
	o.Status = Created.String()
}

type Status int

const (
	Created Status = iota
	Transiting
	Arrived
	Terminated
)

func (s Status) String() string {
	return [...]string{"created", "transiting", "arrived", "terminated"}[s]
}

const OrdersCollection = "orders"
