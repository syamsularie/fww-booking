package handler

import (
	"booking-engine/internal/usecase"

	"github.com/gofiber/fiber"
)

const (
	awsRegion = "ap-southeast-1"
	sender    = "syams.arie@gmail.com"
	recipient = "syams.arie@gmail.com"
)

type Email struct {
	EmailUsecase usecase.EmailExecutor
}

type EmailHandler interface {
	SendEmail(c *fiber.Ctx) error
}

func NewEmailHandler(email Email) EmailHandler {
	return &email
}

// func sendEmail(c *fiber.Ctx) error {

// 	// Replace these values with your SMTP server details
// 	smtpHost := "email-smtp.ap-southeast-1.amazonaws.com"
// 	smtpPort := "465"
// 	smtpUsername := os.Getenv("SMTP_USER")
// 	smtpPassword := os.Getenv("SMTP_PASSWORD")
// 	sender := "syams.arie@gmail.com"
// 	recipient := "syams.arie@gmail.com"

// 	// Payment code details
// 	paymentCode := "ABC123456"
// 	recipientName := "John Doe" // Replace with the actual recipient's name

// 	// Load HTML template
// 	templateData := EmailData{
// 		RecipientName: recipientName,
// 		PaymentCode:   paymentCode,
// 	}

// 	// Compose the email message
// 	subject := "Reservation Ticket"
// 	// body := "Hello, this is the body of the email."
// 	body := `
// 	<!DOCTYPE html>
// 	<html lang="en">
// 	<head>
// 		<meta charset="UTF-8">
// 		<meta http-equiv="X-UA-Compatible" content="IE=edge">
// 		<meta name="viewport" content="width=device-width, initial-scale=1.0">
// 		<title>Payment Code Email</title>
// 	</head>
// 	<body>
// 		<h1>Payment Code Email</h1>
// 		<p>Hello ` + templateData.RecipientName + `,</p>
// 		<p>Your payment code is: <strong>{{.PaymentCode}}</strong></p>
// 		<p>Thank you for your payment!</p>
// 	</body>
// 	</html>
// `
// 	message := "Subject: " + subject + "\r\n\r\n" + body

// 	// Set up authentication information
// 	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

// 	// Connect to the SMTP server
// 	smtpAddr := smtpHost + ":" + smtpPort
// 	err := smtp.SendMail(smtpAddr, auth, sender, []string{recipient}, []byte(message))
// 	if err != nil {
// 		log.Fatal("Failed to send email:", err)
// 	}

// 	log.Println("Email sent successfully.")

// 	return c.JSON(fiber.Map{"message": "Email sent successfully"})

// }

// sendEmail implements EmailHandler.
func (e *Email) SendEmail(c *fiber.Ctx) error {
	// Create a new AWS session using credentials from environment variables, IAM role, or AWS credentials file.
	if err := e.EmailUsecase.SendEmail(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email sent successfully"})
	// return c.JSON(fiber.Map{"message": "Email sent successfully"})
}
