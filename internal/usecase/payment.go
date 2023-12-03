package usecase

import (
	"booking-engine/internal/model"
	"booking-engine/internal/repository"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PaymentUsecase struct {
	PaymentRepo     repository.PaymentPersister
	ReservationRepo repository.ReservationPersister
}

type PaymentExecutor interface {
	GetPaymentDetailByPaymentID(paymentID int) (model.PaymentDetailResponse, error)
	UpdatePaymentStatus(payment model.PaymentPayRequest, status bool) error
	GetTicketDetailByBookingCode(bookingCode string) (model.TicketDetailResponse, error)
}

func NewPaymentUsecaseService(paymentUsecase *PaymentUsecase) PaymentExecutor {
	return paymentUsecase
}

func (s *PaymentUsecase) GetPaymentDetailByPaymentID(paymentID int) (model.PaymentDetailResponse, error) {
	var paymentDetailResponse model.PaymentDetailResponse

	paymentDetail, err := s.PaymentRepo.GetPaymentDetailByPaymentID(paymentID)
	if err != nil {
		return paymentDetailResponse, err
	}

	passengerIDRequest := strconv.Itoa(paymentDetail.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return paymentDetailResponse, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return paymentDetailResponse, err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return paymentDetailResponse, err
	}

	paymentDetailResponse.FlightNumber = paymentDetail.FlightNumber
	paymentDetailResponse.PassengerFirstName = passenger.FirstName
	paymentDetailResponse.PassengerLastName = passenger.LastName
	paymentDetailResponse.SeatNumber = paymentDetail.SeatNumber
	paymentDetailResponse.Price = paymentDetail.Price
	paymentDetailResponse.PaymentStatus = paymentDetail.PaymentStatus
	paymentDetailResponse.PaymentMethod = paymentDetail.PaymentMethod
	paymentDetailResponse.PaymentCode = paymentDetail.PaymentCode

	return paymentDetailResponse, nil
}

func (s *PaymentUsecase) UpdatePaymentStatus(payment model.PaymentPayRequest, status bool) error {
	if err := s.PaymentRepo.UpdatePaymentStatus(payment, status); err != nil {
		return err
	}

	reservationID, err := s.PaymentRepo.GetReservationIDByPaymentCode(payment.PaymentCode)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bookCode := make([]byte, 5)
	for i := range bookCode {
		bookCode[i] = charset[rand.Intn(len(charset))]
	}

	err = s.ReservationRepo.UpdateBookingCode(reservationID, string(bookCode))
	if err != nil {
		return err
	}

	return nil
}

func (s *PaymentUsecase) GetTicketDetailByBookingCode(bookingCode string) (model.TicketDetailResponse, error) {
	var ticketDetailResponse model.TicketDetailResponse

	reservation, err := s.ReservationRepo.GetReservationByBookingCode(bookingCode)
	if err != nil {
		return ticketDetailResponse, err
	}

	//Fetch passenger
	passengerIDRequest := strconv.Itoa(reservation.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		return ticketDetailResponse, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ticketDetailResponse, err
	}

	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		return ticketDetailResponse, err
	}

	//Fetch flight
	fwwCoreApiURL = os.Getenv("FWW_CORE_URL") + "/flights/" + reservation.FlightNumber
	response, err = http.Get(fwwCoreApiURL)
	if err != nil {
		return ticketDetailResponse, err
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return ticketDetailResponse, err
	}

	var flight model.FlightResponse
	err = json.Unmarshal(body, &flight)
	if err != nil {
		return ticketDetailResponse, err
	}

	ticketDetailResponse.FlightNumber = reservation.FlightNumber
	ticketDetailResponse.BookingCode = bookingCode
	ticketDetailResponse.PassengerFirstName = passenger.FirstName
	ticketDetailResponse.PassengerLastName = passenger.LastName
	ticketDetailResponse.FlightNumber = flight.FlightNumber
	ticketDetailResponse.SeatNumber = reservation.SeatNumber
	ticketDetailResponse.DepartureAirportCode = flight.DepartureAirportCode
	ticketDetailResponse.ArrivalAirportCode = flight.ArrivalAirportCode
	ticketDetailResponse.DepartureDateTime = flight.DepartureDateTime
	ticketDetailResponse.ArrivalDateTime = flight.ArrivalDateTime

	return ticketDetailResponse, nil
}
