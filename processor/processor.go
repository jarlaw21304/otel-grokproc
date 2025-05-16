package processor

import (
    "context"
    "fmt"
    "strings"

    "github.com/yourorg/otel-grokproc/processor/grokparse"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/pdata/plog"
    "go.opentelemetry.io/collector/processor"
)

// GrokProcessor contains config, loader, and parser
type GrokProcessor struct {
    cfg      *Config
    grokInst *grokparse.GrokWrapper
}

// Wrapper helps expose pattern parsing
type GrokWrapper struct {
    Parser *grokparse.Grok
    // You can add caching here if needed for performance
}

// NewProcessor builds the processor, loads patterns
func NewProcessor(cfg *Config) (*GrokProcessor, error) {
    g, err := grokparse.LoadAllPatternFiles(cfg.PatternDirectory)
    if err != nil {
        return nil, fmt.Errorf("failed loading grok patterns: %w", err)
    }
    return &GrokProcessor{
        cfg:      cfg,
        grokInst: &grokparse.GrokWrapper{Parser: g},
    }, nil
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
                if strings.TrimSpace(pattern) == "" {
                    // Skip processing if not configured
                    continue
                }
                // Call the parser
                parsed, err := p.grokInst.Parser.Parse(pattern, body)
                if err != nil {
                    continue // Or log error
                }
                for key, val := range parsed {
                    if dest, ok := p.cfg.FieldMap[key]; ok {
                        lr.Attributes().PutStr(dest, val)
                    } else {
                        lr.Attributes().PutStr(key, val)
                    }
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

func (p *grokProcProcessor) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
    outLogs, err := p.proc.ProcessLogs(ctx, logs)
    if err != nil {
        return err
    }
    return p.next.ConsumeLogs(ctx, outLogs)
}
