package mule_test_exporter

import (
	"go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	FileName                string `mapstructure:"filename"`
}

func (cfg *Config) Validate() error {
	return nil
}
