package test

import (
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"path/filepath"
	"testing"
)

func TestASA302014Pattern(t *testing.T) {
	// Load patterns from patterns/ directory (relative from test/)
	patternDir := filepath.Join("..", "patterns")
	err := grokparse.LoadAllPatternFiles(patternDir)
	if err != nil {
		t.Fatalf("Pattern loading failed: %v", err)
	}
	// Pattern here must exist in ASA302014.grok in patterns/
	pattern := "%{ASA302014}"
	sample := "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message"
	fields, err := grokparse.ParseLine(pattern, sample)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if fields["timestamp"] != "2023-05-13T10:15:14" {
		t.Errorf("Got timestamp: %q", fields["timestamp"])
	}
	if fields["src_ip"] != "192.168.1.1" {
		t.Errorf("Got src_ip: %q", fields["src_ip"])
	}
	if fields["dst_ip"] != "10.1.1.1" {
		t.Errorf("Got dst_ip: %q", fields["dst_ip"])
	}
	if fields["msg"] != "Some test message" {
		t.Errorf("Got msg: %q", fields["msg"])
	}
}
