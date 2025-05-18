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
    // (Assume ASA302014.grok contains: ASA302014 %{TIMESTAMP_ISO8601:timestamp} %{IP:src_ip} %{IP:dst_ip} %{GREEDYDATA:msg})
    pattern := "%{ASA302014}"
    logline := "2023-05-13T10:15:14 192.168.1.1 10.1.1.1 Some test message"

    fields, err := grokparse.ParseLine(pattern, logline)
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }

    fmt.Println("Parsed fields:")
    for k, v := range fields {
        fmt.Printf("  %s: %s\n", k, v)
    }

    // Example: Dynamically choosing the pattern based on extracted or configured Cisco ASA code
    // code := detectCodeFromLogline(logline) // implement detection as needed
    // pattern := fmt.Sprintf("%%{ASA%s}", code)
    // fields, err := grokparse.ParseLine(pattern, logline)
}
