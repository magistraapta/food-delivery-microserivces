package main

import (
	"context"
	"log"
	"net/http"
	"order-service/client"
	"order-service/config"
	"order-service/controller"
	"order-service/messaging"
	"order-service/middleware"
	"order-service/repository"
	"order-service/service"
	"os"
	"os/signal"
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

	// Initialize repository and service
	orderRepository := repository.NewOrderRepositoryImpl(db)
	orderService := service.NewOrderServiceImpl(orderRepository, rabbitmqClient)
	foodClient := client.NewFoodClientImpl()
	orderController := controller.NewOrderController(orderService, foodClient)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming payment events from RabbitMQ
	err = rabbitmqClient.ConsumePaymentEvents(
		ctx,
		orderService.ProcessPaymentSuccess,
		orderService.ProcessPaymentFailed,
	)
	if err != nil {
		log.Fatalf("Failed to start payment events consumer: %v", err)
	}
	log.Println("Started consuming payment events")

	// Start consuming payment timeout events from RabbitMQ
	// These events arrive after a 5-minute delay to check if payment was completed
	err = rabbitmqClient.ConsumePaymentTimeoutEvents(ctx, orderService.ProcessPaymentTimeout)
	if err != nil {
		log.Fatalf("Failed to start payment timeout consumer: %v", err)
	}
	log.Println("Started consuming payment timeout events")

	// Setup HTTP server with Gin
	router := gin.Default()

	// Order routes
	router.POST("/orders", middleware.AuthMiddleware(), orderController.CreateOrder)
	router.GET("/orders/:id", middleware.AuthMiddleware(), orderController.GetOrderById)

	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "order-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8083",
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Order Service starting on port 8083...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Order Service...")

	// Cancel context to stop consumers
	cancel()

	// Graceful HTTP server shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Order Service stopped")
}
