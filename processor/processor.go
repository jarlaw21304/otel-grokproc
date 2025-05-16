package processor

import (
    "context"
    "github.com/yourorg/otel-grokproc/grokparse"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/pdata/plog"
    "go.opentelemetry.io/collector/processor/processorhelper"
)

type GrokProcessor struct {
    cfg *Config
}

func NewProcessor(cfg *Config) *GrokProcessor {
    return &GrokProcessor{cfg: cfg}
}

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
                if pattern == "" { continue }
                parsed, err := grokparse.ParseLine(pattern, body)
                if err != nil { continue }
                if p.cfg.ExtractCEF {
                    if ext, ok := parsed["cef_ext"]; ok {
                        extMap := grokparse.ParseCEFExtension(ext)
                        for k, v := range extMap { parsed[k] = v }
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
