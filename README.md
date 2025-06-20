
# goldmark-tgmd ✨

Fork of github.com/Mad-Pixels/goldmark-tgmd
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24.1-blue.svg)](https://golang.org)

goldmark-tgmd is an extension for the goldmark Markdown parser
that adds support for Telegram-specific Markdown features 🚀.
This library allows developers to render Markdown text according
to Telegram's formatting options, making it easier to create content
for bots 🤖 and applications integrated with Telegram 📱.

## Features 🌟

- Support for Telegram Markdown features including custom entities like hidden text and strikethrough text. 📝
- **Document Quoting**: A special feature to format the entire document as a single blockquote, with an optional expandable version.
- Easy integration with goldmark-based projects. 🔌
- Extensible architecture for further customizations. 🔨

## Getting Started 🚀

### Prerequisites 📋

- Go 1.19 or higher

### Installation 💽

To install goldmark-tgmd, use the following go get command:

```shell
go get github.com/hentaiOS-Infrastructure/goldmark-tgmd
```

### Usage 🛠️

The library provides a convenient `Convert` function that handles all rendering.

```go
package main

import (
   "fmt"
   "os"

   tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
   content, _ := os.ReadFile("./example/source.md")

   // Standard conversion
   output, _ := tgmd.Convert(content)
   fmt.Println(string(output))
}
```

### Document Quoting

To format the entire document as a blockquote, you can enable the feature through the global `Config`. This is useful for creating self-contained, quoted messages.

The feature is controlled by two flags in `tgmd.QuoteConfig`:

- `Enable`: A `bool` that turns the document quoting feature on or off.
- `Expandable`: A `bool` that makes the quote expandable (collapsible) in Telegram.

#### Example Usage

Here is how to use the `SetQuoteOptions` function to configure the quoting behavior:

```go
package main

import (
   "fmt"
   "os"

   tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
   content, _ := os.ReadFile("./example/source.md")

   // 1. Standard Conversion (quoting disabled)
   tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{ Enable: false })
   standardOutput, _ := tgmd.Convert(content)
   fmt.Println("--- Standard ---\n", string(standardOutput))

   // 2. Blockquote Conversion
   tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{ Enable: true, Expandable: false })
   quotedOutput, _ := tgmd.Convert(content)
   fmt.Println("\n--- Quoted ---\n", string(quotedOutput))

   // 3. Expandable Blockquote Conversion
   tgmd.Config.SetQuoteOptions(tgmd.QuoteConfig{ Enable: true, Expandable: true })
   expandableOutput, _ := tgmd.Convert(content)
   fmt.Println("\n--- Expandable Quoted ---\n", string(expandableOutput))
}
```

- When `Enable` is `true` and `Expandable` is `false`, the entire output is converted into a standard blockquote.
- When both `Enable` and `Expandable` are `true`, the output is converted into an expandable blockquote, wrapped with `**` and `||`.

You can try the full [example](./example) to see this in action.

## Contributing

We're open to any new ideas and contributions. We also have some rules and taboos here, so please read this page and our [Code of Conduct](/CODE_OF_CONDUCT.md) carefully.

### I want to report an issue

If you've found an issue and want to report it, please check our [Issues](https://github.com/hentaiOS-Infrastructure/goldmark-tgmd/issues) page.
