package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-service/config"
	"payment-service/controller"
	"payment-service/messaging"
	"payment-service/repository"
	"payment-service/service"
	"payment-service/stripe"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize RabbitMQ client
	rabbitmqClient, err := messaging.NewRabbitMQClient()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmqClient.Close()
	log.Println("Connected to RabbitMQ")

	// Setup queues and exchanges
	if err := rabbitmqClient.SetupQueuesAndExchanges(); err != nil {
		log.Fatalf("Failed to setup queues and exchanges: %v", err)
	}
	log.Println("RabbitMQ queues and exchanges setup complete")

	// Initialize Stripe client
	stripeClient := stripe.NewStripeClient()
	log.Println("Stripe client initialized")

	// Initialize repository, service, and controller
	paymentRepository := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepository, stripeClient, rabbitmqClient)
	paymentController := controller.NewPaymentController(paymentService)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming order.created events from RabbitMQ
	err = rabbitmqClient.ConsumeOrderEvents(ctx, paymentService.ProcessOrderCreatedEvent)
	if err != nil {
		log.Fatalf("Failed to start order events consumer: %v", err)
	}
	log.Println("Started consuming order.created events")

	// Setup HTTP server with Gin
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "payment-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Payment API endpoints
	router.GET("/checkout/:orderId", paymentController.GetCheckoutURL)
	router.GET("/status/:orderId", paymentController.GetPaymentStatus)

	// Stripe webhook endpoint
	router.POST("/webhook/stripe", paymentController.HandleStripeWebhook)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8084",
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Payment Service starting on port 8084...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Payment Service...")

	// Cancel context to stop consumers
	cancel()

	// Graceful HTTP server shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Payment Service stopped")
}
