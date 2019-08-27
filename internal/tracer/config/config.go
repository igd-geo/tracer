package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	envDeploymentEnvironment = "DEPLOYMENT_ENVIRONMENT"
	envDatabaseURL           = "DATABASE_URL"
	envBrokerURL             = "BROKER_URL"
	envBrokerUser            = "BROKER_USER"
	envBrokerPassword        = "BROKER_PASSWORD"
	envBatchSizeLimit        = "BATCH_SIZE_LIMIT"
	envBatchTimeout          = "BATCH_TIMEOUT"

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
	config := Config{
		DB:             defaultDB,
		Broker:         defaultBroker,
		BatchSizeLimit: defaultBatchSizeLimit,
		BatchTimeout:   defaultBatchTimeout,
	}
	if os.Getenv(envDeploymentEnvironment) == "PROD" {
		brokerUser := os.Getenv(envBrokerUser)
		brokerPassword := os.Getenv(envBrokerPassword)
		brokerURL := os.Getenv(envBrokerURL)

		config.DB = os.Getenv(envDatabaseURL)
		config.Broker = fmt.Sprintf("amqp://%s:%s@%s", brokerUser, brokerPassword, brokerURL)

		batchSizeLimit, err := strconv.Atoi(os.Getenv(envBatchSizeLimit))
		if err != nil || batchSizeLimit < 0 {
			log.Fatal("could not parse batch size limit, value must be > 0")
		}
		config.BatchSizeLimit = batchSizeLimit

		batchTimeout, err := strconv.Atoi(os.Getenv(envBatchTimeout))
		if err != nil || batchTimeout < 0 {
			log.Fatal("could not parse batch timeout, value must be > 0")
		}
		config.BatchTimeout = batchTimeout
	}
	return &config
}
