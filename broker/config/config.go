package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds all the configurations for the broker
type Config struct {
	Port          string `yaml:"port"`
	Discovery     bool   `yaml:"discovery"`
	LabelSelector string `yaml:"labelSelector"`
}

func (c *Config) applyDefaults() {
	if c.Port == "" {
		c.Port = "50051"
	}
}

// Load loads the configuration from a config file
func Load(filename string) (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, err
	}

	config.applyDefaults()
	return config, nil
}
