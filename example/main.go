package main

import (
	"fmt"
	"log"

	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
	// Step 1: Load Grok patterns from the "patterns" directory
	err := grokparse.LoadAllPatternFiles("patterns")
	if err != nil {
		log.Fatalf("Failed to load grok patterns: %v", err)
	}
	fmt.Println("Grok patterns loaded successfully.")

	// Optional: Show the loaded patterns
	fmt.Println("Loaded pattern names:")
	for k := range grokparse.Patterns {
		fmt.Println("  ", k)
	}
	fmt.Println()

	// Step 2: Print expanded regex for ASA302014
	regex, err := grokparse.GetExpandedRegex("ASA302014")
	if err != nil {
		fmt.Printf("Pattern ASA302014: ERROR: %v\n", err)
	} else {
		fmt.Printf("Pattern ASA302014: Expanded regex:\n%s\n\n", regex)
	}

	// Step 3: Try parsing a matching log line
	logLine := "(TCP): 192.168.1.1/1234 -> 10.1.1.1/80 Some test message"
	fields, err := grokparse.ParseLine("ASA302014", logLine)
	if err != nil {
		fmt.Printf("Log line did not match pattern ASA302014: %v\n", err)
	} else {
		fmt.Printf("Log line matched pattern ASA302014. Fields:\n")
		for k, v := range fields {
			fmt.Printf("  %-12s: %v\n", k, v)
		}
	}
}
