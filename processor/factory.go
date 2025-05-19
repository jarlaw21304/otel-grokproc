package processor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

// Type is how this processor is represented in OTEL Collector YAML configs.
const Type = "grokproc"

// createDefaultConfig returns a default config struct.
func createDefaultConfig() component.Config {
	return &Config{
		Pattern:          "",
		PatternDirectory: "",
		FieldMap:         make(map[string]string),
		ExtractCEF:       false,
	}
}

// createLogsProcessor constructs the processor when the collector starts its pipeline.
func createLogsProcessor(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	next consumer.Logs,
) (processor.Logs, error) {
	c := cfg.(*Config)
	p, err := NewProcessor(c)
	if err != nil {
		return nil, err
	}
	p.next = next // Wire up the next consumer/exporter in the pipeline.
	return p, nil
}

// NewFactory registers the processor type so it’s available to the Collector.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		Type,
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, component.StabilityLevelBeta),
	)
}
