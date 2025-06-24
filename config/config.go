package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Service       Service
	Mode          Mode
	SearchIndex   string
	FlushInterval time.Duration
	Timeout       time.Duration

	Port string

	DefaultFields *map[string]string

	hostname string
}

type Service struct {
	Name    string
	Version string
}

func (cfg *Config) Ensure() error {
	if cfg.Service.Name == "" || cfg.Service.Version == "" {
		return errors.New("app not configured")
	}

	if cfg.Mode != Noop && cfg.Mode != Local && cfg.Mode != Debug && cfg.Mode != Development && cfg.Mode != Production {
		return errors.New("invalid log mode")
	}

	if cfg.SearchIndex == "" {
		cfg.SearchIndex = cfg.Service.Name
	}

	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = 30
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 10
	}

	if cfg.Port == "" {
		cfg.Port = "80"
	}

	var err error
	cfg.hostname, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("invalid hostname: %w", err)
	}

	return nil
}

func (cfg Config) GetHostname() string {
	return cfg.hostname
}

func (cfg Config) GetSearchIndex() string {
	if cfg.Mode == Development {
		return cfg.SearchIndex
	}
	return cfg.SearchIndex
}
