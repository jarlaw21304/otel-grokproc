package grokparse

var Patterns = map[string]string{
    "YEAR":         `\d{4}`,
    "MONTH":        `[A-Za-z]{3}`,
    "NUMBER":       `\d+`,
    "IPV4":         `\d{1,3}(?:\.\d{1,3}){3}`,
    "HOSTNAME":     `[A-Za-z0-9\.\-]+`,
    "WORD":         `\w+`,
    "SYSLOGTIMESTAMP": `%{MONTH} +%{NUMBER} %{NUMBER}:%{NUMBER}:%{NUMBER}`,
    // ...add all other patterns from previous answers here!
}
