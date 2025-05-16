package grokparse

import (
    "bufio"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

// LoadAllPatternFiles reads *.grok pattern files from a directory.
// Each file may contain one or more lines: NAME GROK_PATTERN
func LoadAllPatternFiles(patternDir string) error {
    files, err := ioutil.ReadDir(patternDir)
    if err != nil {
        return err
    }
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".grok") {
            continue
        }
        fullPath := filepath.Join(patternDir, file.Name())
        if err := loadPatternFile(fullPath); err != nil {
            return err
        }
    }
    return nil
}

func loadPatternFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }
        idx := strings.Index(line, " ")
        if idx > 0 {
            name := line[:idx]
            pat := strings.TrimSpace(line[idx+1:])
            Patterns[name] = pat
        }
    }
    return scanner.Err()
}

