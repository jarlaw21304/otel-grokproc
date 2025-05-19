package otelgrokproc

import (
	"context"
	"fmt"
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"otel-grokproc/processor/grokparse"
)

// Config definition for this processor (add struct tags as needed for yaml/OTel config).
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	PatternDirectory         string `mapstructure:"pattern_directory"`
	PatternName              string `mapstructure:"pattern_name"` // if static, or map field for dynamic cases
}

// Factory for processor registration.
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		"otelgrokproc",
		createDefaultConfig,
		component.WithLogsProcessor(createLogsProcessor),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID("otelgrokproc")),
		PatternDirectory:  "patterns",
		PatternName:       "", // Configure via yaml
	}
}

func createLogsProcessor(
	ctx context.Context,
	settings component.ProcessorCreateSettings,
	cfg component.Config,
	nextConsumer component.LogsConsumer,
) (component.LogsProcessor, error) {
	oCfg := cfg.(*Config)
	// Load grok patterns at startup.
	if err := grokparse.LoadAllPatternFiles(oCfg.PatternDirectory); err != nil {
		return nil, fmt.Errorf("could not load grok patterns: %w", err)
	}
	return &Processor{
		patternName:  oCfg.PatternName,
		nextConsumer: nextConsumer,
	}, nil
}

// Processor implements the OTel log processor API.
type Processor struct {
	patternName  string // example: "ASA302014"
	nextConsumer component.LogsConsumer
}

func (p *Processor) Capabilities() component.ProcessorCapabilities {
	return component.ProcessorCapabilities{MutatesData: true}
}

func (p *Processor) Start(ctx context.Context, host component.Host) error { return nil }
func (p *Processor) Shutdown(ctx context.Context) error                   { return nil }

func (p *Processor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	// Iterate all log records and parse using grok.
	resSlice := ld.ResourceLogs()
	for i := 0; i < resSlice.Len(); i++ {
		scopeSlice := resSlice.At(i).ScopeLogs()
		for j := 0; j < scopeSlice.Len(); j++ {
			logs := scopeSlice.At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				bodyStr := lr.Body().AsString()
				grokPattern := "%{" + p.patternName + "}"
				fields, err := grokparse.ParseLine(grokPattern, bodyStr)
				if err == nil {
					// Attach parsed fields as attributes.
					attrs := lr.Attributes()
					for k, v := range fields {
						attrs.PutStr(k, v)
					}
				} else {
					// Optionally, annotate parse errors or drop unmatched logs.
				}
			}
		}
	}
	return p.nextConsumer.ConsumeLogs(ctx, ld)
}
