package cointop

import (
	"strconv"

	fcolor "github.com/fatih/color"
	gocui "github.com/miguelmota/gocui"
	xtermcolor "github.com/tomnomnom/xtermcolor"
)

// TODO: fix hex color support

// colorschemeColors is a map of color string names to Attribute types
type colorschemeColors map[string]interface{}

// ISprintf is a sprintf interface
type ISprintf func(...interface{}) string

// colorCache is a map of color string names to sprintf functions
type colorCache map[string]ISprintf

// Colorscheme is the struct for colorscheme
type Colorscheme struct {
	colors colorschemeColors
	cache  colorCache
}

var fgcolorschemeColorsMap = map[string]fcolor.Attribute{
	"black":   fcolor.FgBlack,
	"blue":    fcolor.FgBlue,
	"cyan":    fcolor.FgCyan,
	"green":   fcolor.FgGreen,
	"magenta": fcolor.FgMagenta,
	"red":     fcolor.FgRed,
	"white":   fcolor.FgWhite,
	"yellow":  fcolor.FgYellow,
}

var bgcolorschemeColorsMap = map[string]fcolor.Attribute{
	"black":   fcolor.BgBlack,
	"blue":    fcolor.BgBlue,
	"cyan":    fcolor.BgCyan,
	"green":   fcolor.BgGreen,
	"magenta": fcolor.BgMagenta,
	"red":     fcolor.BgRed,
	"white":   fcolor.BgWhite,
	"yellow":  fcolor.BgYellow,
}

var gocuiColorschemeColorsMap = map[string]gocui.Attribute{
	"black":   gocui.ColorBlack,
	"blue":    gocui.ColorBlue,
	"cyan":    gocui.ColorCyan,
	"green":   gocui.ColorGreen,
	"magenta": gocui.ColorMagenta,
	"red":     gocui.ColorRed,
	"white":   gocui.ColorWhite,
	"yellow":  gocui.ColorYellow,
}

// NewColorscheme ...
func NewColorscheme(colors colorschemeColors) *Colorscheme {
	return &Colorscheme{
		colors: colors,
		cache:  make(colorCache),
	}
}

// BaseFg ...
func (c *Colorscheme) BaseFg() gocui.Attribute {
	return c.gocuiFgColor("base")
}

// BaseBg ...
func (c *Colorscheme) BaseBg() gocui.Attribute {
	return c.gocuiBgColor("base")
}

// Chart ...
func (c *Colorscheme) Chart(a ...interface{}) string {
	return c.color("chart", a...)
}

// Marketbar ...
func (c *Colorscheme) Marketbar(a ...interface{}) string {
	return c.color("marketbar", a...)
}

// MarketbarSprintf ...
func (c *Colorscheme) MarketbarSprintf() ISprintf {
	return c.toSprintf("marketbar")
}

// MarketbarChangeSprintf ...
func (c *Colorscheme) MarketbarChangeSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change")
}

// MarketbarChangeDownSprintf ...
func (c *Colorscheme) MarketbarChangeDownSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change_down")
}

// MarketbarChangeUpSprintf ...
func (c *Colorscheme) MarketbarChangeUpSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change_up")
}

// MarketBarLabelActive ...
func (c *Colorscheme) MarketBarLabelActive(a ...interface{}) string {
	return c.color("marketbar_label_active", a...)
}

// Menu ...
func (c *Colorscheme) Menu(a ...interface{}) string {
	return c.color("menu", a...)
}

// MenuHeader ...
func (c *Colorscheme) MenuHeader(a ...interface{}) string {
	return c.color("menu_header", a...)
}

// MenuLabel ...
func (c *Colorscheme) MenuLabel(a ...interface{}) string {
	return c.color("menu_label", a...)
}

// MenuLabelActive ...
func (c *Colorscheme) MenuLabelActive(a ...interface{}) string {
	return c.color("menu_label_active", a...)
}

// Searchbar ...
func (c *Colorscheme) Searchbar(a ...interface{}) string {
	return c.color("searchbar", a...)
}

// Statusbar ...
func (c *Colorscheme) Statusbar(a ...interface{}) string {
	return c.color("statusbar", a...)
}

// TableColumnPrice ...
func (c *Colorscheme) TableColumnPrice(a ...interface{}) string {
	return c.color("table_column_price", a...)
}

// TableColumnPriceSprintf ...
func (c *Colorscheme) TableColumnPriceSprintf() ISprintf {
	return c.toSprintf("table_column_price")
}

// TableColumnChange ...
func (c *Colorscheme) TableColumnChange(a ...interface{}) string {
	return c.color("table_column_change", a...)
}

// TableColumnChangeSprintf ...
func (c *Colorscheme) TableColumnChangeSprintf() ISprintf {
	return c.toSprintf("table_column_change")
}

// TableColumnChangeDown ...
func (c *Colorscheme) TableColumnChangeDown(a ...interface{}) string {
	return c.color("table_column_change_down", a...)
}

// TableColumnChangeDownSprintf ...
func (c *Colorscheme) TableColumnChangeDownSprintf() ISprintf {
	return c.toSprintf("table_column_change_down")
}

// TableColumnChangeUp ...
func (c *Colorscheme) TableColumnChangeUp(a ...interface{}) string {
	return c.color("table_column_change_up", a...)
}

// TableColumnChangeUpSprintf ...
func (c *Colorscheme) TableColumnChangeUpSprintf() ISprintf {
	return c.toSprintf("table_column_change_up")
}

// TableHeader ...
func (c *Colorscheme) TableHeader(a ...interface{}) string {
	return c.color("table_header", a...)
}

// TableHeaderSprintf ...
func (c *Colorscheme) TableHeaderSprintf() ISprintf {
	return c.toSprintf("table_header")
}

// TableHeaderColumnActive ...
func (c *Colorscheme) TableHeaderColumnActive(a ...interface{}) string {
	return c.color("table_header_column_active", a...)
}

// TableHeaderColumnActiveSprintf ...
func (c *Colorscheme) TableHeaderColumnActiveSprintf() ISprintf {
	return c.toSprintf("table_header_column_active")
}

// TableRow ...
func (c *Colorscheme) TableRow(a ...interface{}) string {
	return c.color("table_row", a...)
}

// TableRowSprintf ...
func (c *Colorscheme) TableRowSprintf() ISprintf {
	return c.toSprintf("table_row")
}

// TableRowActive ...
func (c *Colorscheme) TableRowActive(a ...interface{}) string {
	return c.color("table_row_active", a...)
}

// TableRowFavorite ...
func (c *Colorscheme) TableRowFavorite(a ...interface{}) string {
	return c.color("table_row_favorite", a...)
}

// TableRowFavoriteSprintf ...
func (c *Colorscheme) TableRowFavoriteSprintf() ISprintf {
	return c.toSprintf("table_row_favorite")
}

// SetViewColor ...
func (c *Colorscheme) SetViewColor(view *gocui.View, name string) {
	view.FgColor = c.gocuiFgColor(name)
	view.BgColor = c.gocuiBgColor(name)
}

// SetViewActiveColor ...
func (c *Colorscheme) SetViewActiveColor(view *gocui.View, name string) {
	view.SelFgColor = c.gocuiFgColor(name)
	view.SelBgColor = c.gocuiBgColor(name)
}

func (c *Colorscheme) toSprintf(name string) ISprintf {
	if cached, ok := c.cache[name]; ok {
		return cached
	}

	var attrs []fcolor.Attribute
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := c.toFgAttr(v); ok {
			attrs = append(attrs, fg)
		}
	}
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := c.toBgAttr(v); ok {
			attrs = append(attrs, bg)
		}
	}
	if v, ok := c.colors[name+"_bold"].(bool); ok {
		if bold, ok := c.toBoldAttr(v); ok {
			attrs = append(attrs, bold)
		}
	}
	if v, ok := c.colors[name+"_underline"].(bool); ok {
		if underline, ok := c.toUnderlineAttr(v); ok {
			attrs = append(attrs, underline)
		}
	}

	c.cache[name] = fcolor.New(attrs...).SprintFunc()
	return c.cache[name]
}

func (c *Colorscheme) color(name string, a ...interface{}) string {
	return c.toSprintf(name)(a...)
}

func (c *Colorscheme) gocuiFgColor(name string) gocui.Attribute {
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := c.toGocuiAttr(v); ok {
			return fg
		}
	}

	return gocui.ColorDefault
}

func (c *Colorscheme) gocuiBgColor(name string) gocui.Attribute {
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := c.toGocuiAttr(v); ok {
			return bg
		}
	}

	return gocui.ColorDefault
}

func (c *Colorscheme) toFgAttr(v string) (fcolor.Attribute, bool) {
	if attr, ok := fgcolorschemeColorsMap[v]; ok {
		return attr, true
	}

	if code, ok := HexToAnsi(v); ok {
		return fcolor.Attribute(code), true
	}

	return 0, false
}

func (c *Colorscheme) toBgAttr(v string) (fcolor.Attribute, bool) {
	if attr, ok := bgcolorschemeColorsMap[v]; ok {
		return attr, true
	}

	if code, ok := HexToAnsi(v); ok {
		return fcolor.Attribute(code), true
	}

	return 0, false
}

// toBoldAttr converts a boolean to an Attribute type
func (c *Colorscheme) toBoldAttr(v bool) (fcolor.Attribute, bool) {
	return fcolor.Bold, v
}

// toUnderlineAttr converts a boolean to an Attribute type
func (c *Colorscheme) toUnderlineAttr(v bool) (fcolor.Attribute, bool) {
	return fcolor.Underline, v
}

// toGocuiAttr converts a color string name to a gocui Attribute type
func (c *Colorscheme) toGocuiAttr(v string) (gocui.Attribute, bool) {
	if attr, ok := gocuiColorschemeColorsMap[v]; ok {
		return attr, true
	}

	if code, ok := HexToAnsi(v); ok {
		return gocui.Attribute(code), true
	}

	return 0, false
}

// HexToAnsi converts a hex color string to a uint8 ansi code
func HexToAnsi(h string) (uint8, bool) {
	if h == "" {
		return 0, false
	}

	n, err := strconv.Atoi(h)
	if err == nil {
		if n <= 255 {
			return uint8(n), true
		}
	}

	code, err := xtermcolor.FromHexStr(h)
	if err != nil {
		return 0, false
	}

	return code, true
}

// gocui can use xterm colors
