package grokparse

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	Patterns = map[string]string{}
	mu       sync.RWMutex
)

// LoadAllPatternFiles loads all .grok files in the specified directory
func LoadAllPatternFiles(patternDir string) error {
	files, err := os.ReadDir(patternDir)
	if err != nil {
		return fmt.Errorf("could not read pattern directory: %w", err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".grok") {
			fullPath := filepath.Join(patternDir, file.Name())
			f, err := os.Open(fullPath)
			if err != nil {
				return fmt.Errorf("could not open %s: %w", fullPath, err)
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				// Split into pattern name and body (by first space)
				parts := strings.SplitN(line, " ", 2)
				if len(parts) != 2 {
					continue
				}
				Patterns[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
			f.Close()
		}
	}
	return nil
}

// expandPattern recursively expands Grok patterns in '%{PATTERN:field}' or '%{PATTERN}' format
func expandPattern(pattern string, depth int) (string, error) {
	if depth > 10 {
		return "", errors.New("pattern recursion too deep")
	}
	// Matches %{PATTERN} or %{PATTERN:field}
	re := regexp.MustCompile(`%{(\w+)(?::(\w+))?}`)
	for {
		loc := re.FindStringIndex(pattern)
		if loc == nil {
			break
		}
		match := pattern[loc[0]:loc[1]]
		submatches := re.FindStringSubmatch(match)
		patName := submatches[1]
		var field string
		if len(submatches) > 2 {
			field = submatches[2]
		}
		mu.RLock()
		pat, ok := Patterns[patName]
		mu.RUnlock()
		if !ok {
			return "", fmt.Errorf("unknown pattern: %s", patName)
		}
		expPat, err := expandPattern(pat, depth+1)
		if err != nil {
			return "", err
		}
		var rep string
		if field != "" {
			rep = "(?P<" + field + ">" + expPat + ")"
		} else {
			rep = "(" + expPat + ")"
		}
		pattern = pattern[:loc[0]] + rep + pattern[loc[1]:]
	}
	return pattern, nil
}

// ParseLine attempts to apply the Grok pattern to the log line and extract fields
func ParseLine(pattern, line string) (map[string]string, error) {
	expanded, err := expandPattern(pattern, 0)
	if err != nil {
		return nil, err
	}
	// Accept optional trailing whitespace (newline, carriage return)
	fmt.Printf("[Regex compile] ^%s\\s*$\n", expanded)
	re, err := regexp.Compile("^" + expanded + "\\s*$")
	if err != nil {
		return nil, fmt.Errorf("pattern compile failed: %w", err)
	}
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("no match")
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}
	return result, nil
}

// Optional: for testing/debugging, expand a pattern without parsing a line
func ExpandPatternForTest(pattern string) (string, error) {
	return expandPattern(pattern, 0)
}

// GetExpandedRegex returns the expanded regex string for a Grok pattern.
func GetExpandedRegex(patternOrName string) (string, error) {
	mu.RLock()
	pat, ok := Patterns[patternOrName]
	mu.RUnlock()
	if ok {
		return expandPattern(pat, 0)
	}
	// If not found, try to expand literal pattern
	return expandPattern(patternOrName, 0)
}
