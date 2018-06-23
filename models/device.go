package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Device struct {
	Id         string `json:"id" bson:"_id,omitempty" valid:"-"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
	Name       string `json:"name" bson:"name" valid:"-"`
	Type       string `json:"type" bson:"type"`
	LastAccess int64  `json:"last_access" bson:"last_access" valid:"-"`
	Active     bool   `json:"active" bson:"active" valid:"-"`
}

func (d *Device) BeforeCreate(user *User) {
	d.Id = bson.NewObjectId().Hex()
	d.LastAccess = time.Now().Unix()
	d.OrganizationId = user.Id
}

const DevicesCollection = "devices"
