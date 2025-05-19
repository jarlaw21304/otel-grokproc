package main

import (
	"fmt"
	"log"

	// IMPORTANT: Replace the import below with your actual module name from go.mod!
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
	// 1. Load Grok patterns from the patterns directory in the root project
	err := grokparse.LoadAllPatternFiles("patterns")
	if err != nil {
		log.Fatalf("Failed to load grok patterns: %v", err)
	}
	fmt.Println("Grok patterns loaded successfully.")

	// 2. Example: Parse a sample Cisco ASA log using the appropriate pattern name
	pattern := "%{ASA302014}"
	logline := "(TCP): 192.168.1.1/12345 -> 10.1.1.1/53 Some test message"

	// --- DEBUG BLOCK START ---
	expanded, _ := grokparse.GetExpandedRegex(pattern)
	fmt.Printf("\n==DEBUG MATCH ATTEMPT==\n")
	fmt.Printf("Pattern: %s\n", pattern)
	fmt.Printf("Expanded regex: %q\n", expanded)
	fmt.Printf("LogLine bytes: %v\n", []byte(logline))
	fmt.Printf("LogLine string: %q\n", logline)
	fmt.Printf("======================\n")
	// --- DEBUG BLOCK END ---

	fields, err := grokparse.ParseLine(pattern, logline)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	fmt.Println("Parsed fields:")
	for k, v := range fields {
		fmt.Printf("  %s: %s\n", k, v)
	}
}
