package usecase

import (
	"booking-engine/internal/model"
	"booking-engine/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
)

// Service handles business logic for flights and bookings
type FlightUsecase struct {
	FlightRepo repository.FlightPersister
}

type FlightExecutor interface {
	GetFlightByID(id string) (*model.Flight, error)
	BookFlight(bookingRequest model.BookingRequest) (model.Reservation, error)
	GetAllReservations() ([]model.Reservation, error)
}

// NewService creates a new instance of the flight service
func NewFlightUsecaseService(flightUsecase *FlightUsecase) FlightExecutor {
	return flightUsecase
}

const ZeebeAddr = "0.0.0.0:26500"

// GetFlightByID returns details of a specific flight by ID
func (s *FlightUsecase) GetFlightByID(id string) (*model.Flight, error) {
	return s.GetFlightByID(id)
}

// BookFlight books a flight and returns the booking details
func (s *FlightUsecase) BookFlight(bookingRequest model.BookingRequest) (model.Reservation, error) {

	reservationId, err := s.FlightRepo.SaveBooking(bookingRequest)
	if err != nil {
		return model.Reservation{}, err
	}

	passengerIDRequest := strconv.Itoa(bookingRequest.PassengerID)
	fwwCoreApiURL := os.Getenv("FWW_CORE_URL") + "/passengers/" + passengerIDRequest
	fmt.Println(fwwCoreApiURL)
	response, err := http.Get(fwwCoreApiURL)
	if err != nil {
		fmt.Println(err)
		return model.Reservation{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return model.Reservation{}, err
	}
	var passenger model.PassengerResponse
	err = json.Unmarshal(body, &passenger)
	if err != nil {
		fmt.Println(err)
		return model.Reservation{}, err
	}

	fmt.Println(passenger)
	newBooking := model.Reservation{
		ReservationID: reservationId,
		FlightNumber:  bookingRequest.FlightNumber,
		PassengerID:   bookingRequest.PassengerID,
		SeatNumber:    bookingRequest.SeatNumber,
		Price:         bookingRequest.Price,
	}

	fmt.Println(reservationId, err)

	gatewayAddr := os.Getenv("ZEEBE_ADDRESS")
	plainText := false

	if gatewayAddr == "" {
		gatewayAddr = ZeebeAddr
		plainText = true
	}

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddr,
		UsePlaintextConnection: plainText,
	})

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	variables := model.BookingVariables{
		ReservationID:  reservationId,
		StatusPayment:  false,
		BlacklistUser:  false,
		Dukcapil:       "not valid",
		PeduliLindungi: "not vaksin",
		PassengerID:    passenger.Ktp,
	}

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("fww-reservation").LatestVersion().VariablesFromObject(variables)
	if err != nil {
		panic(err)
	}

	resp, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	err = s.FlightRepo.UpdateInstanceID(reservationId, resp.ProcessInstanceKey)
	return newBooking, nil
}

// GetBookings returns all flight bookings
func (s *FlightUsecase) GetAllReservations() ([]model.Reservation, error) {
	reservations, err := s.FlightRepo.GetAllReservations()
	if err != nil {
		return []model.Reservation{}, err
	}
	return reservations, nil
}

// Confirm Payment
func (s *FlightUsecase) ConfirmPayment(instanceID int64) error {

	return nil
}
