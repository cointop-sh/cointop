package cointop

import (
	"fmt"
	"strconv"
	"sync"

	fcolor "github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/tomnomnom/xtermcolor"
)

// TODO: fix hex color support

// ColorschemeColors is a map of color string names to Attribute types
type ColorschemeColors map[string]interface{}

// ISprintf is a sprintf interface
type ISprintf func(...interface{}) string

// ColorCache is a map of color string names to sprintf functions
type ColorCache map[string]ISprintf

// Colorscheme is the struct for colorscheme
type Colorscheme struct {
	colors     ColorschemeColors
	cache      ColorCache
	cacheMutex sync.RWMutex
}

var FgColorschemeColorsMap = map[string]fcolor.Attribute{
	"black":   fcolor.FgBlack,
	"blue":    fcolor.FgBlue,
	"cyan":    fcolor.FgCyan,
	"green":   fcolor.FgGreen,
	"magenta": fcolor.FgMagenta,
	"red":     fcolor.FgRed,
	"white":   fcolor.FgWhite,
	"yellow":  fcolor.FgYellow,
}

var BgColorschemeColorsMap = map[string]fcolor.Attribute{
	"black":   fcolor.BgBlack,
	"blue":    fcolor.BgBlue,
	"cyan":    fcolor.BgCyan,
	"green":   fcolor.BgGreen,
	"magenta": fcolor.BgMagenta,
	"red":     fcolor.BgRed,
	"white":   fcolor.BgWhite,
	"yellow":  fcolor.BgYellow,
}

// See more: vendor/github.com/mattn/go-colorable/colorable_windows.go:905
// any new color for the below mapping should be compatible with this above list
var TcellColorschemeColorsMap = map[string]tcell.Color{
	"black":   tcell.ColorBlack,
	"blue":    tcell.ColorNavy,
	"cyan":    tcell.ColorTeal,
	"green":   tcell.ColorGreen,
	"magenta": tcell.ColorPurple,
	"red":     tcell.ColorMaroon,
	"white":   tcell.ColorSilver,
	"yellow":  tcell.ColorOlive,
}

// NewColorscheme ...
func NewColorscheme(colors ColorschemeColors) *Colorscheme {
	return &Colorscheme{
		colors:     colors,
		cache:      make(ColorCache),
		cacheMutex: sync.RWMutex{},
	}
}

// BaseFg ...
func (c *Colorscheme) BaseFg() tcell.Color {
	return c.FgColor("base")
}

// BaseBg ...
func (c *Colorscheme) BaseBg() tcell.Color {
	return c.BgColor("base")
}

// Chart ...
func (c *Colorscheme) Chart(a ...interface{}) string {
	return c.Color("chart", a...)
}

// Marketbar ...
func (c *Colorscheme) Marketbar(a ...interface{}) string {
	return c.Color("marketbar", a...)
}

// MarketbarSprintf ...
func (c *Colorscheme) MarketbarSprintf() ISprintf {
	return c.ToSprintf("marketbar")
}

// MarketbarChangeSprintf ...
func (c *Colorscheme) MarketbarChangeSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.ToSprintf("table_column_change")
}

// MarketbarChangeDownSprintf ...
func (c *Colorscheme) MarketbarChangeDownSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.ToSprintf("table_column_change_down")
}

// MarketbarChangeUpSprintf ...
func (c *Colorscheme) MarketbarChangeUpSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.ToSprintf("table_column_change_up")
}

// MarketBarLabelActive ...
func (c *Colorscheme) MarketBarLabelActive(a ...interface{}) string {
	return c.Color("marketbar_label_active", a...)
}

// Menu ...
func (c *Colorscheme) Menu(a ...interface{}) string {
	return c.Color("menu", a...)
}

// MenuHeader ...
func (c *Colorscheme) MenuHeader(a ...interface{}) string {
	return c.Color("menu_header", a...)
}

// MenuLabel ...
func (c *Colorscheme) MenuLabel(a ...interface{}) string {
	return c.Color("menu_label", a...)
}

// MenuLabelActive ...
func (c *Colorscheme) MenuLabelActive(a ...interface{}) string {
	return c.Color("menu_label_active", a...)
}

// Searchbar ...
func (c *Colorscheme) Searchbar(a ...interface{}) string {
	return c.Color("searchbar", a...)
}

// Statusbar ...
func (c *Colorscheme) Statusbar(a ...interface{}) string {
	return c.Color("statusbar", a...)
}

// TableColumnPrice ...
func (c *Colorscheme) TableColumnPrice(a ...interface{}) string {
	return c.Color("table_column_price", a...)
}

// TableColumnPriceSprintf ...
func (c *Colorscheme) TableColumnPriceSprintf() ISprintf {
	return c.ToSprintf("table_column_price")
}

// TableColumnChange ...
func (c *Colorscheme) TableColumnChange(a ...interface{}) string {
	return c.Color("table_column_change", a...)
}

// TableColumnChangeSprintf ...
func (c *Colorscheme) TableColumnChangeSprintf() ISprintf {
	return c.ToSprintf("table_column_change")
}

// TableColumnChangeDown ...
func (c *Colorscheme) TableColumnChangeDown(a ...interface{}) string {
	return c.Color("table_column_change_down", a...)
}

// TableColumnChangeDownSprintf ...
func (c *Colorscheme) TableColumnChangeDownSprintf() ISprintf {
	return c.ToSprintf("table_column_change_down")
}

// TableColumnChangeUp ...
func (c *Colorscheme) TableColumnChangeUp(a ...interface{}) string {
	return c.Color("table_column_change_up", a...)
}

// TableColumnChangeUpSprintf ...
func (c *Colorscheme) TableColumnChangeUpSprintf() ISprintf {
	return c.ToSprintf("table_column_change_up")
}

// TableHeader ...
func (c *Colorscheme) TableHeader(a ...interface{}) string {
	return c.Color("table_header", a...)
}

// TableHeaderSprintf ...
func (c *Colorscheme) TableHeaderSprintf() ISprintf {
	return c.ToSprintf("table_header")
}

// TableHeaderColumnActive ...
func (c *Colorscheme) TableHeaderColumnActive(a ...interface{}) string {
	return c.Color("table_header_column_active", a...)
}

// TableHeaderColumnActiveSprintf ...
func (c *Colorscheme) TableHeaderColumnActiveSprintf() ISprintf {
	return c.ToSprintf("table_header_column_active")
}

// TableRow ...
func (c *Colorscheme) TableRow(a ...interface{}) string {
	return c.Color("table_row", a...)
}

// TableRowSprintf ...
func (c *Colorscheme) TableRowSprintf() ISprintf {
	return c.ToSprintf("table_row")
}

// TableRowActive ...
func (c *Colorscheme) TableRowActive(a ...interface{}) string {
	return c.Color("table_row_active", a...)
}

// TableRowFavorite ...
func (c *Colorscheme) TableRowFavorite(a ...interface{}) string {
	return c.Color("table_row_favorite", a...)
}

// TableRowFavoriteSprintf ...
func (c *Colorscheme) TableRowFavoriteSprintf() ISprintf {
	return c.ToSprintf("table_row_favorite")
}

// Default ...
func (c *Colorscheme) Default(a ...interface{}) string {
	return fmt.Sprintf(a[0].(string), a[1:]...)
}

func (c *Colorscheme) ToSprintf(name string) ISprintf {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	if cached, ok := c.cache[name]; ok {
		return cached
	}

	var attrs []fcolor.Attribute
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := c.ToFgAttr(v); ok {
			attrs = append(attrs, fg)
		}
	}
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := c.ToBgAttr(v); ok {
			attrs = append(attrs, bg)
		}
	}
	if v, ok := c.colors[name+"_bold"].(bool); ok {
		if bold, ok := c.ToBoldAttr(v); ok {
			attrs = append(attrs, bold)
		}
	}
	if v, ok := c.colors[name+"_underline"].(bool); ok {
		if underline, ok := c.ToUnderlineAttr(v); ok {
			attrs = append(attrs, underline)
		}
	}

	c.cache[name] = fcolor.New(attrs...).SprintFunc()
	return c.cache[name]
}

func (c *Colorscheme) Color(name string, a ...interface{}) string {
	return c.ToSprintf(name)(a...)
}

func (c *Colorscheme) FgColor(name string) tcell.Color {
	fg := c.tcellColor(name + "_fg")

	// TODO: fixme
	// if v, ok := c.colors[name+"_bold"].(bool); ok {
	// 	if v {
	// 		attrs = append(attrs, gocui.AttrBold)
	// 	}
	// }
	// if v, ok := c.colors[name+"_underline"].(bool); ok {
	// 	if v {
	// 		attrs = append(attrs, gocui.AttrUnderline)
	// 	}
	// }

	return fg
}

func (c *Colorscheme) BgColor(name string) tcell.Color {
	return c.tcellColor(name + "_bg")
}

// tcellColor can supply for types of color name: specific mapped name, tcell color name, hex
// Examples: black, honeydew, #000000
func (c *Colorscheme) tcellColor(name string) tcell.Color {
	v, ok := c.colors[name].(string)
	if !ok {
		return tcell.ColorDefault
	}

	if color, found := TcellColorschemeColorsMap[v]; found {
		return color
	}

	color := tcell.GetColor(name)
	if color != tcell.ColorDefault {
		return color
	}

	code, _ := HexToAnsi(v)
	return tcell.PaletteColor(int(code) & 0xff)
}

func (c *Colorscheme) ToFgAttr(v string) (fcolor.Attribute, bool) {
	if attr, ok := FgColorschemeColorsMap[v]; ok {
		return attr, true
	}

	if code, ok := HexToAnsi(v); ok {
		return fcolor.Attribute(code), true
	}

	return 0, false
}

func (c *Colorscheme) ToBgAttr(v string) (fcolor.Attribute, bool) {
	if attr, ok := BgColorschemeColorsMap[v]; ok {
		return attr, true
	}

	if code, ok := HexToAnsi(v); ok {
		return fcolor.Attribute(code), true
	}

	return 0, false
}

// ToBoldAttr converts a boolean to an Attribute type
func (c *Colorscheme) ToBoldAttr(v bool) (fcolor.Attribute, bool) {
	return fcolor.Bold, v
}

// ToUnderlineAttr converts a boolean to an Attribute type
func (c *Colorscheme) ToUnderlineAttr(v bool) (fcolor.Attribute, bool) {
	return fcolor.Underline, v
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
