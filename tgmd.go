package tgmd

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	ext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	textm "github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Convert is a custom function that wraps the standard Goldmark conversion.
// It allows for post-processing to quote the entire document.
func Convert(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	md := TGMD()
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}

	if !Config.Quote.Enable {
		return buf.Bytes(), nil
	}

	contentToProcess := buf.Bytes()
	// If goldmark rendered nothing, but the source is not empty,
	// it's likely whitespace. Use the source directly.
	if buf.Len() == 0 && len(source) > 0 {
		contentToProcess = source
	}

	// Trim trailing newlines, but preserve whitespace content
	processedBuf := bytes.TrimRight(contentToProcess, "\n")

	// If the result is effectively empty, return nothing.
	if len(processedBuf) == 0 {
		return []byte{}, nil
	}

	// Post-processing for QuoteDocument
	var result bytes.Buffer
	if Config.Quote.Expandable {
		result.Write([]byte{'*', '*'})
	}

	lines := bytes.Split(processedBuf, []byte{'\n'})
	for i, line := range lines {
		result.WriteByte(GreaterThanChar.Byte())
		result.Write(line)
		if i < len(lines)-1 {
			result.WriteByte(NewLineChar.Byte())
		}
	}

	if Config.Quote.Expandable {
		result.Write([]byte{'|', '|'})
	}

	return result.Bytes(), nil
}

// TGMD (telegramMarkdown) endpoint.
func TGMD() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(util.Prioritized(NewRenderer(), 1000)),
			),
		),
		goldmark.WithExtensions(Strikethroughs),
		goldmark.WithExtensions(Hidden),
	)
}

// Renderer implement renderer.NodeRenderer object.
type Renderer struct{}

// NewRenderer initialize Renderer as renderer.NodeRenderer.
func NewRenderer() renderer.NodeRenderer {
	return &Renderer{}
}

// RegisterFuncs add AST objects to Renderer.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDocument, r.document)
	reg.Register(ast.KindParagraph, r.paragraph)

	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindEmphasis, r.emphasis)

	reg.Register(ast.KindHeading, r.heading)
	reg.Register(ast.KindList, r.renderList) // Changed r.list to r.renderList
	reg.Register(ast.KindListItem, r.listItem)
	reg.Register(ast.KindLink, r.link)

	reg.Register(ast.KindBlockquote, r.blockquote)
	reg.Register(ast.KindFencedCodeBlock, r.code)
	reg.Register(ast.KindCodeSpan, r.codeSpan)

	reg.Register(ext.KindStrikethrough, r.strikethrough)
	reg.Register(KindHidden, r.hidden)
	reg.Register(KindDoubleSpace, r.doubleSpace)
}

// isEffectivelyEmpty checks if a node is an empty paragraph.
func isEffectivelyEmpty(node ast.Node) bool {
	if node == nil {
		return true // Or false, depending on how nil should be treated. For "first child", nil means no first child.
	}
	if p, ok := node.(*ast.Paragraph); ok && p.ChildCount() == 0 {
		return true
	}
	return false
}

// isFirstVisibleBlock checks if the given node is the first block-level element
// in the document that would produce visible output. It skips over an initial
// empty paragraph (often from a BOM).
func isFirstVisibleBlock(node ast.Node) bool {
	if node == nil || node.Parent() == nil || node.Parent().Kind() != ast.KindDocument {
		return false // Not a direct child of the document
	}
	if node.PreviousSibling() == nil { // It's the first child
		return !isEffectivelyEmpty(node) // True if it's not an empty paragraph itself
	}
	// It's not the first child, check if the actual first child was an empty paragraph
	// and this node is the second child.
	if isEffectivelyEmpty(node.PreviousSibling()) && node.PreviousSibling().PreviousSibling() == nil {
		return !isEffectivelyEmpty(node) // True if this node (the second child) is not empty
	}
	return false
}

// writeBlockSeparationNewLines handles the newline logic for block elements.
func writeBlockSeparationNewLines(w util.BufWriter, n ast.Node) {
	if isFirstVisibleBlock(n) {
		// No leading newlines for the very first visible block
		return
	}
	if n.HasBlankPreviousLines() {
		writeNewLine(w)
		writeNewLine(w)
	} else if n.PreviousSibling() != nil && !isEffectivelyEmpty(n.PreviousSibling()) {
		// Single newline if immediately follows another non-empty block
		writeNewLine(w)
	}
}

func (r *Renderer) heading(w util.BufWriter, _ []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	n := node.(*ast.Heading)
	if entering {
		writeBlockSeparationNewLines(w, n)
		Config.headings[n.Level-1].writeStart(w)
	} else {
		Config.headings[n.Level-1].writeEnd(w)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) paragraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	n := node.(*ast.Paragraph)
	if entering {
		// Rule 0: Skip empty first paragraph in document (BOM handling)
		if n.PreviousSibling() == nil && n.Parent().Kind() == ast.KindDocument && isEffectivelyEmpty(n) {
			return ast.WalkContinue, nil
		}

		parentKind := n.Parent().Kind()

		if parentKind == ast.KindListItem || parentKind == ast.KindBlockquote {
			// Paragraphs inside ListItems or Blockquotes:
			// Only add \n\n if the paragraph itself has HBL (blank line *within* the container).
			// Otherwise, no leading newlines from the paragraph itself. Content flows after bullet/>.
			// Only add \n\n if the paragraph has HBL AND it's not the first child of its container.
			// (n.PreviousSibling() == nil implies it's the first child paragraph within the container)
			if n.HasBlankPreviousLines() && n.PreviousSibling() != nil {
				writeNewLine(w)
				writeNewLine(w)
			}
		} else {
			if !isFirstVisibleBlock(n) { // Not the first visible block in the document
				if n.HasBlankPreviousLines() { // Preceded by blank line(s) in source
					writeNewLine(w)
					writeNewLine(w)
				} else if n.PreviousSibling() != nil && !isEffectivelyEmpty(n.PreviousSibling()) { // Immediately follows another non-empty block
					writeNewLine(w)
				}
			}
		}
	}
	return ast.WalkContinue, nil
}

// renderList handles the ast.KindList node.
// It's responsible for the newlines *before* the entire list block.
func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	if entering {
		n := node.(*ast.List)
		writeBlockSeparationNewLines(w, n)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) listItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	n := node.(*ast.ListItem)
	if entering {
		// Newline separation from previous item
		if n.PreviousSibling() != nil { // If not the first list item in this list
			if n.HasBlankPreviousLines() { // If blank lines were present in source between items
				writeNewLine(w)
				writeNewLine(w)
			} else {
				writeNewLine(w) // Default single newline between items
			}
		}
		// Else (it's the first list item), newlines before the whole list are handled by renderList's writeBlockSeparationNewLines

		// Indentation and bullet logic
		listLevel := -1
		for p := n.Parent(); p != nil; p = p.Parent() {
			if p.Kind() == ast.KindList {
				listLevel++
			}
		}

		bulletIndex := listLevel
		if bulletIndex >= len(Config.listBullets) {
			bulletIndex = len(Config.listBullets) - 1
		}

		indentation := (listLevel * 2) + 2

		writeRowBytes(w, SpaceChar.Bytes(indentation))
		writeRune(w, Config.listBullets[bulletIndex])
		writeRowBytes(w, SpaceChar.Bytes(1)) // Single space after bullet
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) code(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	n := node.(interface {
		Lines() *textm.Segments
	})
	var content []byte
	l := n.Lines().Len()
	for i := range l {
		line := n.Lines().At(i)
		content = append(content, line.Value(source)...)
	}
	content = bytes.ReplaceAll(
		content,
		[]byte{TabChar.Byte()},
		[]byte{SpaceChar.Byte(), SpaceChar.Byte(), SpaceChar.Byte()},
	)
	nn := node.(*ast.FencedCodeBlock)
	if entering {
		writeBlockSeparationNewLines(w, nn)
		writeWrapperArr(w.Write(CodeTg.Bytes()))
		writeWrapperArr(w.Write(nn.Language(source)))
		writeNewLine(w)
	} else {
		writeWrapperArr(w.Write(content))
		writeWrapperArr(w.Write(CodeTg.Bytes()))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	render(w, n.Segment.Value(source))
	if n.SoftLineBreak() || n.HardLineBreak() {
		writeNewLine(w)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	_, _ = w.Write(n.Value)
	return ast.WalkContinue, nil
}

func (r *Renderer) emphasis(w util.BufWriter, _ []byte, node ast.Node, _ bool) (
	ast.WalkStatus, error,
) {
	n := node.(*ast.Emphasis)
	if n.Level == 2 {
		writeRowBytes(w, BoldTg.Bytes())
	}
	if n.Level == 1 {
		writeRowBytes(w, ItalicsTg.Bytes())
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) link(w util.BufWriter, _ []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	n := node.(*ast.Link)
	if entering {
		writeRowBytes(w, []byte{OpenBracketChar.Byte()})
	} else {
		writeRowBytes(w, []byte{CloseBracketChar.Byte(), OpenParenChar.Byte()})
		writeRowBytes(w, n.Destination)
		writeRowBytes(w, []byte{CloseParenChar.Byte()})
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) blockquote(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		writeBlockSeparationNewLines(w, n)
		writeRowBytes(w, []byte{GreaterThanChar.Byte()})
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) codeSpan(w util.BufWriter, _ []byte, _ ast.Node, _ bool) (
	ast.WalkStatus, error,
) {
	writeWrapperArr(w.Write(SpanTg.Bytes()))
	return ast.WalkContinue, nil
}

func (r *Renderer) strikethrough(w util.BufWriter, _ []byte, _ ast.Node, _ bool) (
	ast.WalkStatus, error,
) {
	writeWrapperArr(w.Write(StrikethroughTg.Bytes()))
	return ast.WalkContinue, nil
}

func (r *Renderer) hidden(w util.BufWriter, _ []byte, _ ast.Node, _ bool) (
	ast.WalkStatus, error,
) {
	writeWrapperArr(w.Write(HiddenTg.Bytes()))
	return ast.WalkContinue, nil
}

func (r *Renderer) doubleSpace(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (
	ast.WalkStatus, error,
) {
	return ast.WalkContinue, nil
}

func (r *Renderer) document(w util.BufWriter, _ []byte, node ast.Node, entering bool) (
	ast.WalkStatus, error,
) {
	if entering {
		return ast.WalkContinue, nil
	}

	// Add a final newline for multi-block documents.
	if node.ChildCount() > 1 {
		writeNewLine(w)
	}

	return ast.WalkContinue, nil
}
