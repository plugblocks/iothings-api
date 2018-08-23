package models

type EmailData struct {
	User    *User
	Subject string
	ApiUrl  string
	AppName string
}