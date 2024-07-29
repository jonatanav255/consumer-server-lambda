package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka-consumer/pkg/config"

	"github.com/IBM/sarama"
)

// KafkaConsumer represents a Sarama consumer group consumer
type KafkaConsumer struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
    return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *KafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
        sess.MarkMessage(message, "")
    }
    return nil
}

// StartKafkaConsumer initializes and starts the Kafka consumer
func StartKafkaConsumer(cfg *config.Config) {
    // Create a new Sarama configuration
    saramaConfig := sarama.NewConfig()
    saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

    // Create a new consumer group
    consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "my-consumer-group", saramaConfig)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    }

    // Create a new consumer instance
    consumer := &KafkaConsumer{}

    // Set up a signal notifier for graceful shutdown
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    // Consume messages in a loop
    ctx := context.Background()
    go func() {
        for {
            if err := consumerGroup.Consume(ctx, []string{cfg.Kafka.Topic}, consumer); err != nil {
                log.Fatalf("Error from consumer: %v", err)
            }
            // Check if context was cancelled, signaling a shutdown
            if ctx.Err() != nil {
                return
            }
        }
    }()

    <-signals
    log.Println("Shutting down consumer...")

    if err := consumerGroup.Close(); err != nil {
        log.Fatalf("Error closing consumer group: %v", err)
    }
}
