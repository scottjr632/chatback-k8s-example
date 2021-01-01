package config

import "github.com/scottjr632/chatback/server/broker"

// Config holds all the configuration options for
// chatback
type Config struct {
	DB     *DB
	Broker *broker.Config
}

// New creates a new config
func New() Config {
	config := Config{
		newDB(),
		&broker.Config{},
	}
	return config
}
