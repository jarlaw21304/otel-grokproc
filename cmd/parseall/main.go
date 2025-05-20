package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jarlaw21304/otel-grokproc/processor/grokparse"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LogConfig struct {
	Logs []struct {
		File      string   `yaml:"file"`
		ParseWith string   `yaml:"parse_with"`
		Pattern   string   `yaml:"pattern,omitempty"`
		Regex     string   `yaml:"regex,omitempty"`
		Fields    []string `yaml:"fields,omitempty"`
	} `yaml:"logs"`
}

func main() {
	_ = grokparse.LoadAllPatternFiles("patterns")
	data, err := ioutil.ReadFile("collector-config.yaml")
	if err != nil {
		fmt.Println("Could not read config:", err)
		os.Exit(1)
	}
	var config LogConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Println("YAML error:", err)
		os.Exit(1)
	}

	specs := map[string]LogConfig{}
	for _, logSpec := range config.Logs {
		specs[logSpec.File] = config
	}

	files, _ := ioutil.ReadDir("./logs")
	for _, f := range files {
		if f.IsDir() || (!strings.HasSuffix(f.Name(), ".txt") && !strings.HasSuffix(f.Name(), ".log")) {
			continue
		}
		var spec *LogConfig
		for _, s := range config.Logs {
			if s.File == f.Name() {
				spec = &LogConfig{Logs: []struct {
					File      string   "yaml:\"file\""
					ParseWith string   "yaml:\"parse_with\""
					Pattern   string   "yaml:\"pattern,omitempty\""
					Regex     string   "yaml:\"regex,omitempty\""
					Fields    []string "yaml:\"fields,omitempty\""
				}{s}}
				break
			}
		}
		if spec == nil {
			fmt.Printf("No config for file %s\n", f.Name())
			continue
		}
		fmt.Println("---- Parsing", f.Name(), "----")
		fp, _ := os.Open(filepath.Join("./logs", f.Name()))
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for lnum := 1; scanner.Scan(); lnum++ {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			s := spec.Logs[0]
			switch s.ParseWith {
			case "grok":
				fields, err := grokparse.ParseLine(s.Pattern, line)
				if err != nil {
					fmt.Printf("[line %d] Parse error: %v\n", lnum, err)
				} else {
					fmt.Printf("[line %d] %v\n", lnum, fields)
				}
			case "regex":
				re := regexp.MustCompile(s.Regex)
				matches := re.FindStringSubmatch(line)
				out := map[string]string{}
				for i, name := range s.Fields {
					if i+1 < len(matches) {
						out[name] = matches[i+1]
					}
				}
				if len(out) != len(s.Fields) {
					fmt.Printf("[line %d] Regex did not match or field count off\n", lnum)
				} else {
					fmt.Printf("[line %d] %v\n", lnum, out)
				}
			default:
				fmt.Println("Unknown parser type:", s.ParseWith)
			}
		}
	}
}

