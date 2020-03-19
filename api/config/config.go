package config

import (
	"io/ioutil"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/broker"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"
	"gopkg.in/yaml.v2"
)

// Config Configuration for the api service
type Config struct {
	Log      *LogConfig       `yaml:"log"`
	Database *database.Config `yaml:"database"`
	Broker   *broker.Config   `yaml:"broker"`
	API      *APIConfig       `yaml:"api"`
}

// APIConfig Configuration for the api's http server
type APIConfig struct {
	Port string `yaml:"port"`
}

// LogConfig Configuration for the lgger
type LogConfig struct {
	Debug    bool   `yaml:"debug"`
	Exchange string `yaml:"exchange"`
}

// Parse Parses an yaml file to a config struct and returns it
func Parse(configPath string) (*Config, error) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configFile, config)
	return config, err
}
