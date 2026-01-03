package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"payment-service/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type PaymentController struct {
	paymentService *service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// GetCheckoutURL returns the Stripe Checkout URL for an order
// GET /checkout/:orderId
func (c *PaymentController) GetCheckoutURL(ctx *gin.Context) {
	orderIDStr := ctx.Param("orderId")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := c.paymentService.GetPaymentByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.CheckoutURL == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Checkout URL not yet available"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order_id":     payment.OrderID,
		"checkout_url": payment.CheckoutURL,
		"status":       payment.Status,
		"amount":       payment.Amount,
		"currency":     payment.Currency,
	})
}

// GetPaymentStatus returns the payment status for an order
// GET /status/:orderId
func (c *PaymentController) GetPaymentStatus(ctx *gin.Context) {
	orderIDStr := ctx.Param("orderId")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := c.paymentService.GetPaymentByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order_id":   payment.OrderID,
		"status":     payment.Status,
		"amount":     payment.Amount,
		"currency":   payment.Currency,
		"created_at": payment.CreatedAt,
		"updated_at": payment.UpdatedAt,
	})
}

// HandleStripeWebhook processes Stripe webhook events
// POST /webhook/stripe
func (c *PaymentController) HandleStripeWebhook(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading body"})
		return
	}

	// Get webhook signing secret from environment
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	var event stripe.Event

	if endpointSecret != "" {
		// Verify webhook signature
		sigHeader := ctx.GetHeader("Stripe-Signature")
		event, err = webhook.ConstructEvent(payload, sigHeader, endpointSecret)
		if err != nil {
			log.Printf("Webhook signature verification failed: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
			return
		}
	} else {
		// For development without signature verification
		log.Println("Warning: STRIPE_WEBHOOK_SECRET not set, skipping signature verification")
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
	}

	log.Printf("Received Stripe webhook event: %s", event.Type)

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing checkout.session.completed: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing event"})
			return
		}

		if err := c.paymentService.HandleCheckoutSessionCompleted(&session); err != nil {
			log.Printf("Error handling checkout.session.completed: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing event"})
			return
		}

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing checkout.session.expired: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing event"})
			return
		}

		if err := c.paymentService.HandleCheckoutSessionExpired(&session); err != nil {
			log.Printf("Error handling checkout.session.expired: %v", err)
			// Don't return error for expired sessions
		}

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	ctx.JSON(http.StatusOK, gin.H{"received": true})
}
