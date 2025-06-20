package tgmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func TestTGMDConvert_VariousCases(t *testing.T) {
	// Define default configurations for restoration
	defaultH1Config := tgmd.Element{Style: tgmd.BoldTg, Prefix: "", Postfix: ""}
	defaultPrimaryListBullet := tgmd.CircleSymbol.Rune()

	// Read source and expected result from files
	sourceMdContent, err := os.ReadFile("example/source.md")
	if err != nil {
		t.Fatalf("Failed to read source.md: %v", err)
	}

	sourceResultMdContent, err := os.ReadFile("example/source-result.md")
	if err != nil {
		t.Fatalf("Failed to read source-result.md: %v", err)
	}
	// Normalize newlines to LF for consistent comparison
	normalizedResult := strings.ReplaceAll(string(sourceResultMdContent), "\r\n", "\n")

	testCases := []struct {
		name          string
		input         string
		setupConfig   func()
		cleanupConfig func()
		expected      string
	}{
		{
			name:     "Single Line User Input",
			input:    "Test Update: AAAA.000000.001",
			expected: "Test Update: AAAA\\.000000\\.001",
		},
		{
			name:     "Simple Single Line",
			input:    "AAAA.BBBB",
			expected: "AAAA\\.BBBB",
		},
		{
			name:  "Multi-paragraph with blank line",
			input: "Para 1.\n\nPara 2.",
			// Paragraph renderer adds \n\n if HasBlankPreviousLines.
			// The blank line in input causes HBL for Para 2.
			// Para 1 is first, no HBL, so no \n\n from para renderer.
			// Output: Para 1.\n\nPara 2. (newlines from source blank line)
			expected: "Para 1\\.\n\nPara 2\\.\n",
		},
		{
			name:     "Heading 1 (Default Config)",
			input:    "# Heading One",
			expected: "*Heading One*",
		},
		{
			name:  "Heading 1 with Custom Config (from example/main.go)",
			input: "# Heading1 🎉",
			setupConfig: func() {
				tgmd.Config.UpdateHeading1(tgmd.Element{
					Style:   tgmd.BoldTg,
					Prefix:  "!!!",
					Postfix: "!!!",
				})
			},
			cleanupConfig: func() {
				tgmd.Config.UpdateHeading1(defaultH1Config)
			},
			expected: "*\\!\\!\\!Heading1 🎉\\!\\!\\!*",
		},
		{
			name:     "Strikethrough in paragraph",
			input:    "~~strike~~",
			expected: "~strike~",
		},
		{
			name:     "Code span in paragraph",
			input:    "text `code` text",
			expected: "text `code` text",
		},
		{
			name:     "Link in paragraph",
			input:    "[goldmark](url)",
			expected: "[goldmark](url)",
		},
		{
			name:     "Blockquote",
			input:    "> BQ",
			expected: ">BQ",
		},
		{
			name:     "Fenced Code Block (as first element)",
			input:    "```go\nfunc main() {}\n```",
			expected: "```go\nfunc main() {}\n```",
		},
		{
			name:     "List Item (as first element)",
			input:    "- Item 1",
			expected: "  • Item 1",
		},
		{
			name:  "Full Example Source Document",
			input: string(sourceMdContent),
			setupConfig: func() {
				tgmd.Config.UpdateHeading1(tgmd.Element{
					Style:   tgmd.BoldTg,
					Prefix:  "!!!",
					Postfix: "!!!",
				})
				tgmd.Config.UpdatePrimaryListBullet('•')
			},
			cleanupConfig: func() {
				tgmd.Config.UpdateHeading1(defaultH1Config)
				tgmd.Config.UpdatePrimaryListBullet(defaultPrimaryListBullet)
			},
			expected: normalizedResult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupConfig != nil {
				tc.setupConfig()
			}
			// Ensure cleanupConfig is called even if the test panics
			if tc.cleanupConfig != nil {
				defer tc.cleanupConfig()
			}

			md := tgmd.TGMD()
			var buf bytes.Buffer
			if err := md.Convert([]byte(tc.input), &buf); err != nil {
				t.Fatalf("Convert failed: %v", err)
			}
			got := buf.String()

			// Use strict comparison for newlines
			if got != tc.expected {
				t.Errorf("Output mismatch:\nInput:    %q\nExpected: %q\nGot:      %q", tc.input, tc.expected, got)
			}
		})
	}
}
