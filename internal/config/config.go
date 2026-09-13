package config

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	HTTPAddress string
	DatabaseURL string
	LogLevel    string
}

// Load reads an optional .env file; process environment values take precedence.
func Load() (Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("LOG_LEVEL", "info")
	if err := v.ReadInConfig(); err != nil {
		var missing viper.ConfigFileNotFoundError
		if !errors.As(err, &missing) {
			return Config{}, fmt.Errorf("read configuration: %w", err)
		}
	}
	c := Config{Environment: v.GetString("APP_ENV"), HTTPAddress: v.GetString("HTTP_ADDR"), DatabaseURL: v.GetString("DATABASE_URL"), LogLevel: v.GetString("LOG_LEVEL")}
	return c, c.Validate()
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return errors.New("DATABASE_URL is required")
	}
	if _, _, err := net.SplitHostPort(c.HTTPAddress); err != nil {
		return errors.New("HTTP_ADDR must use host:port format, for example :8080")
	}
	if c.Environment != "development" && c.Environment != "test" && c.Environment != "production" {
		return errors.New("APP_ENV must be development, test, or production")
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return errors.New("LOG_LEVEL must be debug, info, warn, or error")
	}
	return nil
}
