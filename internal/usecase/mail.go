package usecase

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type EmailUsecase struct {
}

type EmailExecutor interface {
	SendEmail(reservationId int) error
}

func NewEmailUsecaseService(emailUsecase *EmailUsecase) EmailExecutor {
	return emailUsecase
}

const (
	awsRegion = "ap-southeast-1"
	sender    = "syams.arie@gmail.com"
	recipient = "syams.arie@gmail.com"
)

func (e *EmailUsecase) SendEmail(reservationId int) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return err
	}

	// Create an SES client
	sesClient := ses.New(sess)

	// Construct the email input
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String("Hello, this is the body of the email."),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Test Email"),
			},
		},
		Source: aws.String(sender),
	}

	// Send the email
	_, err = sesClient.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}
