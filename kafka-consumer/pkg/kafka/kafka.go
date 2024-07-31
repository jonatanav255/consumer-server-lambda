package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"kafka-consumer/pkg/config" // Your local config package

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config" // Alias the AWS SDK's config package
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Consumer struct {
	s3Client *s3.Client
	bucket   string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		// Upload message to S3
		_, err := consumer.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(consumer.bucket),
			Key:    aws.String(fmt.Sprintf("%d-%s", message.Offset, message.Topic)),
			Body:   strings.NewReader(string(message.Value)),
		})
		if err != nil {
			log.Printf("Failed to upload data to S3: %v", err)
		} else {
			log.Printf("Successfully uploaded data to S3")
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func StartKafkaConsumer(cfg *config.Config) {
	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "consumer-group", nil)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Create a new AWS config
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	s3Client := s3.NewFromConfig(awsCfg)

	// Create a new consumer instance
	consumer := Consumer{
		s3Client: s3Client,
		bucket:   "your-s3-bucket-name", // Update with your bucket name
	}

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
