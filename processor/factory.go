package processor

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/config"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/processor"
)

const (
    TypeStr = "grokproc"
)

func Factory() processor.Factory {
    return processor.NewFactory(
        TypeStr,
        createDefaultConfig,
        processor.WithLogs(createLogsProcessor, component.StabilityLevelBeta),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        Pattern:      "",
        VendorFormat: "",
        FieldMap:     map[string]string{},
        ExtractCEF:   false,
    }
}

func createLogsProcessor(
    ctx context.Context,
    settings processor.CreateSettings,
    cfg component.Config,
    next consumer.Logs,
) (processor.Logs, error) {
    pcfg := cfg.(*Config)
    proc := NewProcessor(pcfg)
    return &grokProcProcessor{
        next: next,
        proc: proc,
    }, nil
}
