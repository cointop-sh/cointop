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

// ColorScheme ...
type ColorScheme struct {
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

// NewColorScheme ...
func NewColorScheme(colors colorschemeColors) *ColorScheme {
	return &ColorScheme{
		colors: colors,
		cache:  make(colorCache),
	}
}

// BaseFg ...
func (c *ColorScheme) BaseFg() gocui.Attribute {
	return c.gocuiFgColor("base")
}

// BaseBg ...
func (c *ColorScheme) BaseBg() gocui.Attribute {
	return c.gocuiBgColor("base")
}

// Chart ...
func (c *ColorScheme) Chart(a ...interface{}) string {
	return c.color("chart", a...)
}

// Marketbar ...
func (c *ColorScheme) Marketbar(a ...interface{}) string {
	return c.color("marketbar", a...)
}

// MarketbarSprintf ...
func (c *ColorScheme) MarketbarSprintf() ISprintf {
	return c.toSprintf("marketbar")
}

// MarketbarChangeSprintf ...
func (c *ColorScheme) MarketbarChangeSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change")
}

// MarketbarChangeDownSprintf ...
func (c *ColorScheme) MarketbarChangeDownSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change_down")
}

// MarketbarChangeUpSprintf ...
func (c *ColorScheme) MarketbarChangeUpSprintf() ISprintf {
	// NOTE: reusing table styles
	return c.toSprintf("table_column_change_up")
}

// MarketBarLabelActive ...
func (c *ColorScheme) MarketBarLabelActive(a ...interface{}) string {
	return c.color("marketbar_label_active", a...)
}

// Menu ...
func (c *ColorScheme) Menu(a ...interface{}) string {
	return c.color("menu", a...)
}

// MenuHeader ...
func (c *ColorScheme) MenuHeader(a ...interface{}) string {
	return c.color("menu_header", a...)
}

// MenuLabel ...
func (c *ColorScheme) MenuLabel(a ...interface{}) string {
	return c.color("menu_label", a...)
}

// MenuLabelActive ...
func (c *ColorScheme) MenuLabelActive(a ...interface{}) string {
	return c.color("menu_label_active", a...)
}

// Searchbar ...
func (c *ColorScheme) Searchbar(a ...interface{}) string {
	return c.color("searchbar", a...)
}

// Statusbar ...
func (c *ColorScheme) Statusbar(a ...interface{}) string {
	return c.color("statusbar", a...)
}

// TableColumnPrice ...
func (c *ColorScheme) TableColumnPrice(a ...interface{}) string {
	return c.color("table_column_price", a...)
}

// TableColumnPriceSprintf ...
func (c *ColorScheme) TableColumnPriceSprintf() ISprintf {
	return c.toSprintf("table_column_price")
}

// TableColumnChange ...
func (c *ColorScheme) TableColumnChange(a ...interface{}) string {
	return c.color("table_column_change", a...)
}

// TableColumnChangeSprintf ...
func (c *ColorScheme) TableColumnChangeSprintf() ISprintf {
	return c.toSprintf("table_column_change")
}

// TableColumnChangeDown ...
func (c *ColorScheme) TableColumnChangeDown(a ...interface{}) string {
	return c.color("table_column_change_down", a...)
}

// TableColumnChangeDownSprintf ...
func (c *ColorScheme) TableColumnChangeDownSprintf() ISprintf {
	return c.toSprintf("table_column_change_down")
}

// TableColumnChangeUp ...
func (c *ColorScheme) TableColumnChangeUp(a ...interface{}) string {
	return c.color("table_column_change_up", a...)
}

// TableColumnChangeUpSprintf ...
func (c *ColorScheme) TableColumnChangeUpSprintf() ISprintf {
	return c.toSprintf("table_column_change_up")
}

// TableHeader ...
func (c *ColorScheme) TableHeader(a ...interface{}) string {
	return c.color("table_header", a...)
}

// TableHeaderSprintf ...
func (c *ColorScheme) TableHeaderSprintf() ISprintf {
	return c.toSprintf("table_header")
}

// TableHeaderColumnActive ...
func (c *ColorScheme) TableHeaderColumnActive(a ...interface{}) string {
	return c.color("table_header_column_active", a...)
}

// TableHeaderColumnActiveSprintf ...
func (c *ColorScheme) TableHeaderColumnActiveSprintf() ISprintf {
	return c.toSprintf("table_header_column_active")
}

// TableRow ...
func (c *ColorScheme) TableRow(a ...interface{}) string {
	return c.color("table_row", a...)
}

// TableRowSprintf ...
func (c *ColorScheme) TableRowSprintf() ISprintf {
	return c.toSprintf("table_row")
}

// TableRowActive ...
func (c *ColorScheme) TableRowActive(a ...interface{}) string {
	return c.color("table_row_active", a...)
}

// TableRowFavorite ...
func (c *ColorScheme) TableRowFavorite(a ...interface{}) string {
	return c.color("table_row_favorite", a...)
}

// TableRowFavoriteSprintf ...
func (c *ColorScheme) TableRowFavoriteSprintf() ISprintf {
	return c.toSprintf("table_row_favorite")
}

// SetViewColor ...
func (c *ColorScheme) SetViewColor(view *gocui.View, name string) {
	view.FgColor = c.gocuiFgColor(name)
	view.BgColor = c.gocuiBgColor(name)
}

// SetViewActiveColor ...
func (c *ColorScheme) SetViewActiveColor(view *gocui.View, name string) {
	view.SelFgColor = c.gocuiFgColor(name)
	view.SelBgColor = c.gocuiBgColor(name)
}

func (c *ColorScheme) toSprintf(name string) ISprintf {
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

func (c *ColorScheme) color(name string, a ...interface{}) string {
	return c.toSprintf(name)(a...)
}

func (c *ColorScheme) gocuiFgColor(name string) gocui.Attribute {
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := c.toGocuiAttr(v); ok {
			return fg
		}
	}

	return gocui.ColorDefault
}

func (c *ColorScheme) gocuiBgColor(name string) gocui.Attribute {
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := c.toGocuiAttr(v); ok {
			return bg
		}
	}

	return gocui.ColorDefault
}

func (c *ColorScheme) toFgAttr(k string) (color.Attribute, bool) {
	attr, ok := fgcolorschemeColorsMap[k]
	return attr, ok
}

func (c *ColorScheme) toBgAttr(k string) (color.Attribute, bool) {
	attr, ok := bgcolorschemeColorsMap[k]
	return attr, ok
}

func (c *ColorScheme) toBoldAttr(v bool) (color.Attribute, bool) {
	return color.Bold, v
}

func (c *ColorScheme) toUnderlineAttr(v bool) (color.Attribute, bool) {
	return color.Underline, v
}

func (c *ColorScheme) toGocuiAttr(k string) (gocui.Attribute, bool) {
	attr, ok := gocuicolorschemeColorsMap[k]
	return attr, ok
}
