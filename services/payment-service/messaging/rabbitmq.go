package messaging

import (
	"log"
	"os"
	"payment-service/messaging/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchanges
	OrderEventsExchange   = "order.events"
	PaymentEventsExchange = "payment.events"

	// Queues
	OrderCreatedQueue   = "payment.order.created" // Payment service's queue for order.created events
	PaymentSuccessQueue = "payment.success"
	PaymentFailedQueue  = "payment.failed"

	// Routing Keys
	OrderCreatedRoutingKey           = "order.created"
	PaymentSuccessRoutingKey         = "payment.success"
	PaymentFailedRoutingKey          = "payment.failed"
	PaymentCheckoutCreatedRoutingKey = "payment.checkout.created"
)

type RabbitmqClient interface {
	PublishPaymentSuccess(event events.PaymentSuccessEvent) error
	PublishPaymentFailed(event events.PaymentFailedEvent) error
	PublishPaymentCheckoutCreated(event events.PaymentCheckoutCreatedEvent) error
}

type RabbitmqClientImpl struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient() (*RabbitmqClientImpl, error) {
	// Get RabbitMQ URL from environment variable, fallback to default
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	client := &RabbitmqClientImpl{
		Conn:    conn,
		Channel: channel,
	}

	return client, nil
}

// SetupQueuesAndExchanges declares all necessary exchanges and queues for payment service
func (c *RabbitmqClientImpl) SetupQueuesAndExchanges() error {
	// Declare order events exchange (published by order-service, consumed by payment-service)
	err := c.Channel.ExchangeDeclare(
		OrderEventsExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("Declared exchange: %s", OrderEventsExchange)

	// Declare payment events exchange (published by payment-service, consumed by order-service)
	err = c.Channel.ExchangeDeclare(
		PaymentEventsExchange, // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("Declared exchange: %s", PaymentEventsExchange)

	// Declare queue for consuming order.created events
	_, err = c.Channel.QueueDeclare(
		OrderCreatedQueue, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("Declared queue: %s", OrderCreatedQueue)

	// Bind order.created queue to order events exchange
	err = c.Channel.QueueBind(
		OrderCreatedQueue,      // queue name
		OrderCreatedRoutingKey, // routing key
		OrderEventsExchange,    // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log.Printf("Bound queue %s to exchange %s with routing key %s", OrderCreatedQueue, OrderEventsExchange, OrderCreatedRoutingKey)

	return nil
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitmqClientImpl) Close() {
	if c.Channel != nil {
		c.Channel.Close()
		log.Println("RabbitMQ channel closed")
	}
	if c.Conn != nil {
		c.Conn.Close()
		log.Println("RabbitMQ connection closed")
	}
}
