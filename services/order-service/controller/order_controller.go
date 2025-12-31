package controller

import (
	"net/http"
	"order-service/client"
	"order-service/dto"
	"order-service/models"
	"order-service/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	orderService service.OrderService
	foodClient   client.FoodClient
}

func NewOrderController(orderService service.OrderService, foodClient client.FoodClient) *OrderController {
	return &OrderController{orderService: orderService, foodClient: foodClient}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var request dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate order ID
	orderID := uuid.New()

	// Convert DTO order items to model order items
	var orderItems []models.OrderItem
	var totalAmount float64
	for _, item := range request.OrderItems {
		food, err := c.foodClient.GetFoodById(item.FoodID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderItem := models.OrderItem{
			ID:       uuid.New(),
			OrderID:  orderID,
			FoodID:   item.FoodID,
			Quantity: item.Quantity,
			Price:    food.Price,
		}
		// TODO: Check stock before creating order
		orderItems = append(orderItems, orderItem)
		totalAmount += orderItem.Price * float64(orderItem.Quantity)
	}

	order := models.Order{
		ID:          orderID,
		UserID:      request.UserID,
		OrderItems:  orderItems,
		Status:      models.PENDING,
		TotalAmount: totalAmount,
	}

	if err := c.orderService.CreateOrder(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

func (c *OrderController) GetOrderById(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := c.orderService.GetOrderById(parsedID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
