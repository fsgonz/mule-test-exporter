package mule_test_exporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	typeStr       = "mule-exporter"
	defaultPrefix = "Mule-"
)

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		Prefix:           defaultPrefix,
	}
}

func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(typeStr,
		createDefaultConfig)
}
