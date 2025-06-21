package tgmd

import "github.com/yuin/goldmark/util"

var Config = &config{
	headings: [6]Element{
		{
			Style:  BoldTg,
			Prefix: "",
		},
		{
			Style:  BoldTg,
			Prefix: "",
		},
		{
			Style:  BoldTg,
			Prefix: "",
		},
		{
			Style:  ItalicsTg,
			Prefix: "",
		},
		{
			Style:  ItalicsTg,
			Prefix: "",
		},
		{
			Style:  ItalicsTg,
			Prefix: "",
		},
	},
	listBullets: [3]rune{
		CircleSymbol.Rune(),
		SquareSymbol.Rune(),
		TriangleSymbol.Rune(),
	},
	Quote: QuoteConfig{
		Enable:     false,
		Expandable: false,
	},
}

type config struct {
	headings    [6]Element
	listBullets [3]rune
	// Quote holds configuration for the document quoting feature.
	Quote QuoteConfig
}

// QuoteConfig holds configuration for the document quoting feature.
type QuoteConfig struct {
	// Enable determines whether the document quoting feature is enabled.
	Enable bool
	// Expandable determines whether the expandable quote feature is enabled.
	Expandable bool
}

// UpdateHeading1 change default H1 style.
func (c *config) UpdateHeading1(e Element) {
	c.headings[0] = e
}

// UpdateHeading2 change default H2 style.
func (c *config) UpdateHeading2(e Element) {
	c.headings[1] = e
}

// UpdateHeading3 change default H3 style.
func (c *config) UpdateHeading3(e Element) {
	c.headings[2] = e
}

// UpdateHeading4 change default H4 style.
func (c *config) UpdateHeading4(e Element) {
	c.headings[3] = e
}

// UpdateHeading5 change default H5 style.
func (c *config) UpdateHeading5(e Element) {
	c.headings[4] = e
}

// UpdateHeading6 change default H6 style.
func (c *config) UpdateHeading6(e Element) {
	c.headings[5] = e
}

// UpdatePrimaryListBullet change default primary bullet.
func (c *config) UpdatePrimaryListBullet(r rune) {
	c.listBullets[0] = r
}

// UpdateSecondaryListBullet change default primary bullet.
func (c *config) UpdateSecondaryListBullet(r rune) {
	c.listBullets[1] = r
}

// UpdateAdditionalListBullet change default primary bullet.
func (c *config) UpdateAdditionalListBullet(r rune) {
	c.listBullets[2] = r
}

// SetQuoteOptions sets the configuration for the document quoting feature.
func (c *config) SetQuoteOptions(q QuoteConfig) {
	c.Quote = q
}

// Element styles object.
type Element struct {
	Style   SpecialTag
	Prefix  string
	Postfix string
}

func (e Element) writeStart(w util.BufWriter) {
	writeSpecialTagStart(w, e.Style, StringToBytes(e.Prefix))
}

func (e Element) writeEnd(w util.BufWriter) {
	writeSpecialTagEnd(w, e.Style, StringToBytes(e.Postfix))
}

// An Option configures a Renderer.
type Option func(*config)

// WithQuote sets the quote options.
func WithQuote(q QuoteConfig) Option {
	return func(c *config) {
		c.SetQuoteOptions(q)
	}
}

// WithHeading1 sets the H1 style.
func WithHeading1(e Element) Option {
	return func(c *config) {
		c.UpdateHeading1(e)
	}
}

// WithHeading2 sets the H2 style.
func WithHeading2(e Element) Option {
	return func(c *config) {
		c.UpdateHeading2(e)
	}
}

// WithHeading3 sets the H3 style.
func WithHeading3(e Element) Option {
	return func(c *config) {
		c.UpdateHeading3(e)
	}
}

// WithHeading4 sets the H4 style.
func WithHeading4(e Element) Option {
	return func(c *config) {
		c.UpdateHeading4(e)
	}
}

// WithHeading5 sets the H5 style.
func WithHeading5(e Element) Option {
	return func(c *config) {
		c.UpdateHeading5(e)
	}
}

// WithHeading6 sets the H6 style.
func WithHeading6(e Element) Option {
	return func(c *config) {
		c.UpdateHeading6(e)
	}
}

// WithPrimaryListBullet sets the primary list bullet.
func WithPrimaryListBullet(r rune) Option {
	return func(c *config) {
		c.UpdatePrimaryListBullet(r)
	}
}

// WithSecondaryListBullet sets the secondary list bullet.
func WithSecondaryListBullet(r rune) Option {
	return func(c *config) {
		c.UpdateSecondaryListBullet(r)
	}
}

// WithAdditionalListBullet sets the additional list bullet.
func WithAdditionalListBullet(r rune) Option {
	return func(c *config) {
		c.UpdateAdditionalListBullet(r)
	}
}
