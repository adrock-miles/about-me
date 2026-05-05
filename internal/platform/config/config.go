// Package config loads runtime configuration from a YAML file with
// PORTFOLIO_-prefixed environment variable overrides.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds every runtime setting for the portfolio server.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Email  EmailConfig  `mapstructure:"email"`
	Site   SiteConfig   `mapstructure:"site"`
}

// ServerConfig controls the HTTP listener and lifecycle.
type ServerConfig struct {
	Host          string        `mapstructure:"host"`
	Port          int           `mapstructure:"port"`
	ReadTimeout   time.Duration `mapstructure:"read_timeout"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`
	IdleTimeout   time.Duration `mapstructure:"idle_timeout"`
	ShutdownGrace time.Duration `mapstructure:"shutdown_grace"`
}

// EmailConfig selects a delivery provider and supplies its credentials.
//
// Provider values:
//   - ""        → dev mode, messages are logged via slog (no network)
//   - "smtp"    → net/smtp via STARTTLS on Email.SMTP.*
//   - "resend"  → Resend HTTP API using Email.Resend.APIKey
type EmailConfig struct {
	Provider string       `mapstructure:"provider"`
	From     string       `mapstructure:"from"`
	To       string       `mapstructure:"to"`
	SMTP     SMTPConfig   `mapstructure:"smtp"`
	Resend   ResendConfig `mapstructure:"resend"`
}

// SMTPConfig holds SMTP credentials.
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ResendConfig holds Resend API credentials.
type ResendConfig struct {
	APIKey string `mapstructure:"api_key"`
}

// SiteConfig holds public-facing site metadata.
type SiteConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

// Load reads ./configs/config.yaml (and the working dir as a fallback) and
// applies PORTFOLIO_-prefixed environment overrides on top.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetEnvPrefix("PORTFOLIO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Railway / generic PaaS still pass PORT directly.
	if err := v.BindEnv("server.port", "PORT"); err != nil {
		return nil, fmt.Errorf("binding PORT: %w", err)
	}

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)
	v.SetDefault("server.shutdown_grace", 10*time.Second)

	v.SetDefault("email.provider", "")
	v.SetDefault("email.from", "noreply@adrock.dev")
	v.SetDefault("email.to", "miles.adrock@gmail.com")
	v.SetDefault("email.smtp.port", 587)

	v.SetDefault("site.base_url", "https://adrock.dev")
}
