package usecase

import (
	"booking-engine/internal/model"
	"booking-engine/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type EmailUsecase struct {
	PaymentRepo     repository.PaymentPersister
	ReservationRepo repository.ReservationPersister
}

type EmailExecutor interface {
	SendEmail(reservationId int) error
	SendEmailUnpaid(reservationId int) error
	SendEmailFailedPayment(reservationId int) error
	SendEmailReservationCode(reservationId int) error
}

func NewEmailUsecaseService(emailUsecase *EmailUsecase) EmailExecutor {
	return emailUsecase
}

func (e *EmailUsecase) SendEmail(reservationId int) error {

	reservation, err := e.ReservationRepo.GetReservationById(reservationId)
	if err != nil {
		return err
	}

	var paymentDetailResponse model.PaymentDetailResponse

	paymentDetail, err := e.PaymentRepo.GetPaymentDetailByReservationID(reservation.ReservationID)
	if err != nil {
		return err
	}

	passengerIDRequest := strconv.Itoa(paymentDetail.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return err
	}

	paymentDetailResponse.FlightNumber = paymentDetail.FlightNumber
	paymentDetailResponse.PassengerFirstName = passenger.FirstName
	paymentDetailResponse.PassengerLastName = passenger.LastName
	paymentDetailResponse.SeatNumber = paymentDetail.SeatNumber
	paymentDetailResponse.Price = paymentDetail.Price
	paymentDetailResponse.PaymentStatus = paymentDetail.PaymentStatus
	paymentDetailResponse.PaymentMethod = paymentDetail.PaymentMethod
	paymentDetailResponse.PaymentCode = paymentDetail.PaymentCode

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return err
	}
	// Parse the HTML template
	tmpl, err := template.New("reservationTemplate").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template to generate HTML content
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, paymentDetailResponse); err != nil {
		log.Fatal(err)
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
					Data: aws.String(emailBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Reservation Ticket"),
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

func (e *EmailUsecase) SendEmailUnpaid(reservationId int) error {

	reservation, err := e.ReservationRepo.GetReservationById(reservationId)
	if err != nil {
		return err
	}

	var paymentDetailResponse model.PaymentDetailResponse

	paymentDetail, err := e.PaymentRepo.GetPaymentDetailByReservationID(reservation.ReservationID)
	if err != nil {
		return err
	}

	passengerIDRequest := strconv.Itoa(paymentDetail.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return err
	}

	paymentDetailResponse.FlightNumber = paymentDetail.FlightNumber
	paymentDetailResponse.PassengerFirstName = passenger.FirstName
	paymentDetailResponse.PassengerLastName = passenger.LastName
	paymentDetailResponse.SeatNumber = paymentDetail.SeatNumber
	paymentDetailResponse.Price = paymentDetail.Price
	paymentDetailResponse.PaymentStatus = paymentDetail.PaymentStatus
	paymentDetailResponse.PaymentMethod = paymentDetail.PaymentMethod
	paymentDetailResponse.PaymentCode = paymentDetail.PaymentCode

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return err
	}

	// Parse the HTML template
	tmpl, err := template.New("reservationTemplateUnpaid").Parse(htmlTemplateUnpaid)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template to generate HTML content
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, paymentDetailResponse); err != nil {
		log.Fatal(err)
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
					Data: aws.String(emailBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Reservation Ticket Unpaid"),
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

func (e *EmailUsecase) SendEmailFailedPayment(reservationId int) error {

	reservation, err := e.ReservationRepo.GetReservationById(reservationId)
	if err != nil {
		return err
	}

	var paymentDetailResponse model.PaymentDetailResponse

	paymentDetail, err := e.PaymentRepo.GetPaymentDetailByReservationID(reservation.ReservationID)
	if err != nil {
		return err
	}

	passengerIDRequest := strconv.Itoa(paymentDetail.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return err
	}

	paymentDetailResponse.FlightNumber = paymentDetail.FlightNumber
	paymentDetailResponse.PassengerFirstName = passenger.FirstName
	paymentDetailResponse.PassengerLastName = passenger.LastName
	paymentDetailResponse.SeatNumber = paymentDetail.SeatNumber
	paymentDetailResponse.Price = paymentDetail.Price
	paymentDetailResponse.PaymentStatus = paymentDetail.PaymentStatus
	paymentDetailResponse.PaymentMethod = paymentDetail.PaymentMethod
	paymentDetailResponse.PaymentCode = paymentDetail.PaymentCode
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return err
	}

	// Parse the HTML template
	tmpl, err := template.New("reservationTemplateFailedPayment").Parse(htmlTemplateFailedPayment)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template to generate HTML content
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, paymentDetailResponse); err != nil {
		log.Fatal(err)
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
					Data: aws.String(emailBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Failed Payment Reservation Ticket"),
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

func (e *EmailUsecase) SendEmailReservationCode(reservationId int) error {
	var ticketDetailResponse model.TicketDetailResponse
	reservation, err := e.ReservationRepo.GetReservationById(reservationId)
	if err != nil {
		return err
	}

	fmt.Println("BLA", reservation)
	//Fetch passenger
	passengerIDRequest := strconv.Itoa(reservation.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return err
	}

	//Fetch flight
	fwwCoreApiURL = os.Getenv("FWW_CORE_URL") + "/flights/" + reservation.FlightNumber
	response, err = http.Get(fwwCoreApiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var flight model.FlightResponse
	err = json.Unmarshal(body, &flight)
	if err != nil {
		return err
	}

	fmt.Println("BLA", flight)
	ticketDetailResponse.FlightNumber = reservation.FlightNumber
	ticketDetailResponse.BookingCode = reservation.BookingCode
	ticketDetailResponse.PassengerFirstName = passenger.FirstName
	ticketDetailResponse.PassengerLastName = passenger.LastName
	ticketDetailResponse.FlightNumber = flight.FlightNumber
	ticketDetailResponse.SeatNumber = reservation.SeatNumber
	ticketDetailResponse.DepartureAirportCode = flight.DepartureAirportCode
	ticketDetailResponse.ArrivalAirportCode = flight.ArrivalAirportCode
	ticketDetailResponse.DepartureDateTime = flight.DepartureDateTime
	ticketDetailResponse.ArrivalDateTime = flight.ArrivalDateTime

	fmt.Println("BLA", ticketDetailResponse)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return err
	}
	// Parse the HTML template
	tmpl, err := template.New("reservationCodeTemplate").Parse(htmlBookingTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template to generate HTML content
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, ticketDetailResponse); err != nil {
		log.Fatal(err)
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
					Data: aws.String(emailBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Booking Ticket Confirmation"),
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

const (
	awsRegion    = "ap-southeast-1"
	sender       = "syams.arie@gmail.com"
	recipient    = "syams.arie@gmail.com"
	htmlTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Ticket Reservation Confirmation</title>
	</head>
	<body>
	    <div>
	        <h1>Ticket Reservation Confirmation</h1>
	        <p>Your ticket reservation has been created. Below are the details of your reservation:</p>
	        <p><strong>Flight Number:</strong> {{.FlightNumber}}</p>
	        <p><strong>Passenger Name:</strong> {{.PassengerFirstName}} {{.PassengerLastName}}</p>
	        <p><strong>Seat Number:</strong> {{.SeatNumber}}</p>
	        <p><strong>Price</strong> {{.Price}}</p>
	        <p><strong>Payment Code</strong> {{.PaymentCode}}</p>
	        <p>Safe travels!</p>
	    </div>
	</body>
	</html>`

	htmlBookingTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ticket Reservation Confirmation</title>
</head>
<body>
    <div>
        <h1>Ticket Reservation Confirmation</h1>
        <p>Dear {{.PassengerFirstName}} {{.PassengerLastName}},</p>
        <p>Your ticket reservation has been confirmed. Below are the details of your reservation:</p>
        <p><strong>Reservation Code:</strong> {{.BookingCode}}</p>
        <p><strong>Seat Number:</strong> {{.SeatNumber}}</p>
        <p><strong>Departure Airport:</strong> {{.DepartureAirportCode}}</p>
        <p><strong>Arrival Airport:</strong> {{.ArrivalAirportCode}}</p>
        <p><strong>Departure Date:</strong> {{.DepartureDateTime}}</p>
        <p>We look forward to having you on board. Please ensure you arrive at the airport well in advance.</p>
        <p>Safe travels!</p>
    </div>
</body>
</html>
`

	htmlTemplateUnpaid = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Ticket Reservation Unpaid</title>
</head>
<body>
	<div>
		<h1>Ticket Reservation Unpaid</h1>
		<p>Your ticket reservation has been failed. Below are the details of your reservation:</p>
		<p><strong>Flight Number:</strong> {{.FlightNumber}}</p>
		<p><strong>Passenger Name:</strong> {{.PassengerFirstName}} {{.PassengerLastName}}</p>
		<p><strong>Seat Number:</strong> {{.SeatNumber}}</p>
		<p><strong>Price</strong> {{.Price}}</p>
		<p><strong>Payment Code</strong> {{.PaymentCode}}</p>
		<p>Safe travels!</p>
	</div>
</body>
</html>`

	htmlTemplateFailedPayment = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Failed Payment Ticket Reservation</title>
</head>
<body>
	<div>
		<h1>Failed Payment Ticket Reservation</h1>
		<p>Your ticket reservation has been failed. Below are the details of your reservation:</p>
		<p><strong>Flight Number:</strong> {{.FlightNumber}}</p>
		<p><strong>Passenger Name:</strong> {{.PassengerFirstName}} {{.PassengerLastName}}</p>
		<p><strong>Seat Number:</strong> {{.SeatNumber}}</p>
		<p><strong>Price</strong> {{.Price}}</p>
		<p><strong>Payment Code</strong> {{.PaymentCode}}</p>
		<p>Safe travels!</p>
	</div>
</body>
</html>`
)
