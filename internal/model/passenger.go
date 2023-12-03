package model

type PassengerResponse struct {
	PassengerID int    `json:"passenger_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Ktp         string `json:"ktp"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}
