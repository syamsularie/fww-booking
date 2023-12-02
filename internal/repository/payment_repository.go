package repository

import (
	"booking-engine/internal/model"
	"database/sql"
)

type PaymentRepository struct {
	DB *sql.DB
}

type PaymentPersister interface {
	GetPaymentDetailByPaymentID(paymentID int) (model.PaymentDetail, error)
}

func NewPaymentRepository(payment PaymentRepository) PaymentPersister {
	return &payment
}

func (r *PaymentRepository) GetPaymentDetailByPaymentID(paymentID int) (model.PaymentDetail, error) {
	var paymentDetail model.PaymentDetail
	query := `
				SELECT r.passenger_id, r.flight_number, r.seat_number, r.price, p.payment_status, p.payment_method, p.payment_code
				FROM payments p 
				INNER JOIN reservations r on p.reservation_id = r.reservation_id
				WHERE p.payment_id = ?
			`
	row := r.DB.QueryRow(query, paymentID)
	err := row.Scan(&paymentDetail.PassengerID, &paymentDetail.FlightNumber, &paymentDetail.SeatNumber, &paymentDetail.Price, &paymentDetail.PaymentStatus, &paymentDetail.PaymentMethod, &paymentDetail.PaymentCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return paymentDetail, nil
		}

		return paymentDetail, err
	}

	return paymentDetail, nil

}
