package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"blog-cdc-search/application/service"
	"blog-cdc-search/infrastructure/messagequeue"
	"blog-cdc-search/infrastructure/searchindex"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {
	// Parse command line flags
	var (
		rabbitMQHost     = flag.String("rabbitmq-host", getEnv("RABBITMQ_HOST", "localhost"), "RabbitMQ host")
		rabbitMQPort     = flag.Int("rabbitmq-port", getEnvInt("RABBITMQ_PORT", 5672), "RabbitMQ port")
		rabbitMQUser     = flag.String("rabbitmq-user", getEnv("RABBITMQ_USER", "admin"), "RabbitMQ username")
		rabbitMQPassword = flag.String("rabbitmq-password", getEnv("RABBITMQ_PASSWORD", "admin123"), "RabbitMQ password")
		rabbitMQVHost    = flag.String("rabbitmq-vhost", getEnv("RABBITMQ_VHOST", "/"), "RabbitMQ vhost")
		typesenseHost    = flag.String("typesense-host", getEnv("TYPESENSE_HOST", "localhost"), "Typesense host")
		typesensePort    = flag.Int("typesense-port", getEnvInt("TYPESENSE_PORT", 8108), "Typesense port")
		typesenseAPIKey  = flag.String("typesense-api-key", getEnv("TYPESENSE_API_KEY", "xyz"), "Typesense API key")
		queueName        = flag.String("queue", getEnv("QUEUE_NAME", "cdc-posts"), "Queue name to consume from")
	)
	flag.Parse()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting CDC Service...")

	// Create RabbitMQ repository
	rabbitMQConfig := messagequeue.RabbitMQConfig{
		Host:     *rabbitMQHost,
		Port:     *rabbitMQPort,
		Username: *rabbitMQUser,
		Password: *rabbitMQPassword,
		VHost:    *rabbitMQVHost,
	}
	rabbitMQRepo := messagequeue.NewRabbitMQRepository(rabbitMQConfig)

	// Create Typesense repository
	typesenseConfig := searchindex.TypesenseConfig{
		Host:   *typesenseHost,
		Port:   *typesensePort,
		APIKey: *typesenseAPIKey,
	}
	typesenseRepo := searchindex.NewTypesenseRepository(typesenseConfig)

	// Create CDC service
	cdcService := service.NewCDCService(rabbitMQRepo, typesenseRepo)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	// Start the CDC service
	log.Printf("Starting CDC service with queue: %s", *queueName)
	if err := cdcService.StartCDC(ctx, *queueName); err != nil {
		log.Fatalf("CDC service failed: %v", err)
	}

	log.Println("CDC service stopped")
}
