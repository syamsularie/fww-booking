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
	UpdatePaymentStatus(paymentCode string, status bool) error
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

func (s *PaymentUsecase) UpdatePaymentStatus(paymentCode string, status bool) error {
	if err := s.PaymentRepo.UpdatePaymentStatus(paymentCode, status); err != nil {
		return err
	}

	reservationID, err := s.PaymentRepo.GetReservationIDByPaymentCode(paymentCode)
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
