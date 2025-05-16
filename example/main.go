package main

import (
    "bufio"
    "fmt"
    "os"
    "github.com/yourorg/otel-grokproc/grokparse"
)

func main() {
    pattern := "%{CEF_PALOALTO}"
    fmap := grokparse.FieldMap{"src":"source.ip", "dst":"destination.ip", "severity":"log.severity"}
    fh, _ := os.Open("example/samplelogs.txt")
    sc := bufio.NewScanner(fh)
    for sc.Scan() {
        line := sc.Text()
        parsed, _ := grokparse.ParseLine(pattern, line)
        if ext, ok := parsed["cef_ext"]; ok {
            extfields := grokparse.ParseCEFExtension(ext)
            for k, v := range extfields { parsed[k] = v }
        }
        mapped := grokparse.MapFields(parsed, fmap)
        fmt.Println(mapped)
    }
}
