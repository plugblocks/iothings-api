package models

import (
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             string `json:"id" bson:"_id,omitempty" valid:"-"`
	Firstname      string `json:"first_name" bson:"first_name"`
	Lastname       string `json:"last_name" bson:"last_name"`
	Password       string `json:"password" bson:"password" valid:"required"`
	Email          string `json:"email" bson:"email" valid:"email,required"`
	Phone          string `json:"phone" bson:"phone"`
	Active         bool   `json:"active" bson:"active"`
	Language       string `json:"language" bson:"language"`
	OrganizationId string `json:"organization_id" bson:"organization_id"`
	ActivationKey  string `json:"activationKey" bson:"activationKey"`
	ResetKey       string `json:"resetKey" bson:"resetKey"`
	Admin          bool   `json:"admin" bson:"admin"`
	LastAccess     int64  `json:"last_access" bson:"last_access" valid:"-"`
}

type SanitizedUser struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Firstname string `json:"first_name" bson:"first_name"`
	Lastname  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Admin     bool   `json:"admin" bson:"admin"`
}

func (user *User) Sanitize() SanitizedUser {
	return SanitizedUser{user.Id, user.Firstname, user.Lastname, user.Email, user.Admin}
}

func (user *User) BeforeCreate() error {
	user.Active = false
	user.ActivationKey = helpers.RandomString(20)
	user.Email = strings.ToLower(user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "encryption_failed", "Failed to generate the crypted password", err)
	}
	user.Password = string(hashedPassword)

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}

	return nil
}

const UsersCollection = "users"
