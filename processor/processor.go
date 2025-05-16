package processor

import (
    "context"

    "github.com/yourorg/otel-grokproc/grokparse"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/pdata/plog"
    "go.opentelemetry.io/collector/processor"
)

// GrokProcessor contains config and parsing logic
type GrokProcessor struct {
    cfg *Config
}

// NewProcessor builds the GrokProcessor from config.
func NewProcessor(cfg *Config) *GrokProcessor {
    return &GrokProcessor{cfg: cfg}
}

// ProcessLogs parses log bodies with grok and sets extracted fields as OTEL attributes.
func (p *GrokProcessor) ProcessLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
    for i := 0; i < logs.ResourceLogs().Len(); i++ {
        rl := logs.ResourceLogs().At(i)
        for j := 0; j < rl.ScopeLogs().Len(); j++ {
            sl := rl.ScopeLogs().At(j)
            for k := 0; k < sl.LogRecords().Len(); k++ {
                lr := sl.LogRecords().At(k)
                body := lr.Body().AsString()
                pattern := p.cfg.Pattern
                if pattern == "" && p.cfg.VendorFormat != "" {
                    pattern = "%{" + p.cfg.VendorFormat + "}"
                }
                if pattern == "" {
                    continue
                }
                parsed, err := grokparse.ParseLine(pattern, body)
                if err != nil {
                    continue
                }
                if p.cfg.ExtractCEF {
                    if ext, ok := parsed["cef_ext"]; ok {
                        extMap := grokparse.ParseCEFExtension(ext)
                        for k, v := range extMap {
                            parsed[k] = v
                        }
                    }
                }
                mapped := grokparse.MapFields(parsed, p.cfg.FieldMap)
                for k, v := range mapped {
                    lr.Attributes().PutStr(k, v)
                }
            }
        }
    }
    return logs, nil
}

// grokProcProcessor implements the OTEL Collector processor.Logs interface.
type grokProcProcessor struct {
    next consumer.Logs
    proc *GrokProcessor
}

func (p *grokProcProcessor) Capabilities() consumer.Capabilities {
    return consumer.Capabilities{MutatesData: true}
}

func (p *grokProcProcessor) Start(ctx context.Context, host component.Host) error {
    return nil
}

func (p *grokProcProcessor) Shutdown(ctx context.Context) error {
    return nil
}

// ConsumeLogs is called by the collector pipeline to process logs.
func (p *grokProcProcessor) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
    outLogs, err := p.proc.ProcessLogs(ctx, logs)
    if err != nil {
        return err
    }
    return p.next.ConsumeLogs(ctx, outLogs)
}
```
