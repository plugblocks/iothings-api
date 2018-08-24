package models

type Subscription struct {
	Id              string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name            string `json:"name" bson:"name"`
	Active          bool   `json:"active" bson:"active"`
	PlanType        string `json:"plan_type" bson:"plan_type"`
	PlanExpiration  int    `json:"plan_expiration" bson:"plan_expiration"`
	PlanCreditMails int    `json:"plan_credit_mail" bson:"plan_credit_mail"`
	PlanCreditTexts int    `json:"plan_credit_text" bson:"plan_credit_text"`
	PlanCreditWifi  int    `json:"plan_credit_wifi" bson:"plan_credit_wifi"`
}

const SubscriptionsCollection = "subscriptions"
