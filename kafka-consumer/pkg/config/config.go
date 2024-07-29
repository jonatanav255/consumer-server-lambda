package config

import (
	"github.com/spf13/viper"
)

// Config struct holds the entire configuration for the application.
type Config struct {
    Kafka KafkaConfig
}

// KafkaConfig struct holds the Kafka-specific configuration settings.
type KafkaConfig struct {
    Brokers []string
    Topic   string
}

// LoadConfig function initializes and loads the configuration settings.
func LoadConfig() *Config {
    // Set default values for Kafka configuration
    viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
    viper.SetDefault("kafka.topic", "events")

    // Automatically override the defaults with environment variables if available
    viper.AutomaticEnv()

    // Create and return the configuration object with loaded settings
    return &Config{
        Kafka: KafkaConfig{
            Brokers: viper.GetStringSlice("kafka.brokers"),
            Topic:   viper.GetString("kafka.topic"),
        },
    }
}
