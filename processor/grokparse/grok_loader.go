package grokparse

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"

    "github.com/elastic/go-grok"
)

type Grok struct {
    *grok.Grok
}

// Wrapper struct used above
type GrokWrapper struct {
    Parser *Grok
}

// LoadAllPatternFiles returns a Grok instance loading all .grok patterns in a dir
func LoadAllPatternFiles(patternDir string) (*Grok, error) {
    g, err := grok.New()
    if err != nil {
        return nil, fmt.Errorf("failed to create Grok instance: %w", err)
    }
    files, err := ioutil.ReadDir(patternDir)
    if err != nil {
        return nil, fmt.Errorf("failed to read pattern directory: %w", err)
    }
    loaded := 0
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".grok") {
            continue
        }
        fullPath := filepath.Join(patternDir, file.Name())
        if err := g.AddPatternsFromFile(fullPath); err != nil {
            return nil, fmt.Errorf("failed to load patterns from %s: %w", fullPath, err)
        }
        loaded++
    }
    if loaded == 0 {
        return nil, fmt.Errorf("no .grok files found in %s", patternDir)
    }
    return &Grok{Grok: g}, nil
}
