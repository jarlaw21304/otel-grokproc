package main

import (
    "fmt"
    "log"
    "otel-grokproc/processor/grokparse"
)

func main() {
    err := grokparse.LoadAllPatternFiles("patterns")
    if err != nil {
        log.Fatalf("Failed to load grok patterns: %v", err)
    }
    fmt.Println("Grok patterns loaded successfully.")

    asaCodes := []string{"302014", "302013", "106023"} // add all the ASA codes you have

    // Replace this with test log lines for each type
    loglines := []string{
        "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message",
        // ... add more sample logs, one per code
    }

    for i, logline := range loglines {
        matched := false
        for _, code := range asaCodes {
            pattern := fmt.Sprintf("%%{ASA%s}", code)
            fields, err := grokparse.ParseLine(pattern, logline)
            if err == nil && len(fields) > 0 {
                fmt.Printf("Log #%d: Pattern ASA%s matched. Fields:\n", i, code)
                for k, v := range fields {
                    fmt.Printf("  %s: %s\n", k, v)
                }
                matched = true
                break
            }
        }
        if !matched {
            fmt.Printf("Log #%d: No pattern matched!\n", i)
        }
    }
}

