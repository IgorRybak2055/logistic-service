// Package config describes ragger's configuration.
package config

import (
	"context"

	"github.com/heetch/confita"

	"github.com/IgorRybak2055/logistic-service/internal/ragger"
	"github.com/IgorRybak2055/logistic-service/internal/storage"
	"github.com/IgorRybak2055/logistic-service/pkg/email"
)

// Config stores configs for ragger
type Config struct {
	HTTP *ragger.HTTPConfig
	DB   *storage.Config
	Mail *email.Config
}

// NewConfig returns Config with values from environment variables
func NewConfig(ctx context.Context) (*Config, error) {
	var cfg = &Config{
		HTTP: &ragger.HTTPConfig{},
		DB:   &storage.Config{},
		Mail: &email.Config{},
	}

	if err := confita.NewLoader().Load(ctx, cfg); err != nil {
		return nil, err
	}

	cfg.Mail.RestoreURL = cfg.HTTP.FullRestoreURL()

	return cfg, nil
}
