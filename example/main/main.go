package main

import (
	"fmt"
	"log"
	"otel-grokproc/processor/grokparse"
)

func main() {
	// Step 1: Dynamically load Grok patterns at startup
	err := grokparse.LoadAllPatternFiles("patterns")
	if err != nil {
		log.Fatalf("Failed to load grok patterns: %v", err)
	}
	fmt.Println("Grok patterns loaded successfully.")

	// Step 2: Parse a sample Cisco ASA log using the appropriate pattern name
	pattern := "%{ASA302014}"
	logline := "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message"

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
