package usecase

import (
	"booking-engine/internal/model"
	"booking-engine/internal/repository"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

type PaymentUsecase struct {
	PaymentRepo repository.PaymentPersister
}

type PaymentExecutor interface {
	GetPaymentDetailByPaymentID(paymentID int) (model.PaymentDetailResponse, error)
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
