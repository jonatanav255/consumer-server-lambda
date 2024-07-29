package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka-consumer/pkg/config"
	"kafka-consumer/pkg/consumer"

	"github.com/IBM/sarama"
)

func StartKafkaConsumer(cfg *config.Config) {
    // Create a new consumer group
    consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "consumer-group", nil)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    }
    defer consumerGroup.Close()

    // Create a new consumer instance
    consumer := consumer.Consumer{}

    // Set up a channel to listen for OS signals to gracefully shut down the consumer
    sigterm := make(chan os.Signal, 1)
    signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

    ctx := context.Background()

    go func() {
        for {
            if err := consumerGroup.Consume(ctx, []string{cfg.Kafka.Topic}, &consumer); err != nil {
                fmt.Printf("Error from consumer: %v", err)
            }
        }
    }()

    <-sigterm
    fmt.Println("Terminating: via signal")
}
