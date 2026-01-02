package event

type OrderCreatedEvent struct {
	OrderID         string  `json:"order_id"`
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PaymentMethodID string  `json:"payment_method_id"`
}
