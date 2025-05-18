package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
	// Load all patterns from the patterns directory
	err := grokparse.LoadAllPatternFiles("patterns")
	if err != nil {
		log.Fatalf("Failed to load grok patterns: %v", err)
	}
	fmt.Println("Grok patterns loaded successfully.")

	// Show the loaded pattern names
	fmt.Println("Loaded pattern names:")
	for k := range grokparse.Patterns {
		fmt.Println("  ", k)
	}
	fmt.Println()

	// Example log lines to test matching (add/remove as needed)
	logLines := []string{
		"(TCP): 192.168.1.1/1234 -> 10.1.1.1/80 Some test message",
		"2024-05-18T12:00:00 192.168.2.5 10.0.0.25 Another log format",
		"[10001] login_success user=jane",
		"unexpected or unmatched string for debug",
	}

	// If you want to scan from a file instead, uncomment below and comment the hardcoded logLines:
	/*
		f, err := os.Open("logfile.log")
		if err != nil {
			log.Fatalf("Failed to open logfile: %v", err)
		}
		defer f.Close()
		logLines = nil
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			logLines = append(logLines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading logfile: %v", err)
		}
	*/

	patternNames := []string{"ASA302014", "SYSLOGPARSER", "CEF"}

	for _, rawLine := range logLines {
		logLine := strings.TrimSpace(rawLine)
		matched := false
		for _, pattern := range patternNames {
			fields, err := grokparse.ParseLine(pattern, logLine)
			if err == nil {
				fmt.Printf("\n[MATCH] Pattern: %s\n  LogLine: %s\n  Fields:\n", pattern, logLine)
				for k, v := range fields {
					fmt.Printf("    %-12s: %v\n", k, v)
				}
				matched = true
				break // Stop at first successful match
			}
		}
		if !matched {
			fmt.Printf("\n[NO MATCH] LogLine: %s\n", logLine)
			// Debug: show each expanded regex for the attempted patterns
			for _, pattern := range patternNames {
				regex, err := grokparse.GetExpandedRegex(pattern)
				if err != nil {
					fmt.Printf("  Pattern %s: [error getting regex: %v]\n", pattern, err)
				} else {
					fmt.Printf("  Pattern %s regex: %s\n", pattern, regex)
				}
			}
		}
	}
}
