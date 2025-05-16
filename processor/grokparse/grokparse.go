package grokparse

import (
	"fmt"
	"regexp"
)

// Patterns is populated by the loader at startup.
var Patterns = map[string]string{}

// expandPattern recursively interpolates pattern references.
func expandPattern(p string, seen map[string]struct{}) string {
	re := regexp.MustCompile(`%{(\w+)}`)
	for {
		m := re.FindStringSubmatch(p)
		if m == nil {
			break
		}
		key := m[1]
		if _, already := seen[key]; already {
			p = re.ReplaceAllString(p, ".*")
			break
		}
		seen[key] = struct{}{}
		sub, ok := Patterns[key]
		if !ok {
			p = re.ReplaceAllString(p, ".*")
		} else {
			sub = expandPattern(sub, seen)
			p = re.ReplaceAllString(p, "("+sub+")")
		}
	}
	return p
}

// CompileGrok compiles a Grok pattern into a Go regexp with named groups.
func CompileGrok(pattern string) (*regexp.Regexp, error) {
	re := regexp.MustCompile(`%{(\w+)(?::([\w_]+))?}`)
	result := re.ReplaceAllStringFunc(pattern, func(s string) string {
		m := re.FindStringSubmatch(s)
		base := expandPattern("%{"+m[1]+"}", map[string]struct{}{})
		if m[2] != "" {
			return fmt.Sprintf("(?P<%s>%s)", m[2], base)
		}
		return fmt.Sprintf("(%s)", base)
	})
	return regexp.Compile(result)
}

// ParseLine parses a log line using the given Grok pattern.
func ParseLine(pattern, logline string) (map[string]string, error) {
	re, err := CompileGrok(pattern)
	if err != nil {
		return nil, err
	}
	match := re.FindStringSubmatch(logline)
	if match == nil {
		return nil, fmt.Errorf("no match")
	}
	out := map[string]string{}
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			out[name] = match[i]
		}
	}
	return out, nil
}

// FieldMap enables mapping field names.
type FieldMap map[string]string

func MapFields(src map[string]string, fmap FieldMap) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if mapped, ok := fmap[k]; ok && mapped != "" {
			out[mapped] = v
		} else {
			out[k] = v
		}
	}
	return out
}
