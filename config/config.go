package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Debug bool `envconfig:"ECHELON_DEBUG"`

	Name string `envconfig:"ECHELON_WORKER_NAME"`

	DB struct {
		Host     string `envconfig:"ECHELON_DB_HOST"`     // Only sql is supported for now
		Port     int32  `envconfig:"ECHELON_DB_PORT"`     // Only sql is supported for now
		User     string `envconfig:"ECHELON_DB_USER"`     // Only sql is supported for now
		Password string `envconfig:"ECHELON_DB_PASSWORD"` // Only sql is supported for now
		Database string `envconfig:"ECHELON_DB_DATABASE"` // Only sql is supported for now
	}

	Server struct {
		Host string `envconfig:"ECHELON_SERVER_HOST"`
		Port int    `envconfig:"ECHELON_SERVER_PORT"`
	}

	Serf struct {
		BindAddress string `envconfig:"ECHELON_BIND_ADDRESS"`
		BindPort    int    `envconfig:"ECHELON_BIND_PORT"`
	}
}

func Load() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
