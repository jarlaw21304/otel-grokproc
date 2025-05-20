package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RegexPattern struct {
	File       string   `json:"file"`
	Line       string   `json:"line"`
	Regex      string   `json:"regex"`
	GroupNames []string `json:"group_names"`
}

func regexify(line string) (string, []string) {
	groupNames := []string{}
	i := 0

	// Named group pattern builder
	named := func(name, pat string) string {
		groupNames = append(groupNames, name)
		return fmt.Sprintf("(?P<%s>%s)", name, pat)
	}

	// Date (date only)
	line = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("date%d", i), `\d{4}-\d{2}-\d{2}`)
	})

	// Time (with comma ms)
	line = regexp.MustCompile(`\d{2}:\d{2}:\d{2},\d{3}`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("time%d", i), `\d{2}:\d{2}:\d{2},\d{3}`)
	})

	// Time (with dot ms, e.g., RFC3339 nano)
	line = regexp.MustCompile(`\d{2}:\d{2}:\d{2}\.\d+`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("timeDot%d", i), `\d{2}:\d{2}:\d{2}\.\d+`)
	})

	// Bracketed text (e.g. [WebContainer : 5])
	line = regexp.MustCompile(`
$$
[^
$$
]+\]`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("bracket%d", i), `[^\]]+`)
	})

	// Loglevel (uppercase word)
	line = regexp.MustCompile(`\b(INFO|WARN|DEBUG|ERROR|TRACE|FATAL|SEVERE|NOTICE|CRITICAL)\b`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("level%d", i), `[A-Z]+`)
	})

	// Classnames, hosts, or dotted identifiers
	line = regexp.MustCompile(`[a-zA-Z0-9_]+\.[\w\.]+`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("dotted%d", i), `[\w\.]+`)
	})

	// Parenthesis content
	line = regexp.MustCompile(`$[^$]*\)`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("paren%d", i), `.*?`)
	})

	// IP addresses (v4)
	line = regexp.MustCompile(`\b\d{1,3}(?:\.\d{1,3}){3}\b`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("ip%d", i), `\d{1,3}(?:\.\d{1,3}){3}`)
	})

	// Integers
	line = regexp.MustCompile(`\b\d+\b`).ReplaceAllStringFunc(line, func(_ string) string {
		i++
		return named(fmt.Sprintf("int%d", i), `\d+`)
	})

	// Leave other text as-is, but escape regex special chars, except inside groups!
	specials := `. * ? ^ $ ( ) = ! < > : -`
	for _, s := range strings.Split(specials, " ") {
		if s != "" {
			line = strings.ReplaceAll(line, s, `\`+s)
		}
	}
	return "^" + line + "$", groupNames
}

func main() {
	logDir := "./logs"
	var patterns []RegexPattern

	files, _ := os.ReadDir(logDir)
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".log") || strings.HasSuffix(file.Name(), ".txt")) {
			f, err := os.Open(filepath.Join(logDir, file.Name()))
			if err != nil {
				continue
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				regex, groupNames := regexify(line)
				patterns = append(patterns, RegexPattern{
					File:       file.Name(),
					Line:       line,
					Regex:      regex,
					GroupNames: groupNames,
				})
			}
		}
	}

	outfile, err := os.Create("log_regexes.json")
	if err != nil {
		fmt.Println("Failed to create output:", err)
		return
	}
	defer outfile.Close()

	enc := json.NewEncoder(outfile)
	enc.SetIndent("", "  ")
	enc.Encode(patterns)
	fmt.Println("Exported patterns to log_regexes.json")
}
