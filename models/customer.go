package models

import (
	"github.com/asaskevich/govalidator"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type Customer struct {
	Id             string `json:"id" bson:"_id,omitempty" valid:"-"`
	Firstname      string `json:"first_name" bson:"first_name"`
	Lastname       string `json:"last_name" bson:"last_name"`
	Password       string `json:"password" bson:"password" valid:"required"`
	Email          string `json:"email" bson:"email" valid:"email,required"`
	Phone          string `json:"phone" bson:"phone"`
	Active         bool   `json:"active" bson:"active"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
	ActivationKey  string `json:"activationKey" bson:"activationKey"`
	ResetKey       string `json:"resetKey" bson:"resetKey"`
}

type SanitizedCustomer struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Firstname string `json:"first_name" bson:"first_name"`
	Lastname  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
}

func (customer *Customer) Sanitize() SanitizedCustomer {
	return SanitizedCustomer{customer.Id, customer.Firstname, customer.Lastname, customer.Email}
}

func (customer *Customer) BeforeCreate() error {
	customer.Active = false
	customer.ActivationKey = helpers.RandomString(20)
	customer.Email = strings.ToLower(customer.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "encryption_failed", "Failed to generate the crypted password", err)
	}
	customer.Password = string(hashedPassword)

	_, err = govalidator.ValidateStruct(customer)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}

	return nil
}

const CustomersCollection = "customers"
