package services

import (
				//"net/http"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
		"github.com/spf13/viper"
	"golang.org/x/net/context"
	"github.com/aws/aws-sdk-go/service/sns"
	"gitlab.com/plugblocks/iothings-api/models"
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/store"
	"strconv"
	"gitlab.com/plugblocks/iothings-api/helpers/params"
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
	orga, err := store.GetOrganizationById(c, store.Current(c).OrganizationId)
	if err != nil {
		fmt.Println("Text Check Credit Organization not found", err)
		return false
	}

	mailCredit, _ := strconv.Atoi(orga.PlanCreditTexts)
	fmt.Println("Mail User Creation Organization credit:" + string(mailCredit))
	if mailCredit <= 0 {
		fmt.Println("Text Check Credit Organization no credit")
		return false
	}
	orga.PlanCreditTexts = strconv.Itoa(mailCredit - 1)
	if err := store.UpdateOrganization(c, orga.Id, params.M{"$set": orga}); err != nil {
		c.Error(err)
		c.Abort()
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
		Message: aws.String(data.Message),
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
