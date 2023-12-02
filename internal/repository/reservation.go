package repository

import "database/sql"

type ReservationRepository struct {
	DB *sql.DB
}

type ReservationPersister interface {
	UpdateBookingCode(reservationID int, bookingCode string) error
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
