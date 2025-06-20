package main

import (
	"fmt"
	"os"

	tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
	content, err := os.ReadFile("example/source.md")
	if err != nil {
		fmt.Println("Failed to read source file:", err)
		return
	}

	// --- Standard Conversion ---
	fmt.Println("--- Standard Conversion ---")
	// Ensure quoting is disabled for standard conversion
	tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: false, Expandable: false})
	// Customize other configs for this run
	tgmd.Config.UpdatePrimaryListBullet('◦')
	tgmd.Config.UpdateHeading1(tgmd.Element{
		Style:   tgmd.BoldTg,
		Prefix:  "!!!",
		Postfix: "!!!",
	})

	standardOutput, err := tgmd.Convert(content)
	if err != nil {
		fmt.Println("Standard conversion failed:", err)
		return
	}
	fmt.Println(string(standardOutput))

	// --- Document Quoted Conversion ---
	fmt.Println("\n--- Document Quoted Conversion ---")
	// Enable the document quoting feature
	tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, Expandable: false})

	quotedOutput, err := tgmd.Convert(content)
	if err != nil {
		fmt.Println("Quoted conversion failed:", err)
		return
	}
	fmt.Println(string(quotedOutput))

	// --- Expandable Document Quoted Conversion ---
	fmt.Println("\n--- Expandable Document Quoted Conversion ---")
	// Enable the expandable document quoting feature
	tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{Enable: true, Expandable: true})

	expandableOutput, err := tgmd.Convert(content)
	if err != nil {
		fmt.Println("Expandable quoted conversion failed:", err)
		return
	}
	fmt.Println(string(expandableOutput))
}
