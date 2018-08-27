package services

import (
	//"net/http"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/models"
	"golang.org/x/net/context"
)

const (
	textSenderKey = "textSender"
)

func GetTextSender(c context.Context) TextSender {
	return c.Value(textSenderKey).(TextSender)
}

type TextSender interface {
	SendAlertText(user *models.User, device *models.Device, observation *models.Observation, subject string, templateLink string) error
	CheckTextCredit(c *gin.Context) bool
	SendText(data TextData) error
}

type FakeTextSender struct{}

type TextSenderParams struct {
	senderEmail string
	senderName  string
	apiID       string
	apiKey      string
	apiUrl      string
}

type TextData struct {
	User    *models.User
	Message string
	ApiUrl  string
	AppName string
}

/*func (f *FakeTextSender) SendEmailFromTemplate(user *models.User, subject string, templateLink string) (error) {
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}
}*/

func NewTextSender(config *viper.Viper) TextSender {
	return &TextSenderParams{
		config.GetString("mail_sender_address"),
		config.GetString("mail_sender_name"),
		config.GetString("aws_sns_api_id"),
		config.GetString("aws_sns_api_key"),
		config.GetString("api_url"),
	}
}

func (s *TextSenderParams) SendAlertText(user *models.User, device *models.Device, observation *models.Observation, message string, templateLink string) error {
	data := TextData{User: user, Message: message, ApiUrl: s.apiUrl, AppName: s.senderName}
	s.SendText(data)

	return nil
}

func (s *TextSenderParams) CheckTextCredit(c *gin.Context) bool {
	textCredit := config.GetInt(c, "plan_credit_text")
	fmt.Println("Text Organization credit:", textCredit)
	es := GetEmailSender(c)
	if textCredit > 0 {
		config.Set(c, "plan_credit_text", textCredit-1)
		return true
	} else if textCredit == 0 {
		fmt.Println("Text Check Credit Organization no credit warning mails sent")
		appName := config.GetString(c, "mail_sender_name")
		subject := appName + ", your texts token is empty, we give you 10 texts"
		templateLink := "./templates/html/mail_token_empty.html"
		userData := models.EmailData{ReceiverMail: EmailSender.GetEmailParams(es).senderEmail, ReceiverName: EmailSender.GetEmailParams(es).senderName, Subject: subject, Body: "Texts", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		adminData := models.EmailData{ReceiverMail: "contact@plugblocks.com", ReceiverName: "PlugBlocks Admin", Subject: subject, Body: "Texts", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		EmailSender.SendEmailFromTemplate(es, c, &userData, templateLink)
		EmailSender.SendEmailFromTemplate(es, c, &adminData, templateLink)
		config.Set(c, "plan_credit_text", -1)
		return false
	} else if textCredit > -10 {
		config.Set(c, "plan_credit_text", -100)
		return false
	} else if textCredit == -100 {
		fmt.Println("Text Check Credit Organization no credit disable wifi sent")
		appName := config.GetString(c, "mail_sender_name")
		subject := appName + ", your texts token is empty"
		templateLink := "./templates/html/mail_token_empty.html"
		userData := models.EmailData{ReceiverMail: EmailSender.GetEmailParams(es).senderEmail, ReceiverName: EmailSender.GetEmailParams(es).senderName, Subject: subject, Body: "Texts", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		adminData := models.EmailData{ReceiverMail: "contact@plugblocks.com", ReceiverName: "PlugBlocks Admin", Subject: subject, Body: "Texts", ApiUrl: config.GetString(c, "api_url"), AppName: config.GetString(c, "mail_sender_name")}
		EmailSender.SendEmailFromTemplate(es, c, &userData, templateLink)
		EmailSender.SendEmailFromTemplate(es, c, &adminData, templateLink)
		config.Set(c, "plan_credit_text", -1000)
		return false
	}
	return true
}

func (s *TextSenderParams) SendText(data TextData) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	fmt.Println("Amazon Creds: " + s.apiID + s.apiKey)
	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := sns.New(sess, &aws.Config{Credentials: creds})

	// Assemble the text.

	// Attempt to send the email.
	params := &sns.PublishInput{
		Message:     aws.String(data.Message),
		PhoneNumber: aws.String(data.User.Phone),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("SNS Text Sent to " + data.User.Firstname + " " + data.User.Lastname + " at number: " + data.User.Phone)
	fmt.Println(resp)

	return nil
}
