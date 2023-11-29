package model

import (
	"errors"
	"time"
)

// Flight represents a flight entity
type Flight struct {
	FlightNumber   string    `json:"flight_number"`
	Departure      string    `json:"departure"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departure_time"`
	Price          float64   `json:"price"`
	AvailableSeats int       `json:"available_seats"`
}

// Booking represents a booking entity
type Reservation struct {
	ReservationID int       `json:"reservation_id"`
	FlightNumber  string    `json:"flight_number"`
	PassengerID   int       `json:"passenge_idr"`
	SeatNumber    string    `json:"seat_number"`
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"create_at"`
}

// BookingRequest represents the request structure for booking a flight
type BookingRequest struct {
	FlightNumber string  `json:"flight_number"`
	PassengerID  int     `json:"passenger_id"`
	SeatNumber   string  `json:"seat_number"`
	Price        float64 `json:"price"`
}

// Variable BPMN

type BookingVariables struct {
	ReservationID int  `json:"reservation_id"`
	StatusPayment bool `json:"status_payment"`
}

var (
	ErrFlightNotFound = errors.New("flight not found")
)
