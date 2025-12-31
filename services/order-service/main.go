package main

import (
	"fmt"
	"net/http"
	"order-service/client"
	"order-service/config"
	"order-service/controller"
	"order-service/middleware"
	"order-service/repository"
	"order-service/service"

	"github.com/gin-gonic/gin"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		fmt.Errorf("failed to load environment variables: %v", err)
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	router := gin.Default()

	orderRepository := repository.NewOrderRepositoryImpl(db)
	orderService := service.NewOrderServiceImpl(orderRepository)
	foodClient := client.NewFoodClientImpl()
	orderController := controller.NewOrderController(orderService, foodClient)

	router.POST("/orders", middleware.AuthMiddleware(), orderController.CreateOrder)
	router.GET("/orders/:id", middleware.AuthMiddleware(), orderController.GetOrderById)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Order service is running"})
	})

	router.Run(":8083")
}
