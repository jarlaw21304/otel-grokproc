package processor

import (
    "otel-grokproc/processor/grokparse"
    "fmt"
)

// Processor processes logs with grok patterns.
type Processor struct {
    PatternDir string
}

// NewProcessor creates a processor with dynamic grok pattern loading.
func NewProcessor(patternDir string) (*Processor, error) {
    err := grokparse.LoadAllPatternFiles(patternDir)
    if err != nil {
        return nil, err
    }
    return &Processor{PatternDir: patternDir}, nil
}

// ProcessLogLine parses one log line using the pattern name (e.g. "ASA302014").
func (p *Processor) ProcessLogLine(patternName, logline string) (map[string]string, error) {
    grokPattern := "%{" + patternName + "}"
    fields, err := grokparse.ParseLine(grokPattern, logline)
    if err != nil {
        return nil, fmt.Errorf("Grok parse failed for pattern %s: %w", patternName, err)
    }
    return fields, nil
}

