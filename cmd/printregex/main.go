package main

import (
	"fmt"
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: printregex <pattern-dir> <grok-pattern>")
		os.Exit(1)
	}
	patternDir := os.Args[1]
	pattern := os.Args[2]

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
	fmt.Println(expanded)
}
