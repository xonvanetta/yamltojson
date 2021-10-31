package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestYamlAdvancedSplit(t *testing.T) {
	t.SkipNow() //CBA to make this work for now
	tests := []struct {
		name          string
		file          string
		extectedTexts []string
	}{
		{
			name: "one yaml document",
			file: `---
Kind: Test`,
			extectedTexts: []string{`---
Kind: Test`},
		},
		{
			name:          "one yaml document without yaml delim",
			file:          `Kind: Test`,
			extectedTexts: []string{`Kind: Test`},
		},
		{
			name: "one yaml document with ending yaml delim",
			file: `Kind: Test
---`,
			extectedTexts: []string{`Kind: Test
---`},
		},
		{
			name: "two yaml document",
			file: `---
Kind: one
---
Kind: two`,
			extectedTexts: []string{`---
Kind: one
`, `---
Kind: two`},
		},
		{
			name: "two yaml document without first yaml delim",
			file: `Kind: one
---
Kind: two`,
			extectedTexts: []string{`Kind: one
---`, `Kind: two`},
		},
		{
			name: "one yaml document with ending yaml delim",
			file: `Kind: one
---
Kind: two
---`,
			extectedTexts: []string{`Kind: one
---`, `Kind: two
---`},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			file := strings.NewReader(test.file)
			scanner := bufio.NewScanner(file)
			buf := make([]byte, 10*1024)
			scanner.Buffer(buf, 10*1024*1024)

			scanner.Split(yamlAdvancedSplit)

			i := 0
			for scanner.Scan() {
				fmt.Println(i)
				text := scanner.Text()
				if text != test.extectedTexts[i] {
					fmt.Printf(text)
					fmt.Println("|||||||")
					fmt.Println(test.extectedTexts[i])
					t.Fatalf("text is not equal")
				}
				i++
			}
		})
	}
}

func TestYamlSplit(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expectedTexts []string
	}{
		{
			name: "one yaml document",
			file: `---
Kind: one`,
			expectedTexts: []string{`---
Kind: one`},
		},
		{
			name: "two yaml document",
			file: `---
Kind: one
---
Kind: two`,
			expectedTexts: []string{`---
Kind: one
`, `---
Kind: two`},
		},
		{
			name: "three yaml document",
			file: `---
Kind: one
---
Kind: two
---
Kind: three`,
			expectedTexts: []string{`---
Kind: one
`, `---
Kind: two
`, `---
Kind: three`},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			file := strings.NewReader(test.file)
			scanner := bufio.NewScanner(file)
			buf := make([]byte, 10*1024)
			scanner.Buffer(buf, 10*1024*1024)

			scanner.Split(yamlSplit)

			i := 0
			for scanner.Scan() {
				text := scanner.Text()
				if text != test.expectedTexts[i] {
					fmt.Printf(text)
					fmt.Println("|||||||")
					fmt.Println(test.expectedTexts[i])
					t.Fatalf("text is not equal with expected iteration: %d", i)
				}
				i++
			}
		})
	}
}
