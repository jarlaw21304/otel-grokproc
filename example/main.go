package main

import (
    "fmt"
    "log"

    "github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
    // Load all pattern files from the "patterns" directory
    err := grokparse.LoadAllPatternFiles("patterns")
    if err != nil {
        log.Fatalf("Failed to load grok patterns: %v", err)
    }
    fmt.Println("Grok patterns loaded successfully.")

    // List patterns you want to debug
    patternNames := []string{
        "ASA302014",
        "ASA302013",
        "ASA106023", // Add/remove as desired
    }

    for _, pat := range patternNames {
        regex, err := grokparse.GetExpandedRegex(pat)
        if err != nil {
            fmt.Printf("Pattern %-10s : ERROR: %v\n", pat, err)
        } else {
            fmt.Printf("Pattern %-10s : Expanded regex:\n%s\n\n", pat, regex)
        }
    }

    // Test a log line against a pattern
    logLine := "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message"

    fields, err := grokparse.ParseLine("ASA302014", logLine)
    if err != nil {
        fmt.Printf("Log line did not match pattern ASA302014: %v\n", err)
    } else {
        fmt.Printf("Log line matched pattern ASA302014. Fields:\n")
        for k, v := range fields {
            fmt.Printf("  %-10s: %v\n", k, v)
        }
    }
}
