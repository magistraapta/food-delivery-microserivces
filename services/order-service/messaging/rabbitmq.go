package messaging

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchanges
	OrderEventsExchange   = "order.events"
	PaymentEventsExchange = "payment.events"

	// Queues
	PaymentSuccessQueue = "order.payment.success" // Order service's queue for payment.success events
	PaymentFailedQueue  = "order.payment.failed"  // Order service's queue for payment.failed events

	// Routing Keys
	OrderCreatedRoutingKey   = "order.created"
	PaymentSuccessRoutingKey = "payment.success"
	PaymentFailedRoutingKey  = "payment.failed"
)

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient() (*RabbitMQClient, error) {
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

	client := &RabbitMQClient{
		Conn:    conn,
		Channel: channel,
	}

	return client, nil
}

// SetupQueuesAndExchanges declares all necessary exchanges and queues for order service
func (c *RabbitMQClient) SetupQueuesAndExchanges() error {
	// Declare order events exchange (published by order-service)
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

	// Declare queue for consuming payment.success events
	_, err = c.Channel.QueueDeclare(
		PaymentSuccessQueue, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("Declared queue: %s", PaymentSuccessQueue)

	// Bind payment.success queue to payment events exchange
	err = c.Channel.QueueBind(
		PaymentSuccessQueue,      // queue name
		PaymentSuccessRoutingKey, // routing key
		PaymentEventsExchange,    // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log.Printf("Bound queue %s to exchange %s with routing key %s", PaymentSuccessQueue, PaymentEventsExchange, PaymentSuccessRoutingKey)

	// Declare queue for consuming payment.failed events
	_, err = c.Channel.QueueDeclare(
		PaymentFailedQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("Declared queue: %s", PaymentFailedQueue)

	// Bind payment.failed queue to payment events exchange
	err = c.Channel.QueueBind(
		PaymentFailedQueue,      // queue name
		PaymentFailedRoutingKey, // routing key
		PaymentEventsExchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log.Printf("Bound queue %s to exchange %s with routing key %s", PaymentFailedQueue, PaymentEventsExchange, PaymentFailedRoutingKey)

	return nil
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQClient) Close() {
	if c.Channel != nil {
		c.Channel.Close()
		log.Println("RabbitMQ channel closed")
	}
	if c.Conn != nil {
		c.Conn.Close()
		log.Println("RabbitMQ connection closed")
	}
}
