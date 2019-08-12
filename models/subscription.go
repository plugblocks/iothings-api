package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Subscription struct {
	Id              string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name            string `json:"name" bson:"name"`
	OrganizationId  string `json:"organization_id" bson:"organization_id"`
	Active          bool   `json:"active" bson:"active"`
	LastAccess      int64  `json:"last_access" bson:"last_access" valid:"-"`
	PlanType        string `json:"plan_type" bson:"plan_type"`
	PlanExpiration  int    `json:"plan_expiration" bson:"plan_expiration"`
	PlanCreditMails int    `json:"plan_credit_mail" bson:"plan_credit_mail"`
	PlanCreditTexts int    `json:"plan_credit_text" bson:"plan_credit_text"`
	PlanCreditWifi  int    `json:"plan_credit_wifi" bson:"plan_credit_wifi"`
}

func (s *Subscription) BeforeCreate() {
	s.Id = bson.NewObjectId().Hex()
	s.LastAccess = time.Now().Unix()
}

const SubscriptionsCollection = "subscriptions"
