package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vjeantet/grok"
)

func main() {
	const patternDir = "./patterns"
	const outputFile = "grok_patterns_export.csv"

	grk, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Grok init error:", err)
		os.Exit(1)
	}

	files, err := os.ReadDir(patternDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not open 'patterns' dir:", err)
		os.Exit(1)
	}
	for _, f := range files {
		grk.AddPatternsFromFile(filepath.Join(patternDir, f.Name()))
	}

	// Find all patterns (by name, not contents!)
	patterns := make(map[string]string)
	re := regexp.MustCompile(`^([A-Z0-9_]+)\s+(.+)$`)
	for _, f := range files {
		if f.IsDir() || (!strings.HasSuffix(f.Name(), ".grok") && !strings.HasSuffix(f.Name(), ".patterns")) {
			continue
		}
		fp, err := os.Open(filepath.Join(patternDir, f.Name()))
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", f.Name(), err)
			continue
		}
		defer fp.Close()
		scan := NewScanner(fp)
		for scan.Scan() {
			line := strings.TrimSpace(scan.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			m := re.FindStringSubmatch(line)
			if len(m) == 3 {
				name := m[1]
				raw := m[2]
				patterns[name] = raw
			}
		}
	}

	// Prepare output
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot create output file:", err)
		os.Exit(1)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()
	writer.Write([]string{"GrokPattern", "ExpandedRegex"})

	for patternName := range patterns {
		grokPattern := "%{" + patternName + "}"
		regex, err := grk.Compile(grokPattern, false, false)
		rx := ""
		if err != nil {
			rx = "ERROR: " + err.Error()
		} else {
			rx = regex.String()
		}
		writer.Write([]string{patternName, rx})
	}
	fmt.Printf("Exported %d patterns to %s\n", len(patterns), outputFile)
}

// NewScanner ensures bufio.NewScanner is assigned (to avoid shadowing earlier in closures)
func NewScanner(f *os.File) *bufio.Scanner {
	return bufio.NewScanner(f)
}

