package models

type Order struct {
	Id                  string         `json:"id" bson:"_id,omitempty" valid:"-"`
	Reference           string         `json:"reference" bson:"reference"`
	Status              string         `json:"status" bson:"status"`
	Destination         GeoCoordinates `json:"destination" bson:"destination"`
	DepartureTime       int            `json:"departure_time" bson:"departure_time"`
	ExpectedArrivalTime int            `json:"expected_arrival_time" bson:"expected_arrival_time"`
	CustomerId 			*string 	   `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	ContactPhoneNumber  *string 	   `json:"contact_phone_number" bson:"contact_phone_number"`
	DeviceId            string         `json:"device_id" bson:"device_id"`
	OrganizationId      string         `json:"organization_id" bson:"organization_id"`
}

const OrdersCollection = "orders"
