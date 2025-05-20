package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func escapeRegexMeta(s string) string {
	regexSpecials := `\.+*?()|[]{}^$`
	var out strings.Builder
	for _, r := range s {
		if strings.ContainsRune(regexSpecials, r) {
			out.WriteByte('\\')
		}
		out.WriteRune(r)
	}
	return out.String()
}

func markupTemplateToRegex(template, sample string) (string, []string, error) {
	var fields []string
	var regex strings.Builder
	var prevEnd, samplePos int
	re := regexp.MustCompile(`{(\w+)}`)
	parts := re.FindAllStringSubmatchIndex(template, -1)
	for _, p := range parts {
		fieldName := template[p[2]:p[3]]
		constStart := prevEnd
		constEnd := p[0]
		constText := template[constStart:constEnd]
		constTextEsc := escapeRegexMeta(constText)
		regex.WriteString(constTextEsc)

		if samplePos+len(constText) > len(sample) || sample[samplePos:samplePos+len(constText)] != constText {
			return "", nil, fmt.Errorf("template mismatch at \"%s\" (use identical fixed chars)", constText)
		}
		samplePos += len(constText)
		nextVarStart := p[1]
		after := ""
		if len(parts) > 0 && p != parts[len(parts)-1] {
			after = template[p[1]:parts[len(parts)-1][0]]
		} else if p[1] < len(template) {
			after = template[p[1]:]
		}
		// Guess pattern for common field types
		pat := ".*?"
		if strings.Contains(strings.ToLower(fieldName), "time") {
			pat = `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3}`
		} else if strings.Contains(strings.ToLower(fieldName), "ip") {
			pat = `[0-9.]+`
		} else if strings.Contains(strings.ToLower(fieldName), "level") {
			pat = `[A-Z]+`
		}
		regex.WriteString(fmt.Sprintf("(?P<%s>%s)", fieldName, pat))
		fields = append(fields, fieldName)

		// Best-effort: look ahead in sample for "after" fixed segment
		if after != "" {
			after = after
			ix := strings.Index(sample[samplePos:], after)
			if ix >= 0 {
				samplePos += ix
			}
		}
		prevEnd = p[1]
	}
	// Write any trailing literal
	if prevEnd < len(template) {
		trail := template[prevEnd:]
		regex.WriteString(escapeRegexMeta(trail))
	}
	return "^" + regex.String() + "$", fields, nil
}

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Paste marked-up log template (e.g., {timestamp} [{thread}] {level} {logger} - {message}):")
	tmpl, _ := r.ReadString('\n')
	tmpl = strings.TrimSpace(tmpl)
	fmt.Println("Paste concrete log line matching template:")
	sample, _ := r.ReadString('\n')
	sample = strings.TrimSpace(sample)

	regex, fields, err := markupTemplateToRegex(tmpl, sample)
	if err != nil {
		fmt.Println("Could not generate regex:", err)
		os.Exit(1)
	}
	fmt.Println("Go regex:", regex)
	fmt.Println("Fields:", fields)

	rx := regexp.MustCompile(regex)
	match := rx.FindStringSubmatch(sample)
	if len(match) == len(fields)+1 {
		fmt.Println("Sample Extraction:")
		for i, f := range fields {
			fmt.Printf("  %s: %s\n", f, match[i+1])
		}
	} else {
		fmt.Println("Warning: regex did not yield all fields")
	}
}
