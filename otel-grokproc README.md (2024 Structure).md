
# OTEL Grok Processor

An **OpenTelemetry Collector** log processor that parses unstructured logs using Grok patterns, including support for CEF, LEEF, RFC3164, RFC5424, and common vendor formats for Palo Alto, Cisco, Juniper, and Fortigate Firewalls. Extracted fields are mapped to your desired OpenTelemetry log attributes for downstream analytics, SIEM, and observability pipelines.

---

## 📂 Project Directory Structure

```
otel-grokproc/
│   LICENSE
│   README.md
│   go.mod
│   go.sum
│   collector-config.yaml
│
├── patterns/
│     additional_patterns.grok
│     cisco.asa.grok
│     cisco.ise.grok
│     java.grok
│     linux-sudo.grok
│     postfix.grok
│     rfc5424_patterns.grok
│     tomcat.grok
│
├── processor/
│   └── grokparse/
│         grokparse.go
│         grokparse_test.go
│
├── cmd/
│   ├── grok2regex/
│   │     main.go           # CLI: Expand grok to regex
│   ├── grok2regex-csv/
│   │     main.go           # CLI: Export all patterns to CSV
│   ├── parseall/
│   │     main.go
│   └── rawregexgen/
│         main.go
│
├── test/
│     apache.log
│     asa.log
│     cisco_ios.log
│     cisco_ise.log
│     postfix.log
│
└── .github/
    └── workflows/
          go.yml            # GitHub Actions CI workflow
```

---

## ✨ Features

- **Parse Syslog, CEF, LEEF, and major vendor formats** (Palo Alto, Cisco, etc.) out-of-the-box
- **Field mapping:** Extracted fields can be mapped to [OpenTelemetry semantic attributes](https://opentelemetry.io/docs/reference/specification/logs/attribute-semantics/)
- **Extendable:**  
  Add your own Grok patterns, field maps, or vendor templates for custom log formats
- **High-performance:**  
  All regex-based field extraction; written in pure Go
- **Collector Integration:**  
  Plug directly into any OpenTelemetry Collector pipeline
- **Pattern Tools:**  
  - `grok2regex`: Expand a single Grok pattern into its full regex form  
  - `grok2regex-csv`: Export all loaded Grok patterns as extended regexes in CSV format

---

## 🚀 Getting Started

### 1. Clone the Repository

```shell
git clone https://github.com/YOUR_USERNAME/otel-grokproc.git
cd otel-grokproc
```

### 2. Build

```shell
go mod tidy
go build ./...
```

You can also build your own OTEL Collector distribution with this processor:
```shell
go build -o grokcol ./cmd/mycollector/
```

---

### 3. Convert or Export Grok Patterns

#### a) Interactive: Convert a Grok Pattern to Regex

```sh
go run cmd/grok2regex/main.go
# Enter a Grok pattern name at the prompt (e.g., COMBINEDAPACHELOG)
```

#### b) Batch Export All Patterns to CSV

```sh
go run cmd/grok2regex-csv/main.go
# Outputs grok_patterns_export.csv mapping pattern names to expanded regex strings
```

---

### 4. Usage as Standalone Parser

Example CLI (see your `example/main.go`):

```go
package main

import (
    "fmt"
    "github.com/YOUR_USERNAME/otel-grokproc/processor/grokparse"
)

func main() {
    pattern := "%{CEF_PALOALTO}"
    logline := `CEF:0|Palo Alto Networks|PAN-OS|8.0.0|threat|threat log|1|src=192.0.2.1 dst=198.51.100.1 ...`
    fmap := grokparse.FieldMap{
        "src": "source.ip",
        "dst": "destination.ip",
        "severity": "log.severity",
    }
    parsed, _ := grokparse.ParseLine(pattern, logline)
    if ext, ok := parsed["cef_ext"]; ok {
        extfields := grokparse.ParseCEFExtension(ext)
        for k, v := range extfields {
            parsed[k] = v
        }
    }
    mapped := grokparse.MapFields(parsed, fmap)
    fmt.Println(mapped)
}
```

---

## 🛠️ Configuration for OTEL Collector

Add to your Collector config (`collector-config.yaml`):

```yaml
receivers:
  filelog:
    include: ["test/apache.log"]

processors:
  grokproc:
    vendor_format: "CEF_PALOALTO"
    extract_cef: true
    field_map:
      src: "source.ip"
      dst: "destination.ip"
      severity: "log.severity"

exporters:
  logging:
    loglevel: debug

service:
  pipelines:
    logs:
      receivers: [filelog]
      processors: [grokproc]
      exporters: [logging]
```

---

## 🗝️ Supported Pattern Keys

- `SYSLOG3164` : Standard Unix syslog  
- `SYSLOG5424` : Modern RFC5424 syslog  
- `LEEF_HEADER` : IBM LEEF header  
- `CEF_HEADER` : Generic CEF logs  
- `CEF_PALOALTO`, `CEF_CISCO`, `CEF_FORTIGATE`, `CEF_JUNIPER` – Pre-optimized vendor CEF patterns

Or write your custom patterns using full Grok syntax!

---

## 🔧 Extending and Customization

- **Add new vendor patterns:** Edit or add pattern files under `patterns/`
- **Map new log fields:** Update `field_map` in your OTEL config
- **Custom mapping logic:** Extend `MapFields` in `grokparse.go` as needed

---

## 🤖 Continuous Integration (CI)

- Every PR and `main` branch push triggers GitHub Actions:
  - Build and vet all Go code
  - Run all tests
  - Build CLI tools (`grok2regex` and `grok2regex-csv`)
  - Optionally, export and upload `grok_patterns_export.csv` as a CI artifact for download

---

## 🧪 Development & Contributing

Pull requests, issues, and suggestions are always welcome!

1. Fork & clone the repo  
2. Create your branch:  
   ```shell
   git checkout -b my-feature
   ```
3. Make changes (add new patterns or log samples as needed)
4. Run tests and build:
   ```shell
   go test ./...
   ```
5. Submit a PR and describe your use case

---

## 📜 License

Apache-2.0 or MIT  
See `LICENSE` for details.

---

## 💡 Acknowledgements

Built on top of [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/).  
Inspired by the [Logstash Grok pattern community](https://github.com/logstash-plugins/logstash-filter-grok/tree/main/patterns).

---

**Questions, feature requests, or pattern contributions? Open an Issue or Pull Request!**
