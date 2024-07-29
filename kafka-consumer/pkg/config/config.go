package config

import (
	"github.com/spf13/viper"
)

// Config struct holds the entire configuration for the application.
// In this case, it contains a nested KafkaConfig struct.
type Config struct {
    Kafka KafkaConfig
}

// KafkaConfig struct holds the Kafka-specific configuration settings.
// It includes the list of broker addresses and the topic to consume from.
type KafkaConfig struct {
    Brokers []string
    Topic   string
}

// LoadConfig function initializes and loads the configuration settings.
func LoadConfig() *Config {
    // Set default values for Kafka configuration
    // If no environment variables are set, these defaults will be used.
    viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
    viper.SetDefault("kafka.topic", "events")

    // Automatically override the defaults with environment variables if available
    // For example, setting KAFKA_BROKERS or KAFKA_TOPIC environment variables will override the defaults.
    viper.AutomaticEnv()

    // Create and return the configuration object with loaded settings
    // viper.GetStringSlice and viper.GetString retrieve the values for the specified keys.
    return &Config{
        Kafka: KafkaConfig{
            Brokers: viper.GetStringSlice("kafka.brokers"),
            Topic:   viper.GetString("kafka.topic"),
        },
    }
}
