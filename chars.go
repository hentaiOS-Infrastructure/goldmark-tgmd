package tgmd

import (
	"unsafe"
)

// SpecialRune define custom rune object.
type SpecialRune rune

// Rune from SpecialRune.
func (sr SpecialRune) Rune() rune {
	return rune(sr)
}

// SpecialChar define custom byte object.
type SpecialChar byte

// Byte from SpecialChar.
func (sc SpecialChar) Byte() byte {
	return byte(sc)
}

func (sc SpecialChar) Bytes(num int) []byte {
	bytes := make([]byte, num)
	for i := range num {
		bytes[i] = byte(sc) // Correctly assign the character to the slice elements
	}
	return bytes
}

// Escaped return SpecialChar as escaped byte char.
func (sc SpecialChar) Escaped() []byte {
	return append([]byte{SlashChar.Byte()}, sc.Byte())
}

// SpecialTag define Markdown formatting characters.
type SpecialTag []SpecialChar

// Bytes from SpecialTags.
func (st SpecialTag) Bytes() []byte {
	return unsafe.Slice((*byte)(unsafe.SliceData(st)), len(st))
}

// define characters.
const (
	UnderscoreChar   SpecialChar = '_'
	AsteriskChar     SpecialChar = '*'
	OpenBracketChar  SpecialChar = '['
	CloseBracketChar SpecialChar = ']'
	OpenParenChar    SpecialChar = '('
	CloseParenChar   SpecialChar = ')'
	OpenBraceChar    SpecialChar = '{'
	CloseBraceChar   SpecialChar = '}'
	HashChar         SpecialChar = '#'
	PlusChar         SpecialChar = '+'
	MinusChar        SpecialChar = '-'
	EqualChar        SpecialChar = '='
	DotChar          SpecialChar = '.'
	TildeChar        SpecialChar = '~'
	PipeChar         SpecialChar = '|'
	ExclamationChar  SpecialChar = '!'
	GreaterThanChar  SpecialChar = '>'
	LessThanChar     SpecialChar = '<'
	BackqouteChar    SpecialChar = '`'
	SpaceChar        SpecialChar = ' '
	NewLineChar      SpecialChar = '\n'
	SlashChar        SpecialChar = '\\'
	TabChar          SpecialChar = '\t'
)

// define symbols.
const (
	CircleSymbol   SpecialRune = '•'
	TriangleSymbol SpecialRune = '⁃'
	SquareSymbol   SpecialRune = '‣'
)

// define Telegram Markdown formatting tags.
var (
	BoldTg          SpecialTag = []SpecialChar{AsteriskChar}
	StrikethroughTg SpecialTag = []SpecialChar{TildeChar}
	UnderlineTg     SpecialTag = []SpecialChar{UnderscoreChar, UnderscoreChar}
	HiddenTg        SpecialTag = []SpecialChar{PipeChar, PipeChar}
	ItalicsTg       SpecialTag = []SpecialChar{UnderscoreChar}
	CodeTg          SpecialTag = []SpecialChar{BackqouteChar, BackqouteChar, BackqouteChar}
	SpanTg          SpecialTag = []SpecialChar{BackqouteChar}
)

// define escape map.
var escape = map[byte][]byte{
	UnderscoreChar.Byte():   UnderscoreChar.Escaped(),
	AsteriskChar.Byte():     AsteriskChar.Escaped(),
	OpenBracketChar.Byte():  OpenBracketChar.Escaped(),
	CloseBracketChar.Byte(): CloseBracketChar.Escaped(),
	OpenParenChar.Byte():    OpenParenChar.Escaped(),
	CloseParenChar.Byte():   CloseParenChar.Escaped(),
	OpenBraceChar.Byte():    OpenBraceChar.Escaped(),
	CloseBraceChar.Byte():   CloseBraceChar.Escaped(),
	HashChar.Byte():         HashChar.Escaped(),
	PlusChar.Byte():         PlusChar.Escaped(),
	MinusChar.Byte():        MinusChar.Escaped(),
	EqualChar.Byte():        EqualChar.Escaped(),
	DotChar.Byte():          DotChar.Escaped(),
	ExclamationChar.Byte():  ExclamationChar.Escaped(),
	GreaterThanChar.Byte():  GreaterThanChar.Escaped(),
	LessThanChar.Byte():     LessThanChar.Escaped(),
	TildeChar.Byte():        TildeChar.Escaped(),
	PipeChar.Byte():         PipeChar.Escaped(),
	BackqouteChar.Byte():    BackqouteChar.Escaped(),
	SlashChar.Byte():        SlashChar.Escaped(),
}
