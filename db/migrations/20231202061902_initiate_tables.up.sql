CREATE TABLE reservations (
    reservation_id INT AUTO_INCREMENT PRIMARY KEY,
    flight_number VARCHAR(10),
    passenger_id INT,
    seat_number INT,
    price DECIMAL(10, 2) NOT NULL,
    instance_key BIGINT,
    booking_code VARCHAR(10),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE payments (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50),
    payment_code VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE fww_booking.payments ADD CONSTRAINT payments_FK FOREIGN KEY (reservation_id) REFERENCES fww_booking.reservations(reservation_id);
