
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

The library can be used in two ways: via the simple `Convert` function for quick conversions, or by integrating it as an extension into a `goldmark` instance for more control.

#### Simple Usage

The `Convert` function handles all the rendering and accepts configuration options.

```go
package main

import (
   "fmt"
   "os"

   tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
   content, _ := os.ReadFile("./example/source.md")

   // Standard conversion with custom options
   output, _ := tgmd.Convert(content,
       tgmd.WithHeading1(tgmd.Element{Style: tgmd.BoldTg}),
       tgmd.WithPrimaryListBullet('-'),
   )
   fmt.Println(string(output))
}
```

#### Advanced Usage (with goldmark)

For more complex scenarios, you can use `tgmd` as a `goldmark` extension. This allows you to combine it with other extensions and have full control over the `goldmark` instance.

```go
package main

import (
    "bytes"
    "fmt"
    "os"

    "github.com/yuin/goldmark"
    tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
    content, _ := os.ReadFile("./example/source.md")

    // Create a new goldmark instance with the tgmd renderer
    md := goldmark.New(
        goldmark.WithRenderer(
            tgmd.NewRenderer(
                tgmd.WithQuote(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 20}),
                tgmd.WithHeading1(tgmd.Element{Style: tgmd.ItalicsTg}),
            ),
        ),
        goldmark.WithExtensions(
            tgmd.Strikethroughs,
            tgmd.Hidden,
            tgmd.DoubleSpace,
        ),
    )

    var buf bytes.Buffer
    if err := md.Convert(content, &buf); err != nil {
        panic(err)
    }

    fmt.Println(buf.String())
}
```

### Configuration

Configuration is done via `Option` functions passed to `tgmd.Convert` or `tgmd.NewRenderer`.

- `WithQuote(QuoteConfig)`: Configures document quoting.
- `WithHeading1(Element)` to `WithHeading6(Element)`: Configures heading styles.
- `WithPrimaryListBullet(rune)`, `WithSecondaryListBullet(rune)`, `WithAdditionalListBullet(rune)`: Configures list bullet styles.

### Document Quoting

To format the entire document as a blockquote, use the `WithQuote` option. This is useful for creating self-contained, quoted messages.

The feature is controlled by two fields in `tgmd.QuoteConfig`:

- `Enable`: A `bool` that turns the document quoting feature on or off.
- `ExpandableAfterLines`: An `int` that defines a line-count threshold. If the quote has more lines than this number, it will be made expandable. A value of `0` (the default) disables this feature.

This unique implementation correctly follows the Telegram API's behavior for expandable quotes. Rather than wrapping the entire quote, the expandable marker (`**`) is injected *at the exact line* where the content becomes hidden, and the closing marker (`||`) is appended at the very end.

#### Example Usage

Here is how to use the `WithQuote` option to configure the quoting behavior:

```go
package main

import (
   "fmt"
   "os"

   tgmd "github.com/hentaiOS-Infrastructure/goldmark-tgmd"
)

func main() {
   content, _ := os.ReadFile("./example/source.md")

   // 1. Standard Conversion (quoting disabled by default)
   standardOutput, _ := tgmd.Convert(content)
   fmt.Println("--- Standard ---\n", string(standardOutput))

   // 2. Blockquote Conversion
   quotedOutput, _ := tgmd.Convert(content,
       tgmd.WithQuote(tgmd.QuoteConfig{Enable: true}),
   )
   fmt.Println("\n--- Quoted ---\n", string(quotedOutput))

   // 3. Auto-Expandable Blockquote Conversion
   // The quote will become expandable if the content has more than 20 lines.
   expandableOutput, _ := tgmd.Convert(content,
       tgmd.WithQuote(tgmd.QuoteConfig{Enable: true, ExpandableAfterLines: 20}),
   )
   fmt.Println("\n--- Expandable Quoted ---\n", string(expandableOutput))
}
```

- When `Enable` is `true` and `ExpandableAfterLines` is `0`, the entire output is converted into a standard blockquote.
- When `Enable` is `true` and `ExpandableAfterLines` is greater than `0`, the output will be converted into an expandable blockquote if the line count exceeds the threshold.

You can try the full [example](./example/main.go) to see this in action.

## Contributing

We're open to any new ideas and contributions. We also have some rules and taboos here, so please read this page and our [Code of Conduct](/CODE_OF_CONDUCT.md) carefully.

### I want to report an issue

If you've found an issue and want to report it, please check our [Issues](https://github.com/hentaiOS-Infrastructure/goldmark-tgmd/issues) page.
