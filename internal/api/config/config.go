package config

import (
	"fmt"
	"os"
)

const (
	envDeploymentEnvironment = "DEPLOYMENT_ENVIRONMENT"
	envDatabaseURL           = "DATABASE_URL"
	envBrokerURL             = "BROKER_URL"
	envBrokerUser            = "BROKER_USER"
	envBrokerPassword        = "BROKER_PASSWORD"
	envPort                  = "API_PORT"

	defaultDB     = "localhost:9080"
	defaultBroker = "amqp://guest:guest@localhost:5672/"
	defaultPort   = ":1234"
)

// Config contains configuration information
type Config struct {
	DB     string
	Broker string
	Port   string
}

// New returns an empty Config struct
func New() *Config {
	config := Config{
		DB:     defaultDB,
		Broker: defaultBroker,
		Port:   defaultPort,
	}
	if os.Getenv(envDeploymentEnvironment) == "PROD" {
		brokerUser := os.Getenv(envBrokerUser)
		brokerPassword := os.Getenv(envBrokerPassword)
		brokerURL := os.Getenv(os.Getenv(envBrokerURL))

		config.DB = os.Getenv(envDatabaseURL)
		config.Broker = fmt.Sprintf("amqp://%s:%s@%s", brokerUser, brokerPassword, brokerURL)
		config.Port = os.Getenv(envPort)
	}
	return &config
}
