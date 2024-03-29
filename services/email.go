package services

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/models"
	"golang.org/x/net/context"
)

const (
	emailSenderKey = "emailSender"
)

func GetEmailSender(c context.Context) EmailSender {
	return c.Value(emailSenderKey).(EmailSender)
}

type EmailSender interface {
	GetEmailParams() *EmailSenderParams
	//SendUserValidationEmail(user *models.User, subject string, templateLink string) error
	//SendAlertEmail(user *models.User, device *models.Device, observation *models.Observation, subject string, templateLink string) error
	SendEmailFromTemplate(ctx *gin.Context, data *models.EmailData, templateLink string) error
	SendEmail(data *models.EmailData) error
}

type FakeEmailSender struct{}

type EmailSenderParams struct {
	senderEmail string
	senderName  string
	apiID       string
	apiKey      string
	apiUrl      string
}

/*func (f *FakeEmailSender) SendEmailFromTemplate(user *models.User, subject string, templateLink string) (error) {
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}
}*/

func NewEmailSender(config *viper.Viper) EmailSender {
	return &EmailSenderParams{
		config.GetString("mail_sender_address"),
		config.GetString("mail_sender_name"),
		config.GetString("aws_api_id"),
		config.GetString("aws_api_key"),
		config.GetString("api_url"),
	}
}

func (s *EmailSenderParams) GetEmailParams() *EmailSenderParams {
	return s
}

func (s *EmailSenderParams) SendEmail(data *models.EmailData) error {
	file, err := ioutil.ReadFile("./templates/html/mail_squeletton.html")
	if err != nil {
		return err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(data.ReceiverMail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(buffer.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(data.Subject),
			},
		},
		Source: aws.String(s.senderEmail),
		// Uncomment to use a configuration set
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		fmt.Println(err)
		return err
	}

	fmt.Println("SES Email Sent to " + data.ReceiverName + " at address: " + data.ReceiverMail)

	return nil
}

func (s *EmailSenderParams) SendEmailFromTemplate(ctx *gin.Context, data *models.EmailData, templateLink string) error {
	file, err := ioutil.ReadFile(templateLink)
	if err != nil {
		return err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(data.ReceiverMail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(buffer.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(data.Subject),
			},
		},
		Source: aws.String(s.senderEmail),
		// Uncomment to use a configuration set
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		fmt.Println(err)
		return err
	}

	//fmt.Println("SES Email Sent to " + data.ReceiverName + " at address: " + data.ReceiverMail)
	//fmt.Println(result)

	return nil
}
