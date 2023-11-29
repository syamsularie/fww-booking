package usecase

import (
	"booking-engine/internal/model"
	"booking-engine/internal/repository"
	"context"
	"fmt"
	"os"

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
	// For simplicity, let's assume the booking is successful
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
	// variables := make(map[model.BookingVariables]interface{})
	variables := model.BookingVariables{
		ReservationID: reservationId,
		StatusPayment: false,
	}

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("fww-bpm").LatestVersion().VariablesFromObject(variables)
	if err != nil {
		panic(err)
	}

	resp, err := request.Send(ctx)
	fmt.Println("BLI", resp.ProcessInstanceKey)
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
