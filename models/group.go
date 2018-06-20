package models

type Group struct {
	Id          string   `json:"id" bson:"_id,omitempty" valid:"-"`
	Name        string   `json:"name" bson:"name"`
	UserId      string   `json:"user_id" bson:"user_id"`
	CustomerIds []string `json:"customer_ids" bson:"customer_ids"`
}

const GroupsCollection = "groups"
