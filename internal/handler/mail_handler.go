package handler

import (
	"booking-engine/internal/model"
	"booking-engine/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

const (
	awsRegion = "ap-southeast-1"
	sender    = "syams.arie@gmail.com"
	recipient = "syams.arie@gmail.com"
)

type Email struct {
	EmailUsecase usecase.EmailExecutor
}

type EmailHandler interface {
	SendEmail(c *fiber.Ctx) error
}

func NewEmailHandler(email Email) EmailHandler {
	return &email
}

// sendEmail implements EmailHandler.
func (e *Email) SendEmail(c *fiber.Ctx) error {
	var request model.EmailRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request format")
	}

	if err := e.EmailUsecase.SendEmail(request.ReservationId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email sent successfully"})
	// return c.JSON(fiber.Map{"message": "Email sent successfully"})
}
