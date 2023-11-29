package repository

import (
	"booking-engine/internal/model"
	"database/sql"
)

type FlightRepository struct {
	DB *sql.DB
}

type FlightPersister interface {
	SaveBooking(booking model.BookingRequest) (reservationID int, err error)
	GetAllReservations() ([]model.Reservation, error)
	GetBookingByID(bookingID int) (model.Reservation, error)
	UpdateInstanceID(reservationID int, instanceKey int64) error
}

// NewFlightRepository creates a new instance of FlightRepository
func NewFlightRepository(flight FlightRepository) FlightPersister {
	return &flight
}

// SaveBooking saves a new booking to the MySQL database
func (r *FlightRepository) SaveBooking(booking model.BookingRequest) (reservationID int, err error) {
	query := "INSERT INTO reservations (flight_number, passenger_id, seat_number, price, created_at) VALUES (?, ?, ?, ?, NOW())"
	result, err := r.DB.Exec(query, booking.FlightNumber, booking.PassengerID, booking.SeatNumber, booking.Price)

	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

// GetBookingByID retrieves a booking by ID from the MySQL database
func (r *FlightRepository) GetBookingByID(bookingID int) (model.Reservation, error) {
	query := "SELECT * FROM bookings WHERE id = ?"
	row := r.DB.QueryRow(query, bookingID)

	var booking model.Reservation
	err := row.Scan(&booking.ReservationID, &booking.FlightNumber, &booking.PassengerID, &booking.SeatNumber, &booking.Price, &booking.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return booking, model.ErrFlightNotFound
		}
		return booking, err
	}

	return booking, nil
}

// GetAllBookings retrieves all bookings from the MySQL database
func (r *FlightRepository) GetAllReservations() ([]model.Reservation, error) {
	query := "SELECT * FROM bookings"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Reservation
	for rows.Next() {
		var booking model.Reservation
		err := rows.Scan(&booking.ReservationID, &booking.FlightNumber, &booking.PassengerID, &booking.SeatNumber, &booking.Price, &booking.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// UpdateFlight updates a flight in the database
func (r *FlightRepository) UpdateInstanceID(reservationID int, instanceKey int64) error {
	_, err := r.DB.Exec("UPDATE reservations SET instance_key=? WHERE reservation_id=?",
		instanceKey, reservationID)
	if err != nil {
		return err
	}

	return nil
}
