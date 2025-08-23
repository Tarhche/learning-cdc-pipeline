package messagequeue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQRepository implements the MessageQueueRepository interface
type RabbitMQRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitMQConfig
}

// RabbitMQConfig holds the configuration for RabbitMQ connection
type RabbitMQConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	VHost    string
}

// NewRabbitMQRepository creates a new RabbitMQ repository instance
func NewRabbitMQRepository(config RabbitMQConfig) *RabbitMQRepository {
	return &RabbitMQRepository{
		config: config,
	}
}

// Connect establishes a connection to RabbitMQ
func (r *RabbitMQRepository) Connect() error {
	// Build connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.VHost,
	)

	// Establish connection
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	r.conn = conn
	r.channel = channel

	log.Printf("Connected to RabbitMQ at %s:%d", r.config.Host, r.config.Port)
	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQRepository) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

	return nil
}

// ConsumeMessages starts consuming messages from a queue
func (r *RabbitMQRepository) ConsumeMessages(queueName string, handler func([]byte) error) error {
	// Declare the queue
	queue, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Declare the exchange (Maxwell uses this)
	err = r.channel.ExchangeDeclare(
		"maxwell", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Bind the queue to the exchange with routing key for posts table
	err = r.channel.QueueBind(
		queueName,    // queue name
		"blog.posts", // routing key (database.table)
		"maxwell",    // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Set QoS for reliable message processing
	err = r.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming messages with manual acknowledgment
	msgs, err := r.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack (false for manual acknowledgment)
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("Started consuming messages from queue: %s with reliable delivery (manual ACK)", queueName)

	// Process messages
	for msg := range msgs {
		log.Printf("Received message: %s", string(msg.Body))

		// Handle the message
		if err := handler(msg.Body); err != nil {
			log.Printf("Error processing message: %v", err)
			// Reject the message and requeue it for retry
			if err := msg.Nack(false, true); err != nil {
				log.Printf("Failed to NACK message: %v", err)
			}
			continue
		}

		// Acknowledge the message after successful processing
		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to ACK message: %v", err)
		} else {
			log.Printf("Successfully processed and acknowledged message")
		}
	}

	return nil
}

// PublishMessage publishes a message to an exchange
func (r *RabbitMQRepository) PublishMessage(exchange, routingKey string, message []byte) error {
	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
