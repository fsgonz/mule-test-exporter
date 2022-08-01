package mule_test_exporter

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
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
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter),
	)
}

func createTracesExporter(_ context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.TracesExporter, error) {
	return nopExporterInstance, nil
}

var nopExporterInstance = &nopExporter{
	Consumer: consumertest.NewNop(),
}

type nopExporter struct {
	nopComponent
	consumertest.Consumer
}

type nopComponent struct {
	component.StartFunc
	component.ShutdownFunc
}
