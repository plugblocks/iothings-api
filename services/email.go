package services

import (
	"bytes"
	"html/template"
	"io/ioutil"
	//"net/http"

	"gitlab.com/plugblocks/iothings-api/models"

	/*"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go
	"github.com/sendgrid/sendgrid-go/helpers/mail"*/

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	emailSenderKey = "emailSender"
)

func GetEmailSender(c context.Context) EmailSender {
	return c.Value(emailSenderKey).(EmailSender)
}

type EmailSender interface {
	SendEmailFromTemplate(user *models.User, subject string, templateLink string) (error)
}

type FakeEmailSender struct{}

type EmailSenderParams struct {
	senderEmail string
	senderName  string
	apiID       string
	apiKey      string
	apiUrl      string
}

type Data struct {
	User    *models.User
	ApiUrl  string
	AppName string
}

/*func (f *FakeEmailSender) SendEmailFromTemplate(user *models.User, subject string, templateLink string) (error) {
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}
}*/

func NewEmailSender(config *viper.Viper) EmailSender {
	return &EmailSenderParams{
		/*config.GetString("sendgrid_address"),
		config.GetString("sendgrid_name"),
		config.GetString("sendgrid_api_key"),*/
		config.GetString("mail_sender_address"),
		config.GetString("mail_sender_name"),
		config.GetString("aws_ses_api_id"),
		config.GetString("aws_ses_api_key"),
		config.GetString("api_url"),
	}
}
func (s *EmailSenderParams) SendEmailFromTemplate(user *models.User, subject string, templateLink string) error {
	// Sendgrid Way
	/*to := mail.NewEmail(user.Firstname, user.Email)

	file, err := ioutil.ReadFile(templateLink)
	if err != nil {
		return nil, err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	data := Data{User: user, HostAddress: s.baseUrl, AppName: s.senderName}
	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return s.SendEmail([]*mail.Email{to}, "text/html", subject, buffer.String())*/

	file, err := ioutil.ReadFile(templateLink)
	if err != nil {
		return err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	fmt.Println("SendEmailFromTemplate: Template parsed")

	data := Data{User: user, ApiUrl: s.apiUrl, AppName: s.senderName}
	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	fmt.Println("Amazon Creds: " + s.apiID + s.apiKey)
	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(user.Email),
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
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.senderEmail),
		// Uncomment to use a configuration set
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

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

	fmt.Println("SES Email Sent to " + user.Firstname + " " + user.Lastname + " at address: " + user.Email)
	fmt.Println(result)

	return nil
}
