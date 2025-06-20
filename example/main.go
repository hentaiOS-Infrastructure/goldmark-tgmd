package main

import (
	"bytes"
	"fmt"
	"os"

	tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
	"github.com/yuin/goldmark"
)

func main() {
	content, err := os.ReadFile("example/source.md")
	if err != nil {
		fmt.Println("Failed to read source file:", err)
		return
	}

	// --- Standard Conversion ---
	fmt.Println("--- Standard Conversion ---")
	standardOutput, err := tgmd.Convert(content,
		tgmd.WithQuote(tgmd.QuoteConfig{Enable: false, Expandable: false}),
		tgmd.WithPrimaryListBullet('◦'),
		tgmd.WithHeading1(tgmd.Element{
			Style:   tgmd.BoldTg,
			Prefix:  "!!!",
			Postfix: "!!!",
		}),
	)
	if err != nil {
		fmt.Println("Standard conversion failed:", err)
		return
	}
	fmt.Println(string(standardOutput))

	// --- Document Quoted Conversion ---
	fmt.Println("\n--- Document Quoted Conversion ---")
	quotedOutput, err := tgmd.Convert(content,
		tgmd.WithQuote(tgmd.QuoteConfig{Enable: true, Expandable: false}),
	)
	if err != nil {
		fmt.Println("Quoted conversion failed:", err)
		return
	}
	fmt.Println(string(quotedOutput))

	// --- Expandable Document Quoted Conversion ---
	fmt.Println("\n--- Expandable Document Quoted Conversion ---")
	expandableOutput, err := tgmd.Convert(content,
		tgmd.WithQuote(tgmd.QuoteConfig{Enable: true, Expandable: true}),
	)
	if err != nil {
		fmt.Println("Expandable quoted conversion failed:", err)
		return
	}
	fmt.Println(string(expandableOutput))

	// --- Advanced Usage with custom goldmark instance ---
	fmt.Println("\n--- Advanced Usage ---")
	md := goldmark.New(
		goldmark.WithRenderer(
			tgmd.NewRenderer(
				tgmd.WithHeading1(tgmd.Element{Style: tgmd.ItalicsTg}),
			),
		),
		goldmark.WithExtensions(
			tgmd.Strikethroughs,
			tgmd.Hidden,
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		fmt.Println("Advanced usage failed:", err)
		return
	}
	fmt.Println(buf.String())
}
