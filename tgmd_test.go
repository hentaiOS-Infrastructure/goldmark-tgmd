package tgmd_test

import (
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
			name:     "Standard Blockquote",
			input:    "> BQ",
			expected: ">BQ",
		},
		{
			name:  "Document as Quote",
			input: "Line 1\nLine 2",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false})
			},
			expected: ">Line 1\n>Line 2",
		},
		{
			name:  "Document as Expandable Quote (Forced)",
			input: "Line 1\nLine 2",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 1})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: ">Line 1\n**>Line 2||",
		},
		{
			name:  "Document as Expandable Quote (Threshold Met)",
			input: "Line 1\nLine 2\nLine 3",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 2})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: ">Line 1\n>Line 2\n**>Line 3||",
		},
		{
			name:  "Document as Non-Expandable Quote (Threshold Not Met)",
			input: "Line 1\nLine 2",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 2})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: ">Line 1\n>Line 2",
		},
		{
			name:  "Complex Document as Quote",
			input: "# Title\n\n- Item 1\n- Item 2\n\nSome `code` here.",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false})
			},
			expected: ">*Title*\n>\n>  • Item 1\n>  • Item 2\n>\n>Some `code` here\\.",
		},
		{
			name:  "Complex Document as Expandable Quote",
			input: "# Title\n\n- Item 1\n- Item 2\n\nSome `code` here.",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 4})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: ">*Title*\n>\n>  • Item 1\n>  • Item 2\n**>\n>Some `code` here\\.||",
		},
		{
			name:  "Document with Existing Blockquote as Quote",
			input: "Line 1\n\n> Nested Quote",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false})
			},
			expected: ">Line 1\n>\n>>Nested Quote",
		},
		{
			name:  "Empty Input with Quoting Enabled",
			input: "",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 1})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: "",
		},
		{
			name:  "Whitespace Input with Quoting Enabled",
			input: "   \n\t\n ",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 1})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, ExpandableAfterLines: 0})
			},
			expected: ">   \n**>\t\n> ||",
		},
		{
			name:  "Quote Enabled but Expandable Disabled with Marker",
			input: "Hello\n**",
			setupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true})
			},
			cleanupConfig: func() {
				tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false})
			},
			expected: ">Hello\n>\\*\\*",
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

			got, err := tgmd.Convert([]byte(tc.input))
			if err != nil {
				t.Fatalf("Convert failed: %v", err)
			}

			// Use strict comparison for newlines
			if string(got) != tc.expected {
				t.Errorf(
					"Output mismatch:\nInput:    %q\nExpected: %q\nGot:      %q",
					tc.input,
					tc.expected,
					string(got),
				)
			}
		})
	}
}
