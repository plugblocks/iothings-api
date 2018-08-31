package models

type EmailData struct {
	ReceiverMail string
	ReceiverName string
	User         *User
	Customer     *Customer
	Subject      string
	Body         string
	ApiUrl       string
	AppName      string
}

type TextData struct {
	User    *User
	Subject string
	Message string
}
