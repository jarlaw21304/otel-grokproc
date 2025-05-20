package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vjeantet/grok"
)

func main() {
	const patternDir = "./patterns"

	grk, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Grok init error:", err)
		os.Exit(1)
	}

	files, err := os.ReadDir(patternDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not open 'patterns' dir:", err)
		os.Exit(1)
	}
	for _, f := range files {
		grk.AddPatternsFromFile(patternDir + "/" + f.Name())
	}

	fmt.Println("Paste Grok pattern name (e.g. COMBINEDAPACHELOG) or %{@}:")
	fmt.Print("> ")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		name := sc.Text()
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Print("> ")
			continue
		}
		var grokPattern string
		if !strings.HasPrefix(name, "%{") {
			grokPattern = "%{" + name + "}"
		} else {
			grokPattern = name
		}
		regex, err := grk.Compile(grokPattern, false, false)
		if err != nil {
			fmt.Printf("Expansion error: %v\n", err)
		} else {
			fmt.Println("Expanded regex:")
			fmt.Println(regex.String())
		}
		fmt.Print("> ")
	}
}

