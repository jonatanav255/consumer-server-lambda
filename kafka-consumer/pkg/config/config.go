package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka KafkaConfig
	AWS   AWSConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type AWSConfig struct {
	Region string
}

func LoadConfig() *Config {
	// Set default values for Kafka configuration
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "events")

	// Automatically override the defaults with environment variables if available
	viper.AutomaticEnv()

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	// Create and return the configuration object with loaded settings
	return &Config{
		Kafka: KafkaConfig{
			Brokers: viper.GetStringSlice("kafka.brokers"),
			Topic:   viper.GetString("kafka.topic"),
		},
		AWS: AWSConfig{
			Region: awsCfg.Region,
		},
	}
}
