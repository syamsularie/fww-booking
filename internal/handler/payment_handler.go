package handler

import (
	"booking-engine/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Payment struct {
	PaymentUsecase usecase.PaymentExecutor
}

type PaymentHandler interface {
	GetPaymentDetailByPaymentID(c *fiber.Ctx) error
}

func NewPaymentHandler(handler Payment) PaymentHandler {
	return &handler
}

func (h *Payment) GetPaymentDetailByPaymentID(c *fiber.Ctx) error {
	paymentIDString := c.Params("id")

	paymentId, _ := strconv.Atoi(paymentIDString)
	paymentDetail, err := h.PaymentUsecase.GetPaymentDetailByPaymentID(paymentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(paymentDetail)
}
