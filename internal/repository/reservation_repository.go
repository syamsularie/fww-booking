package repository

import (
	"booking-engine/internal/model"
	"database/sql"
)

type ReservationRepository struct {
	DB *sql.DB
}

type ReservationPersister interface {
	UpdateBookingCode(reservationID int, bookingCode string) error
	GetReservationByBookingCode(bookingCode string) (model.Reservation, error)
	GetReservationById(reservationId int) (model.Reservation, error)
}

func NewReservationRepository(reservation ReservationRepository) ReservationPersister {
	return &reservation
}

func (r *ReservationRepository) UpdateBookingCode(reservationID int, bookingCode string) error {
	query := `UPDATE reservations SET booking_code = ? WHERE reservation_id = ?`
	_, err := r.DB.Exec(query, bookingCode, reservationID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReservationRepository) GetReservationByBookingCode(bookingCode string) (model.Reservation, error) {
	var reservation model.Reservation

	query := `SELECT reservation_id, flight_number, passenger_id, seat_number, price, created_at FROM reservations WHERE booking_code = ?`
	row := r.DB.QueryRow(query, bookingCode)
	err := row.Scan(&reservation.ReservationID, &reservation.FlightNumber, &reservation.PassengerID, &reservation.SeatNumber, &reservation.Price, &reservation.CreatedAt)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

func (r *ReservationRepository) GetReservationById(reservationId int) (model.Reservation, error) {
	var reservation model.Reservation

	query := `SELECT reservation_id, flight_number, passenger_id, seat_number, price, created_at FROM reservations WHERE reservation_id = ?`
	row := r.DB.QueryRow(query, reservationId)
	err := row.Scan(&reservation.ReservationID, &reservation.FlightNumber, &reservation.PassengerID, &reservation.SeatNumber, &reservation.Price, &reservation.CreatedAt)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}
