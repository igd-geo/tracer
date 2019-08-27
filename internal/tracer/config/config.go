package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	envEnvironment    = "ENVIRONMENT"
	envDatabaseURL    = "DATABASE_URL"
	envBrokerURL      = "BROKER_URL"
	envBrokerUser     = "BROKER_USER"
	envBrokerPassword = "BROKER_PASSWORD"
	envBatchSizeLimit = "BATCH_SIZE_LIMIT"
	envBatchTimeout   = "BATCH_TIMEOUT"

	defaultDB             = "localhost:9080"
	defaultBroker         = "amqp://guest:guest@localhost:5672/"
	defaultBatchSizeLimit = 1000
	defaultBatchTimeout   = 100
)

// Config contains configuration information
type Config struct {
	DB             string
	Broker         string
	BatchSizeLimit int
	BatchTimeout   int
}

// New returns an empty Config struct
func New() *Config {
	var err error
	config := Config{
		DB:             defaultDB,
		Broker:         defaultBroker,
		BatchSizeLimit: defaultBatchSizeLimit,
		BatchTimeout:   defaultBatchTimeout,
	}
	if os.Getenv(envEnvironment) == "PROD" {
		brokerUser := os.Getenv(envBrokerUser)
		brokerPassword := os.Getenv(envBrokerPassword)
		brokerURL := os.Getenv(os.Getenv(envBrokerURL))

		config.DB = os.Getenv(envDatabaseURL)
		config.Broker = fmt.Sprintf("amqp://%s:%s@%s", brokerUser, brokerPassword, brokerURL)

		config.BatchSizeLimit, err = strconv.Atoi(os.Getenv(envBatchSizeLimit))
		if err != nil {
			log.Fatal(err)
		}

		config.BatchTimeout, err = strconv.Atoi(os.Getenv(envBatchTimeout))
		if err != nil {
			log.Fatal(err)
		}

	}
	return &config
}
