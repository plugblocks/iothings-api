package models

type Group struct {
	Id       string   `json:"id" bson:"_id,omitempty" valid:"-"`
	Name     string   `json:"name" bson:"name"`
	AdminIds []string `json:"admin_ids" bson:"admin_ids"`
	UserIds  []string `json:"user_ids" bson:"user_ids"`
}

const GroupsCollection = "groups"