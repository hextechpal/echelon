package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug bool `envconfig:"ECHELON_DEBUG"`

	ServerAddress string   `envconfig:"ECHELON_SERVER_ADDRESS"`
	BindAddress   string   `envconfig:"ECHELON_BIND_ADDRESS"`
	JoinAddresses []string `envconfig:"ECHELON_JOIN_ADDRESSES"`

	DB struct {
		Host     string `envconfig:"ECHELON_DB_HOST"`     // Only sql is supported for now
		Port     int32  `envconfig:"ECHELON_DB_PORT"`     // Only sql is supported for now
		User     string `envconfig:"ECHELON_DB_USER"`     // Only sql is supported for now
		Password string `envconfig:"ECHELON_DB_PASSWORD"` // Only sql is supported for now
		Database string `envconfig:"ECHELON_DB_DATABASE"` // Only sql is supported for now
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("server=%s, bind=%s, joinAddrs=%v", c.ServerAddress, c.BindAddress, c.JoinAddresses)
}

func Load() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
