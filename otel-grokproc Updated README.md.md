
# OTEL Grok Processor

An **OpenTelemetry Collector** log processor that parses unstructured logs using Grok patterns, including support for CEF, LEEF, RFC3164, RFC5424, and common vendor formats for Palo Alto, Cisco, Juniper, and Fortigate Firewalls. Extracted fields are mapped to your desired OpenTelemetry log attributes for downstream analytics, SIEM, and observability pipelines.

---

## 📂 Project Directory Structure

```
otel-grokproc/
│   go.mod
│   go.sum
│   grokparser.md          # Additional parsing notes or documentation
│   main.go                # Main entry for standalone execution/parsing
│   printregex             # (Binary or helper tool)
│   README.md
│
├── .idea/                 # IDE configuration (can be ignored)
│
├── cmd/
│   └── printregex/
│         main.go          # CLI utility for printing regex (examples/usage)
│         printregex.exe
│
├── example/
│   ├── config.yaml        # Example configuration for OTEL Collector
│   ├── main.go            # Example standalone usage
│   ├── samplelogs.txt     # Sample log lines for testing
│   └── main/
│         main.go
│
├── patterns/
│       asa_302014.grok    # Vendor and general grok patterns
│       asa_patterns_grok.txt
│       syslog.grok
│
├── processor/
│   ├── config.go
│   ├── factory.go
│   ├── otelgrokproc_processor.go
│   ├── processor.go
│   └── grokparse/
│         grokparse.go     # Grok parsing logic
│
├── test/
│       parser_test.go     # Unit tests
│
└── testdata/
        config.yaml        # Test configuration
---

## ✨ Features

- **Parse Syslog, CEF, LEEF, and major vendor formats** (Palo Alto, Cisco, etc) out-of-the-box
- **Field mapping:** Extracted fields can be mapped to [OpenTelemetry semantic attributes](https://opentelemetry.io/docs/reference/specification/logs/attribute-semantics/)
- **Extendable:**  
  Add your own Grok patterns, field maps, or vendor templates for custom log formats
- **High-performance:**  
  All regex-based field extraction; written in pure Go
- **Collector Integration:**  
  Plug directly into any OpenTelemetry Collector pipeline

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
Or build your own OTEL Collector distribution with this processor:
```shell
go build -o grokcol ./cmd/mycollector/
```

---

### 3. Usage as Standalone Parser

Example CLI (see `example/main.go`):

```go
package main

import (
    "fmt"
    "github.com/YOUR_USERNAME/otel-grokproc/grokparse"
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

Add to your Collector config (`testdata/config.yaml`):

```yaml
receivers:
  filelog:
    include: ["example/samplelogs.txt"]

processors:
  grokproc:
    vendor_format: "CEF_PALOALTO"  # Use a built-in pattern (see below)
    extract_cef: true              # Parse and flatten CEF extension fields
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

- **Add new vendor patterns:** Edit `grokparse/patterns.go`
- **Map new log fields:** Tweak `field_map` in your OTEL config to assign any extracted field to any OTel attribute
- **Custom mapping logic:** Extend `MapFields` in `grokparse/grokparse.go` as needed

---

## 🧪 Development & Contributing

Pull requests, issues, and suggestions are always welcome!

1. Fork & clone the repo
2. Create your branch:  
   ```shell
   git checkout -b my-feature
   ```
3. Make changes (add new log samples to `example/samplelogs.txt` if relevant)
4. Run tests and build:
   ```shell
   go test ./...
   ```
5. Submit a PR and describe your use case

---

## 📜 License

Apache-2.0 or (at your preference) MIT

---

## 💡 Acknowledgements

Built on top of [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/).  
Special thanks to the [Logstash Grok pattern community](https://github.com/logstash-plugins/logstash-filter-grok/tree/main/patterns).

---

**Questions, feature requests, or pattern contributions? Open an Issue or Pull Request!**

