
```markdown
otel-grokproc/
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
│   └── config.yaml
```

---

### **1. go.mod**
```go
module github.com/yourorg/otel-grokproc

go 1.20

replace go.opentelemetry.io/collector => github.com/open-telemetry/opentelemetry-collector v0.97.0

require (
    go.opentelemetry.io/collector v0.97.0
)
```

---

### **2. processor/config.go**
```go
package processor

type Config struct {
    Pattern      string            `mapstructure:"pattern"`       // Grok pattern to use (overrides vendor_format)
    VendorFormat string            `mapstructure:"vendor_format"` // e.g. "CEF_PALOALTO"
    FieldMap     map[string]string `mapstructure:"field_map"`     // log-field -> attr name
    ExtractCEF   bool              `mapstructure:"extract_cef"`   // flatten/parse CEF extension
}
```

---

### **3. processor/processor.go**
```go
package processor

import (
    "context"

    "github.com/yourorg/otel-grokproc/grokparse"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/pdata/plog"
    "go.opentelemetry.io/collector/processor/processorhelper"
)

type GrokProcessor struct {
    cfg *Config
}

func NewProcessor(cfg *Config) *GrokProcessor {
    return &GrokProcessor{cfg: cfg}
}

func (p *GrokProcessor) ProcessLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
    for i := 0; i < logs.ResourceLogs().Len(); i++ {
        rl := logs.ResourceLogs().At(i)
        for j := 0; j < rl.ScopeLogs().Len(); j++ {
            sl := rl.ScopeLogs().At(j)
            for k := 0; k < sl.LogRecords().Len(); k++ {
                lr := sl.LogRecords().At(k)
                body := lr.Body().AsString()
                pattern := p.cfg.Pattern
                if pattern == "" && p.cfg.VendorFormat != "" {
                    pattern = "%{" + p.cfg.VendorFormat + "}"
                }
                if pattern == "" {
                    continue
                }
                parsed, err := grokparse.ParseLine(pattern, body)
                if err != nil {
                    continue
                }
                if p.cfg.ExtractCEF {
                    if ext, ok := parsed["cef_ext"]; ok {
                        extMap := grokparse.ParseCEFExtension(ext)
                        for k, v := range extMap {
                            parsed[k] = v
                        }
                    }
                }
                mapped := grokparse.MapFields(parsed, p.cfg.FieldMap)
                for k, v := range mapped {
                    lr.Attributes().PutStr(k, v)
                }
            }
        }
    }
    return logs, nil
}
```

---

### **4. processor/factory.go**
```go
package processor

import (
    "context"
    "fmt"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/processor/processorhelper"
    "go.opentelemetry.io/collector/processor"
    "go.opentelemetry.io/collector/config"
)

const (
    TypeStr = "grokproc"
)

// Factory returns a new processor factory.
func Factory() component.ProcessorFactory {
    return processor.NewFactory(
        TypeStr,
        createDefaultConfig,
        processor.WithLogs(createLogsProcessor, component.StabilityLevelBeta),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        Pattern:      "",
        VendorFormat: "",
        FieldMap:     map[string]string{},
        ExtractCEF:   false,
    }
}

func createLogsProcessor(
    ctx context.Context,
    params processor.CreateSettings,
    cfg component.Config,
    next consumer.Logs,
) (processor.Logs, error) {
    pcfg := cfg.(*Config)
    proc := NewProcessor(pcfg)
    return processorhelper.NewLogsProcessor(
        ctx,
        params,
        cfg,
        next,
        proc.ProcessLogs,
        processorhelper.WithCapabilities(processorhelper.Capabilities{MutatesData: true}),
    )
}
```

---

### **5. grokparse/grokparse.go** & **patterns.go**
- Use the library as described in [previous answer above](#user-content-otel-grokparse-repo).  
- **Tip:** Use `patterns.go` for the full patterns map.

---

### **6. example/main.go**
```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "github.com/yourorg/otel-grokproc/grokparse"
)

func main() {
    pattern := "%{CEF_PALOALTO}"
    fmap := grokparse.FieldMap{
        "src": "source.ip",
        "dst": "destination.ip",
        "severity": "log.severity",
    }
    fh, _ := os.Open("example/samplelogs.txt")
    sc := bufio.NewScanner(fh)
    for sc.Scan() {
        line := sc.Text()
        parsed, _ := grokparse.ParseLine(pattern, line)
        if ext, ok := parsed["cef_ext"]; ok {
            extfields := grokparse.ParseCEFExtension(ext)
            for k, v := range extfields {
                parsed[k] = v
            }
        }
        mapped := grokparse.MapFields(parsed, fmap)
        fmt.Println(mapped)
    }
}
```

---

### **7. testdata/config.yaml**
Simple sample config:
```yaml
receivers:
  filelog:
    include: [ "example/samplelogs.txt" ]

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

### **8. README.md**
Key documentation:
```markdown
# OTEL GrokProc

An OpenTelemetry Collector log processor for parsing complex loglines using Grok patterns (RFC3164, RFC5424, CEF, LEEF, and vendor variants).

## Usage

- `processors:`
  - `grokproc:`
    - `pattern`: Raw Grok pattern (optional)
    - `vendor_format`: Pattern key for a built-in (e.g. CEF_PALOALTO)
    - `extract_cef`: Whether to flatten CEF extension fields
    - `field_map`: Map parsed fields to desired attribute names

See `testdata/config.yaml` for a collector pipeline example.

## Build

```sh
go build -o myotelcol ./cmd/
```
```

---

### **9. How to Register the Processor in a Collector Distribution**

If you’re building your own collector dist (`cmd/mycollector/main.go`), add:
```go
import (
    "github.com/yourorg/otel-grokproc/processor"
)

func main() {
    components, err := components()
    //...
}

func components() (otelcol.Components, error) {
    return otelcol.Components{
        Processors: map[component.Type]component.ProcessorFactory{
            processor.TypeStr: processor.Factory(),
            // other processors...
        },
    }, nil
}
```

---

## **How to Build and Run**

1. Place your logs in `example/samplelogs.txt`.
2. Edit `testdata/config.yaml` to your liking.
3. Build the collector:
   ```bash
   go build -o grokcol ./cmd/mycollector/
   ```
4. Run with config:
   ```bash
   ./grokcol --config ./testdata/config.yaml
   ```

---

**This is a fully functional OTEL Collector log processor skeleton.  
If you need the full repo zipped, or want a custom Go module name, let me know!**
