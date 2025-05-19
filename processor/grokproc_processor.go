package processor

import (
	"context"
	"fmt"
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
)

type Config struct {
	component.Config
	Pattern          string            `mapstructure:"pattern"`           // Example: "%{MY_PATTERN}"
	PatternDirectory string            `mapstructure:"pattern_directory"` // Where your .grok files are
	FieldMap         map[string]string `mapstructure:"field_map"`         // Field mapping
	ExtractCEF       bool              `mapstructure:"extract_cef"`       // handle CEF ext if needed
}

type grokProcProcessor struct {
	next consumer.Logs
	pcfg *Config
}

func NewProcessor(c *Config) (*grokProcProcessor, error) {
	// Load all patterns at startup
	if err := grokparse.LoadAllPatternFiles(c.PatternDirectory); err != nil {
		return nil, fmt.Errorf("could not load grok patterns: %w", err)
	}
	return &grokProcProcessor{pcfg: c}, nil
}

func (p *grokProcProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
func (p *grokProcProcessor) Start(ctx context.Context, host component.Host) error { return nil }
func (p *grokProcProcessor) Shutdown(ctx context.Context) error                   { return nil }

func (p *grokProcProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	resSlice := ld.ResourceLogs()
	for i := 0; i < resSlice.Len(); i++ {
		scopeSlice := resSlice.At(i).ScopeLogs()
		for j := 0; j < scopeSlice.Len(); j++ {
			logs := scopeSlice.At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				bodyStr := lr.Body().AsString()
				pattern := p.pcfg.Pattern
				if pattern == "" && p.pcfg.FieldMap["pattern_name"] != "" {
					pattern = "%{" + p.pcfg.FieldMap["pattern_name"] + "}"
				}
				if pattern == "" {
					continue // nothing to do
				}
				fields, err := grokparse.ParseLine(pattern, bodyStr)
				if err == nil && fields != nil {
					// Remap fields with FieldMap if provided
					attrs := lr.Attributes()
					for src, val := range fields {
						dst := src
						if len(p.pcfg.FieldMap) > 0 {
							if remap, ok := p.pcfg.FieldMap[src]; ok {
								dst = remap
							}
						}
						attrs.PutStr(dst, val)
					}
					// CEF extension handling optional:
					if p.pcfg.ExtractCEF {
						if ext, ok := fields["cef_ext"]; ok {
							for k, v := range grokparse.ParseCEFExtension(ext) {
								attrs.PutStr(k, v)
							}
						}
					}
				}
			}
		}
	}
	return p.next.ConsumeLogs(ctx, ld)
}
