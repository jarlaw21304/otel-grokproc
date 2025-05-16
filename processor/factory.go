package processor

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/processor"
    "go.opentelemetry.io/collector/processor/processorhelper"
    "go.opentelemetry.io/collector/consumer"
)

const (
    TypeStr = "grokproc"
)

func Factory() component.ProcessorFactory {
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
    params processor.CreateSettings,
    cfg component.Config,
    next consumer.Logs,
) (processor.Logs, error) {
    pcfg := cfg.(*Config)
    proc := NewProcessor(pcfg)
    return processorhelper.NewLogsProcessor(
        ctx,
        params,
        cfg,
        next,
        proc.ProcessLogs,
        processorhelper.WithCapabilities(processorhelper.Capabilities{MutatesData: true}),
    )
}
