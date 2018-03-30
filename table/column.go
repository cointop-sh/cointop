package table

// FormatFn format function
type FormatFn func(interface{}) string

// Col struct
type Col struct {
	name         string
	hide         bool
	format       string
	formatFn     FormatFn
	align        Align
	width        int
	perc         float32
	minWidth     int
	minWidthPerc int
}

// Cols columns
type Cols []*Col

// Hide hide
func (c *Col) Hide() *Col {
	c.hide = true
	return c
}

// SetFormatFn set format function
func (c *Col) SetFormatFn(f FormatFn) *Col {
	c.formatFn = f
	return c
}

// SetFormat sets format
func (c *Col) SetFormat(f string) *Col {
	c.format = f
	return c
}

// AlignLeft align left
func (c *Col) AlignLeft() *Col {
	c.align = AlignLeft
	return c
}

// AlignRight align right
func (c *Col) AlignRight() *Col {
	c.align = AlignRight
	return c
}

// AlignCenter align center
func (c *Col) AlignCenter() *Col {
	c.align = AlignCenter
	return c
}

// SetWidth set width
func (c *Col) SetWidth(w int) *Col {
	c.minWidth = w
	return c
}

// SetWidthPerc  set width percentage
func (c *Col) SetWidthPerc(w int) *Col {
	c.minWidthPerc = w
	return c
}

// Index index
func (c Cols) Index(n string) int {
	for i := range c {
		if c[i].name == n {
			return i
		}
	}
	return -1
}
