package config

import (
	"io/ioutil"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/broker"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"
	"gopkg.in/yaml.v2"
)

// Config contains configuration information
type Config struct {
	Log      *LogConfig       `yaml:"log"`
	Batch    *BatchConfig     `yaml:"batch"`
	Database *database.Config `yaml:"database"`
	Broker   *broker.Config   `yaml:"broker"`
	Arbiter  *ArbiterConfig   `yaml:"arbiter"`
}

// LogConfig Configuration for the logger
type LogConfig struct {
	Debug    bool   `yaml:"debug"`
	Exchange string `yaml:"exchange"`
}

// BatchConfig Configuration for database batching
type BatchConfig struct {
	Size    int `yaml:"size"`
	Timeout int `yaml:"timeout"`
}

// ArbiterConfig Configuration for the arbiter registry endpoints
type ArbiterConfig struct {
	ScenarioRegistry *ScenarioRegistryConfig `yaml:"scenarioRegistry"`
	ServiceRegistry  *ServiceRegistryConfig  `yaml:"serviceRegistry"`
	UserRegistry     *UserRegistryConfig     `yaml:"userRegistry"`
}

// ScenarioRegistryConfig Configuration for the scenario registry endpoint
type ScenarioRegistryConfig struct {
	URL string `yaml:"url"`
}

// ServiceRegistryConfig Cnfiguration for the service registry endpoint
type ServiceRegistryConfig struct {
	URL string `yaml:"url"`
}

// UserRegistryConfig Configuration for the user registry endpoint
type UserRegistryConfig struct {
	URL string `yaml:"url"`
}

// Parse Parses a yaml file to a config struct and returns it
func Parse(configPath string) (*Config, error) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configFile, config)
	return config, err
}
