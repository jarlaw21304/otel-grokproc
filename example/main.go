package main

import (
	"fmt"
	"log"

	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
	// Step 1: Load patterns
	err := grokparse.LoadAllPatternFiles("patterns")
	if err != nil {
		log.Fatalf("Failed to load grok patterns: %v", err)
	}
	fmt.Println("Grok patterns loaded successfully.")

	// Optional: Show loaded pattern names
	fmt.Println("Loaded pattern names:")
	for k := range grokparse.Patterns {
		fmt.Println("  ", k)
	}
	fmt.Println()

	// Optional: Show expanded regex
	regex, err := grokparse.GetExpandedRegex("ASA302014")
	if err != nil {
		fmt.Printf("Pattern ASA302014: ERROR: %v\n", err)
	} else {
		fmt.Printf("Pattern ASA302014: Expanded regex:\n%s\n\n", regex)
	}

	// Step 2: Try parsing multiple log lines with multiple patterns
	logLines := []string{
		"(TCP): 192.168.1.1/1234 -> 10.1.1.1/80 Some test message",
		"2024-05-18T12:00:00 192.168.2.5 10.0.0.25 Another log format",
		"[10001] login_success user=jane",
	}

	patternNames := []string{"ASA302014", "SYSLOGPARSER", "CEF"}

	for _, logLine := range logLines {
		matched := false
		for _, pattern := range patternNames {
			fields, err := grokparse.ParseLine(pattern, logLine)
			if err == nil {
				fmt.Printf("\n[MATCH] Pattern: %s\n  LogLine: %s\n  Fields:\n", pattern, logLine)
				for k, v := range fields {
					fmt.Printf("    %-12s: %v\n", k, v)
				}
				matched = true
				break // Matched this line; move to next
			}
		}
		if !matched {
			fmt.Printf("\n[NO MATCH] %s\n", logLine)
		}
	}
}
