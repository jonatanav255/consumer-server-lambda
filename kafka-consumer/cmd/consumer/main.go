package main

import (
	"kafka-consumer/pkg/config"
	"kafka-consumer/pkg/kafka"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize and start the Kafka consumer
	kafka.StartKafkaConsumer(cfg)

	// Note: The consumer will run until interrupted (Ctrl+C)
	// Cleanup and shutdown logic is handled in StartKafkaConsumer
}
