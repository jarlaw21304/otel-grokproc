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

    // Step 2: Print expanded regexes for listed Grok pattern names
    patternNames := []string{
        "ASA302014",
        "ASA302013",
        "ASA106023", // Add or remove as needed
    }

    for _, pat := range patternNames {
        regex, err := grokparse.GetExpandedRegex(pat)
        if err != nil {
            fmt.Printf("Pattern %-10s : ERROR: %v\n", pat, err)
        } else {
            fmt.Printf("Pattern %-10s : Expanded regex:\n%s\n\n", pat, regex)
        }
    }

    // Step 3: Try parsing an example log line
    logLine := "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message"

    match, fields, err := grokparse.MatchPattern("ASA302014", logLine)
    if err != nil {
        fmt.Printf("Error parsing log: %v\n", err)
    } else if match {
        fmt.Printf("Log line matched pattern ASA302014. Fields:\n")
        for k, v := range fields {
            fmt.Printf("  %-10s: %v\n", k, v)
        }
    } else {
        fmt.Println("Log line did not match pattern ASA302014.")
    }
}
