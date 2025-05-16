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

// LoadAllPatternFiles loads .grok files from the given directory
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
                if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
                    continue
                }
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

// expandPattern recursively expands Grok patterns.
func expandPattern(pattern string, depth int) (string, error) {
    if depth > 10 {
        return "", errors.New("pattern recursion too deep")
    }
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
        rep := "(?P<" + field + ">" + expPat + ")"
        if field == "" {
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
    re, err := regexp.Compile("^" + expanded + "$")
    if err != nil {
        return nil, fmt.Errorf("pattern compile failed: %w", err)
    }
    matches := re.FindStringSubmatch(line)
    if matches == nil {
        return nil, errors.New("no match")
    }
    result := make(map[string]string)
    for i, name := range re.SubexpNames() {
        if i > 0 && name != "" {
            result[name] = matches[i]
        }
    }
    return result, nil
}
