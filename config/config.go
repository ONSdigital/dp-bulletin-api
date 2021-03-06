package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config represents service configuration for dp-bulletin-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   ":24200",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        10 * time.Second,
		HealthCheckCriticalTimeout: time.Minute,
	}

	return cfg, envconfig.Process("", cfg)
}
