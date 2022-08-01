package mule_test_exporter

import (
	"bufio"
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"os"
)

const (
	typeStr         = "mule-exporter"
	defaultFileName = "./spans.txt"
)

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		FileName:         defaultFileName,
	}
}

func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(typeStr,
		createDefaultConfig,
		component.WithTracesExporterAndStabilityLevel(createTracesExporter, component.StabilityLevelBeta),
	)
}

func createTracesExporter(_ context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.TracesExporter, error) {
	return &MuleExporter{config: config.(*Config)}, nil
}

type MuleExporter struct {
	config  *Config
	Started bool
	Stopped bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	f os.File
	w bufio.Writer
)

func (m MuleExporter) Start(ctx context.Context, host component.Host) error {
	m.Started = true
	file, err := os.Create(m.config.FileName)
	w = *bufio.NewWriter(file)
	f = *file

	check(err)

	return nil
}

func (m MuleExporter) Shutdown(ctx context.Context) error {
	m.Stopped = true
	return f.Close()
}

func (m MuleExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (m MuleExporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			spans := ils.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				writeToFile(span)
			}
		}
	}
	return nil
}

func writeToFile(span ptrace.Span) {
	fmt.Fprintln(&w, "")
	fmt.Fprintln(&w, "************************************** BEGIN SPAN **************************************")
	fmt.Fprintln(&w, "Span Name: "+span.Name())
	fmt.Fprintln(&w, "Trace ID: "+span.TraceID().HexString())
	fmt.Fprintln(&w, "Parent ID: "+span.ParentSpanID().HexString())
	fmt.Fprintln(&w, "ID: "+span.SpanID().HexString())
	fmt.Fprintln(&w, "Kind: "+span.Kind().String())
	fmt.Fprintln(&w, "Start time: "+span.StartTimestamp().String())
	fmt.Fprintln(&w, "End time: "+span.EndTimestamp().String())
	fmt.Fprintln(&w, "")
	fmt.Fprintln(&w, "Attributes")
	fmt.Fprintln(&w, "==========")
	fmt.Fprintln(&w, "")
	logAttributes(span.Attributes())
	fmt.Fprintln(&w, "")
	fmt.Fprintln(&w, "************************************** END SPAN **************************************")
	fmt.Fprintln(&w, "")
	w.Flush()
}

func logAttributes(m pcommon.Map) {
	if m.Len() == 0 {
		return
	}

	m.Range(func(k string, v pcommon.Value) bool {
		fmt.Fprintln(&w, k+" => "+v.StringVal())
		return true
	})
}
