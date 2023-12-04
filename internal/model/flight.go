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

type FlightResponse struct {
	FlightNumber         string    `json:"flight_number"`
	DepartureAirportCode string    `json:"departure_airport_code"`
	ArrivalAirportCode   string    `json:"arrival_airport_code"`
	DepartureDateTime    time.Time `json:"departure_date_time"`
	ArrivalDateTime      time.Time `json:"arrival_date_time"`
}

// Booking represents a booking entity
type Reservation struct {
	ReservationID int       `json:"reservation_id"`
	FlightNumber  string    `json:"flight_number"`
	PassengerID   int       `json:"passenge_id"`
	SeatNumber    string    `json:"seat_number"`
	Price         float64   `json:"price"`
	PaymentCode   float64   `json:"payment_code"`
	BookingCode   string    `json:"booking_code"`
	CreatedAt     time.Time `json:"create_at"`
}

// BookingRequest represents the request structure for booking a flight
type BookingRequest struct {
	FlightNumber string  `json:"flight_number"`
	PassengerID  int     `json:"passenger_id"`
	SeatNumber   string  `json:"seat_number"`
	Price        float64 `json:"price"`
}

type EmailRequest struct {
	ReservationId int `json:"reservation_id"`
}

// Variable BPMN

type BookingVariables struct {
	ReservationID  int    `json:"reservationId"`
	BlacklistUser  bool   `json:"blacklistUser"`
	PeduliLindungi string `json:"peduliLindungi"`
	Dukcapil       string `json:"dukcapil"`
	PassengerID    string `json:"passengerId"`
	StatusPayment  bool   `json:"status_payment"`
}

var (
	ErrFlightNotFound = errors.New("flight not found")
)
