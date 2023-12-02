package model

import "time"

type PaymentDetail struct {
	PassengerID   int     `json:"passenger_id"`
	FlightNumber  string  `json:"flight_number"`
	SeatNumber    int     `json:"seat_number"`
	Price         float32 `json:"price"`
	PaymentStatus bool    `json:"payment_status"`
	PaymentMethod string  `json:"payment_method"`
	PaymentCode   string  `json:"payment_code"`
}

type PaymentDetailResponse struct {
	FlightNumber       string  `json:"flight_number"`
	PassengerFirstName string  `json:"passenger_first_name"`
	PassengerLastName  string  `json:"passenger_last_name"`
	SeatNumber         int     `json:"seat_number"`
	Price              float32 `json:"price"`
	PaymentStatus      bool    `json:"payment_status"`
	PaymentMethod      string  `json:"payment_method"`
	PaymentCode        string  `json:"payment_code"`
}

type PaymentPayRequest struct {
	PaymentCode string `json:"payment_code"`
}

type TicketDetailResponse struct {
	FlightNumber         string    `json:"flight_number"`
	BookingCode          string    `json:"booking_code"`
	PassengerFirstName   string    `json:"passenger_first_name"`
	PassengerLastName    string    `json:"passenger_last_name"`
	SeatNumber           string    `json:"seat_number"`
	DepartureAirportCode string    `json:"departure_airport_code"`
	ArrivalAirportCode   string    `json:"arrival_airport_code"`
	DepartureDateTime    time.Time `json:"departure_date_time"`
	ArrivalDateTime      time.Time `json:"arrival_date_time"`
}
