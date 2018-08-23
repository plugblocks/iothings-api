package models

type Organization struct {
	Id              string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name            string `json:"name" bson:"name"`
	Active          bool   `json:"active" bson:"active"`
	Image           string `json:"image" bson:"image"`
	Admin           bool   `json:"admin" bson:"admin"`
	PlanType        string `json:"plan_type" bson:"plan_type"`
	// TODO: Improve front to send int values
	PlanExpiration  string    `json:"plan_expiration" bson:"plan_expiration"`
	PlanCreditTexts string    `json:"plan_credit_texts" bson:"plan_credit_texts"`
	PlanCreditMails string    `json:"plan_credit_mails" bson:"plan_credit_mails"`
	PlanCreditWifi  string    `json:"plan_credit_wifi" bson:"plan_credit_wifi"`
}

const OrganizationsCollection = "organizations"
