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

func LoadAllPatternFiles(patternDir string) (*Grok, error) {
    g, err := grok.New()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize Grok: %w", err)
    }
    files, err := ioutil.ReadDir(patternDir)
    if err != nil {
        return nil, fmt.Errorf("failed to list pattern directory: %w", err)
    }
    found := false
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".grok") { continue }
        if err := g.AddPatternsFromFile(filepath.Join(patternDir, file.Name())); err != nil {
            return nil, fmt.Errorf("failed to load %s: %w", file.Name(), err)
        }
        found = true
    }
    if !found {
        return nil, fmt.Errorf("no .grok pattern files found in %s", patternDir)
    }
    return &Grok{Grok: g}, nil
}
