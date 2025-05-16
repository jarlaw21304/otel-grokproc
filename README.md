# otel-grok-procotel-grokproc/
├── go.mod
├── go.sum
├── README.md
├── grokparse/
│   ├── grokparse.go
│   └── patterns.go
├── processor/
│   ├── config.go
│   ├── processor.go
│   └── factory.go
├── example/
│   ├── main.go
│   └── samplelogs.txt
├── testdata/
│   └── config.yamlessor


# OTEL Grok Processor

An [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) log processor that parses unstructured logs using Grok patterns (including CEF, LEEF, RFC3164, RFC5424, and common vendor formats for Palo Alto, Cisco, Juniper, and Fortigate Firewalls), mapping fields to desired OpenTelemetry log attributes.

---

## Features

- **Parse Syslog, CEF, LEEF, and Vendor-specific formats** (Palo Alto, Cisco, etc) out-of-the-box
- **Map extracted fields** to OpenTelemetry semantic attribute names to support analytics, SIEM, or downstream pipelines
- **Extendable**: Add your own Grok patterns, field maps, or vendor templates to handle any custom log format
- **High-performance**: All field extraction is compiled, regex-based, and written in pure Go
- **Works as a Collector Processor**: Plug into any OTEL collector YAML pipeline

---

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/YOUR_USERNAME/otel-grokproc.git
cd otel-grokproc
```

### Build

```sh
go mod tidy
go build ./...
```

Or, if building your own OTEL Collector distribution with this processor:
```sh
go build -o grokcol ./cmd/mycollector/
```

---

## Usage as Standalone Parser

Example CLI parsing (see `example/main.go`):

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

## Configuration for OTEL Collector

Add to your Collector config (`testdata/config.yaml`):

```yaml
receivers:
  filelog:
    include: [ "example/samplelogs.txt" ]

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

#### Supported Pattern Keys

- `SYSLOG3164` : Standard Unix syslog
- `SYSLOG5424` : Modern RFC5424 syslog
- `LEEF_HEADER` : IBM LEEF header
- `CEF_HEADER` : Generic CEF logs
- `CEF_PALOALTO`, `CEF_CISCO`, `CEF_FORTIGATE`, `CEF_JUNIPER` – Vendor-optimized CEF patterns

Or write your own with full grok syntax in the `pattern:` field.

---

## Extending and Customization

- **Add new vendor patterns**: Edit `grokparse/patterns.go`
- **Map new log fields**: Tweak `field_map` in your config to route any extracted field to any OTel attribute
- **Change attribute data types/complex mapping**: Extend `MapFields` logic in `grokparse/grokparse.go`

---

## Development & Contributing

Pull requests, issues, and suggestions are always welcome!
To contribute:

1. Fork & clone the repo
2. Create a branch (`git checkout -b my-feature`)
3. Make changes, add log samples to `example/samplelogs.txt`
4. Run tests and build
5. PR and describe your use case

---

## License

[Apache-2.0](LICENSE) or (at your preference) MIT.

---

## Acknowledgements

Built on top of [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) – special thanks to Grok pattern inspiration from Logstash and Elastic community!
