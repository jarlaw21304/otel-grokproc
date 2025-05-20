package main

import (
	"bufio"
	"fmt"
	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	patternDir := "./patterns"

	// Load all .grok files for expansion (already supported by your loader)
	err := grokparse.LoadAllPatternFiles(patternDir)
	if err != nil {
		fmt.Printf("Error loading pattern files: %v\n", err)
		os.Exit(1)
	}

	// Print regex for all patterns in each .grok file
	files, _ := os.ReadDir(patternDir)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".grok") {
			f, err := os.Open(filepath.Join(patternDir, file.Name()))
			if err != nil {
				continue
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				split := strings.SplitN(line, " ", 2)
				if len(split) > 0 {
					patname := split[0]
					regex, err := grokparse.GetExpandedRegex("%{" + patname + "}")
					if err != nil {
						fmt.Printf("%s: ERROR %v\n", patname, err)
					} else {
						fmt.Printf("%s: %s\n", patname, regex)
					}
				}
			}
		}
	}
}
