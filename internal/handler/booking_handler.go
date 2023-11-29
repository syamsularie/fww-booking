package handler

import (
	"booking-engine/internal/model"
	"booking-engine/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for flights and bookings
type Handler struct {
	Usecase usecase.FlightExecutor
}
type FlightaHandler interface {
	GetFlightByID(c *fiber.Ctx) error
	BookFlight(c *fiber.Ctx) error
	GetAllReservations(c *fiber.Ctx) error
}

// NewHandler creates a new instance of the flight handler
func NewHandler(handler Handler) FlightaHandler {
	return &handler
}

// GetFlightByID handles the GET /flights/:id endpoint
func (h *Handler) GetFlightByID(c *fiber.Ctx) error {
	id := c.Params("id")
	flight, err := h.Usecase.GetFlightByID(id)
	if err != nil {
		if err == model.ErrFlightNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Flight not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(flight)
}

// BookFlight handles the POST /bookings endpoint
func (h *Handler) BookFlight(c *fiber.Ctx) error {
	var request model.BookingRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request format")
	}

	booking, err := h.Usecase.BookFlight(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.JSON(booking)
}

// GetBookings handles the GET /bookings endpoint
func (h *Handler) GetAllReservations(c *fiber.Ctx) error {
	var reservations []model.Reservation
	reservations, err := h.Usecase.GetAllReservations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(reservations)
}
