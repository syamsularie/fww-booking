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

// @Summary Get Payment Detail by Payment ID
// @Description	Get Payment Detail by Payment ID
// @Tags payment
// @Accept json
// @Produce	json
// @Param id path string true "payment id"
// @Success 200 {object} model.PaymentDetailResponse "OK"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /payment/detail/{id} [get]
func (h *Payment) GetPaymentDetailByPaymentID(c *fiber.Ctx) error {
	paymentIDString := c.Params("id")

	paymentId, _ := strconv.Atoi(paymentIDString)
	paymentDetail, err := h.PaymentUsecase.GetPaymentDetailByPaymentID(paymentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(paymentDetail)
}

// @Summary Post Payment Pay
// @Description	Post Payment Pay
// @Tags payment
// @Accept json
// @Produce	json
// @Param payload body model.PaymentPayRequest true "payment pay request"
// @Success 200 {object} model.StatusResponse "OK"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /payment/pay [post]
func (h *Payment) PostPaymentPay(c *fiber.Ctx) error {
	var request model.PaymentPayRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.PaymentUsecase.UpdatePaymentStatus(request, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// @Summary Get Ticket Detail by Booking Code
// @Description	Get Ticket Detail by Booking Code
// @Tags payment
// @Accept json
// @Produce	json
// @Param booking_code path string true "booking code"
// @Success 200 {object} model.TicketDetailResponse "OK"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Router /ticket/detail/{booking_code} [get]
func (h *Payment) GetTicketDetailByBookingCode(c *fiber.Ctx) error {
	bookingCode := c.Params("booking_code")
	ticketDetail, err := h.PaymentUsecase.GetTicketDetailByBookingCode(bookingCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ticketDetail)
}
