package mule_test_exporter

import (
	"fmt"
	"go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	Prefix                  string `mapstructure:"prefix"`
}

func (cfg *Config) Validate() error {
	if cfg.Prefix == "Error" {
		return fmt.Errorf("cannot be Error")
	}
	return nil
}
