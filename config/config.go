package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug bool `envconfig:"ECHELON_DEBUG"`

	ServerPort int      `envconfig:"ECHELON_SERVER_PORT"`
	BindPort   int      `envconfig:"ECHELON_BIND_PORT"`
	JoinAddrs  []string `envconfig:"ECHELON_JOIN_ADDRS"`

	DB struct {
		Host     string `envconfig:"ECHELON_DB_HOST"`     // Only sql is supported for now
		Port     int32  `envconfig:"ECHELON_DB_PORT"`     // Only sql is supported for now
		User     string `envconfig:"ECHELON_DB_USER"`     // Only sql is supported for now
		Password string `envconfig:"ECHELON_DB_PASSWORD"` // Only sql is supported for now
		Database string `envconfig:"ECHELON_DB_DATABASE"` // Only sql is supported for now
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("serverPort=%d, bindPort=%d, joinAddrs=%v", c.ServerPort, c.BindPort, c.JoinAddrs)
}

func Load() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
