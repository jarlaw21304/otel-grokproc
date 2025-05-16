package main

import (
    "fmt"
    "os"
    "github.com/jarlaw21304/otel-grokproc/processor/grokparse"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Fprintln(os.Stderr, "Usage: printregex <pattern-dir> <grok-pattern> [not-match-at-start]")
        os.Exit(1)
    }
    patternDir := os.Args[1]
    pattern := os.Args[2]

    // Optional: third argument for "not match" prefix
    notMatchPrefix := ""
    if len(os.Args) > 3 {
        notMatch := os.Args[3]
        if notMatch != "" {
            notMatchPrefix = "^(?!" + notMatch + ")"
        }
    }

    err := grokparse.LoadAllPatternFiles(patternDir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading pattern files: %v\n", err)
        os.Exit(1)
    }

    expanded, err := grokparse.GetExpandedRegex(pattern)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error expanding pattern: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(notMatchPrefix + expanded)
}

