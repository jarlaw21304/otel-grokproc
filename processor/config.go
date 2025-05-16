package processor

import "go.opentelemetry.io/collector/component"

type Config struct {
    component.Config
    Pattern          string            `mapstructure:"pattern"`             // Grok parsing pattern
    PatternDirectory string            `mapstructure:"pattern_directory"`   // Directory where .grok pattern files are stored
    FieldMap         map[string]string `mapstructure:"field_map"`           // Field mapping
    ExtractCEF       bool              `mapstructure:"extract_cef"`
}
