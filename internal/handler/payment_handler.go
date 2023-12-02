package handler

import (
	"booking-engine/internal/model"
	"booking-engine/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Payment struct {
	PaymentUsecase usecase.PaymentExecutor
}

type PaymentHandler interface {
	GetPaymentDetailByPaymentID(c *fiber.Ctx) error
	PostPaymentPay(c *fiber.Ctx) error
	GetTicketDetailByBookingCode(c *fiber.Ctx) error
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

func (h *Payment) PostPaymentPay(c *fiber.Ctx) error {
	var request model.PaymentPayRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.PaymentUsecase.UpdatePaymentStatus(request.PaymentCode, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *Payment) GetTicketDetailByBookingCode(c *fiber.Ctx) error {
	bookingCode := c.Params("booking_code")
	ticketDetail, err := h.PaymentUsecase.GetTicketDetailByBookingCode(bookingCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ticketDetail)
}
