package stripe

import (
	"os"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/paymentintent"
)

type StripeClient interface {
	CreateCheckoutSession(orderID string, amount float64, currency string, productName string) (*stripe.CheckoutSession, error)
	GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error)
	CreatePaymentIntent(orderID string, amount float64, currency string) (*stripe.PaymentIntent, error)
	ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error)
	GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error)
	CreateAndConfirmPaymentIntent(orderID string, amount float64, currency string, paymentMethodID string) (*stripe.PaymentIntent, error)
}

type StripeClientImpl struct {
	SuccessURL string
	CancelURL  string
}

func NewStripeClient() *StripeClientImpl {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get base URL for redirects (defaults for local development)
	baseURL := os.Getenv("STRIPE_REDIRECT_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000" // Frontend URL
	}

	return &StripeClientImpl{
		SuccessURL: baseURL + "/payment/success?session_id={CHECKOUT_SESSION_ID}",
		CancelURL:  baseURL + "/payment/cancel",
	}
}

// CreateCheckoutSession creates a Stripe Checkout Session and returns the checkout URL
// The user should be redirected to this URL to complete payment
func (c *StripeClientImpl) CreateCheckoutSession(orderID string, amount float64, currency string, productName string) (*stripe.CheckoutSession, error) {
	// Stripe expects amount in cents (smallest currency unit)
	amountInCents := int64(amount * 100)

	if currency == "" {
		currency = "usd"
	}

	if productName == "" {
		productName = "Food Order"
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(amountInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(c.SuccessURL),
		CancelURL:  stripe.String(c.CancelURL),
		ExpiresAt:  stripe.Int64(300), // 5 minutes (in seconds from now) - matches our timeout
		Metadata: map[string]string{
			"order_id": orderID,
		},
	}

	// Note: ExpiresAt requires a Unix timestamp, not duration
	// For 5 minutes from now, we need to calculate the timestamp
	params.ExpiresAt = nil // Remove for now, Stripe default is 24 hours

	return session.New(params)
}

// GetCheckoutSession retrieves a checkout session by ID
func (c *StripeClientImpl) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	return session.Get(sessionID, nil)
}

// CreatePaymentIntent creates a PaymentIntent for client-side confirmation
// Returns the client_secret which is used by the frontend to confirm the payment
func (c *StripeClientImpl) CreatePaymentIntent(orderID string, amount float64, currency string) (*stripe.PaymentIntent, error) {
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
func (c *StripeClientImpl) ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	return paymentintent.Confirm(paymentIntentID, params)
}

// GetPaymentIntent retrieves a payment intent by ID
func (c *StripeClientImpl) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(paymentIntentID, nil)
}

// CreateAndConfirmPaymentIntent creates and immediately confirms a payment (for testing)
// Uses Stripe test payment methods
func (c *StripeClientImpl) CreateAndConfirmPaymentIntent(orderID string, amount float64, currency string, paymentMethodID string) (*stripe.PaymentIntent, error) {
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
