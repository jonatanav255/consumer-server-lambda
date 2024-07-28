package config

// import (
// 	"github.com/spf13/viper"
//             )

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

func LoadConfig() *Config {
	return &Config{}
}
