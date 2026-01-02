package stripe

import (
	"os"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
)

type StripeClient struct{}

func NewStripeClient() *StripeClient {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &StripeClient{}
}

// CreatePaymentIntent creates a PaymentIntent for client-side confirmation
// Returns the client_secret which is used by the frontend to confirm the payment
func (c *StripeClient) CreatePaymentIntent(orderID string, amount float64, currency string) (*stripe.PaymentIntent, error) {
	// Stripe expects amount in cents (smallest currency unit)
	amountInCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_id": orderID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	return paymentintent.New(params)
}

// ConfirmPaymentIntent confirms a payment intent with a payment method
// This is used for server-side confirmation (testing/sandbox)
func (c *StripeClient) ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	return paymentintent.Confirm(paymentIntentID, params)
}

// GetPaymentIntent retrieves a payment intent by ID
func (c *StripeClient) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(paymentIntentID, nil)
}

// CreateAndConfirmPaymentIntent creates and immediately confirms a payment (for testing)
// Uses Stripe test payment methods
func (c *StripeClient) CreateAndConfirmPaymentIntent(orderID string, amount float64, currency string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	// Stripe expects amount in cents (smallest currency unit)
	amountInCents := int64(amount * 100)

	// If no payment method provided, use Stripe's test card
	if paymentMethodID == "" {
		paymentMethodID = "pm_card_visa" // Stripe test payment method
	}

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amountInCents),
		Currency:      stripe.String(currency),
		PaymentMethod: stripe.String(paymentMethodID),
		Confirm:       stripe.Bool(true), // Automatically confirm
		Metadata: map[string]string{
			"order_id": orderID,
		},
		// Disable redirect-based payment methods for server-side confirmation
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"), // Prevents redirect-based methods
		},
	}

	return paymentintent.New(params)
}
