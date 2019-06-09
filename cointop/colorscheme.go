package cointop

import (
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

// colorschemeColors ..
type colorschemeColors map[string]interface{}

// ISprintf ...
type ISprintf func(...interface{}) string

// colorCache ..
type colorCache map[string]ISprintf

// Colorscheme ...
type Colorscheme struct {
	colors colorschemeColors
	cache  colorCache
}

var fgcolorschemeColorsMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"blue":    color.FgBlue,
	"cyan":    color.FgCyan,
	"green":   color.FgGreen,
	"magenta": color.FgMagenta,
	"red":     color.FgRed,
	"white":   color.FgWhite,
	"yellow":  color.FgYellow,
}

var bgcolorschemeColorsMap = map[string]color.Attribute{
	"black":   color.BgBlack,
	"blue":    color.BgBlue,
	"cyan":    color.BgCyan,
	"green":   color.BgGreen,
	"magenta": color.BgMagenta,
	"red":     color.BgRed,
	"white":   color.BgWhite,
	"yellow":  color.BgYellow,
}

var gocuicolorschemeColorsMap = map[string]gocui.Attribute{
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

	var colors []color.Attribute
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := c.toFgAttr(v); ok {
			colors = append(colors, fg)
		}
	}
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := c.toBgAttr(v); ok {
			colors = append(colors, bg)
		}
	}
	if v, ok := c.colors[name+"_bold"].(bool); ok {
		if bold, ok := c.toBoldAttr(v); ok {
			colors = append(colors, bold)
		}
	}
	if v, ok := c.colors[name+"_underline"].(bool); ok {
		if underline, ok := c.toUnderlineAttr(v); ok {
			colors = append(colors, underline)
		}
	}

	c.cache[name] = color.New(colors...).SprintFunc()
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

func (c *Colorscheme) toFgAttr(k string) (color.Attribute, bool) {
	attr, ok := fgcolorschemeColorsMap[k]
	return attr, ok
}

func (c *Colorscheme) toBgAttr(k string) (color.Attribute, bool) {
	attr, ok := bgcolorschemeColorsMap[k]
	return attr, ok
}

func (c *Colorscheme) toBoldAttr(v bool) (color.Attribute, bool) {
	return color.Bold, v
}

func (c *Colorscheme) toUnderlineAttr(v bool) (color.Attribute, bool) {
	return color.Underline, v
}

func (c *Colorscheme) toGocuiAttr(k string) (gocui.Attribute, bool) {
	attr, ok := gocuicolorschemeColorsMap[k]
	return attr, ok
}
