package main

import (
	"booking-engine/config"
	"booking-engine/config/middleware"
	"booking-engine/internal/handler"
	"booking-engine/internal/repository"
	"booking-engine/internal/usecase"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "booking-engine/docs"

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

// @title FWW-Booking API
// @version 1.0
// @description This is a FWW Booking API
// @host localhost:3002
// @BasePath /
func main() {
	baseDep := config.NewBaseDep()
	loadEnv(baseDep.Logger)
	db, err := config.NewDbPool(baseDep.Logger)

	if err != nil {
		os.Exit(1)
	}

	dbCollector := middleware.NewStatsCollector("fww", db)
	prometheus.MustRegister(dbCollector)
	fiberProm := middleware.NewWithRegistry(prometheus.DefaultRegisterer, "fww-booking", "", "", map[string]string{})

	// Initialize repository
	flightRepo := repository.NewFlightRepository(repository.FlightRepository{
		DB: db,
	})

	paymentRepo := repository.NewPaymentRepository(repository.PaymentRepository{
		DB: db,
	})

	reservationRepo := repository.NewReservationRepository(repository.ReservationRepository{
		DB: db,
	})

	// Initialize usecase
	flightUscase := usecase.NewFlightUsecaseService(&usecase.FlightUsecase{
		FlightRepo: flightRepo,
	})

	paymentUsecase := usecase.NewPaymentUsecaseService(&usecase.PaymentUsecase{
		PaymentRepo:     paymentRepo,
		ReservationRepo: reservationRepo,
	})

	emailUsecase := usecase.NewEmailUsecaseService(&usecase.EmailUsecase{
		PaymentRepo:     paymentRepo,
		ReservationRepo: reservationRepo,
	})

	// Initialize handler
	flightHandler := handler.NewHandler(handler.Handler{
		Usecase: flightUscase,
	})

	paymentHandler := handler.NewPaymentHandler(handler.Payment{
		PaymentUsecase: paymentUsecase,
	})

	emailHandler := handler.NewEmailHandler(handler.Email{
		EmailUsecase: emailUsecase,
	})

	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})

	app.Use(fiberProm.Middleware)
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))
	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	//== Send email payment
	app.Post("/send-email", emailHandler.SendEmail)
	//=== Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)
	//=== healthz route
	app.Get("/healthz", Healthz)
	//=== reservation route
	app.Get("/flights/:id", flightHandler.GetFlightByID)
	app.Post("/bookings", flightHandler.BookFlight)
	app.Get("/bookings", flightHandler.GetAllReservations)
	app.Post("/complete", Complete)
	//=== payment route
	app.Get("/payment/detail/:id", paymentHandler.GetPaymentDetailByPaymentID)
	app.Get("/payment/detail/reservation/:id", paymentHandler.GetPaymentDetailByReservationID)
	app.Post("/payment/pay", paymentHandler.PostPaymentPay)
	app.Get("/ticket/detail/:booking_code", paymentHandler.GetTicketDetailByBookingCode)

	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", "3002")); err != nil {
		log.Fatal(err)
	}

}

func Healthz(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Service is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

func Complete(c *fiber.Ctx) error {
	zeebeBrokerURL := os.Getenv("ZEEBE_ADDRESS")
	// workflowInstanceKey := "2251799814100631"
	// taskID := 2251799814100650

	// Send a POST request to complete the user task
	url := fmt.Sprintf("https://%s/workflow-instance/2251799814100631/complete", zeebeBrokerURL)
	requestBody := fmt.Sprintf(`{"elementId": "%s"}`, "2251799814100650")

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("User task completed successfully.")
	} else {
		fmt.Printf("Failed to complete user task. Status code: %d\n", resp.StatusCode)
	}

	if err := c.JSON("OK"); err != nil {
		return err
	}

	return nil
}

func loadEnv(logger config.Logger) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load()
		if err != nil {
			logger.Error("no .env files provided")
		}
	}
}
