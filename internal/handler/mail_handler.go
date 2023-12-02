package handler

import (
	"net/http"

	"github.com/gofiber/fiber"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	awsRegion = "ap-southeast-1"
	sender    = "syams.arie@gmail.com"
	recipient = "syams.arie@gmail.com"
)

func sendEmail(c *fiber.Ctx) error {
	// Create a new AWS session using credentials from environment variables, IAM role, or AWS credentials file.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create AWS session"})
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
				Text: &ses.Content{
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}
