package config

import (
	"fmt"
	"os"
)

const (
	envAuthUsername = "AUTH_USERNAME"
	envAuthPassword = "AUTH_PASSWORD"
	envTableName    = "TABLE_NAME"
)

type Config struct {
	Username  string
	Password  string
	TableName string
}

func ReadConfig() (*Config, error) {
	cfg := new(Config)
	cfg.Username = os.Getenv(envAuthUsername)
	cfg.Password = os.Getenv(envAuthPassword)
	cfg.TableName = os.Getenv(envTableName)
	err := cfg.validate()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.Username == "" {
		return fmt.Errorf("environment variable '%s' must be specified", envAuthUsername)
	}
	if cfg.Password == "" {
		return fmt.Errorf("environment variable '%s' must be specified", envAuthPassword)
	}
	if cfg.TableName == "" {
		return fmt.Errorf("environment variable '%s' must be specified", envTableName)
	}
	return nil
}
