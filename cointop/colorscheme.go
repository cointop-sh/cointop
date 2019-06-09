package cointop

import (
	"github.com/fatih/color"
)

// Colors ..
type Colors map[string]interface{}

// Cache ..
type Cache map[string]func(...interface{}) string

// ColorScheme ...
type ColorScheme struct {
	colors Colors
	cache  Cache
}

// NewColorScheme ...
func NewColorScheme(colors Colors) *ColorScheme {
	return &ColorScheme{
		colors: colors,
		cache:  make(Cache),
	}
}

// RowText ...
func (c *ColorScheme) RowText(a ...interface{}) string {
	name := "row_text"
	return c.color(name, a...)
}

func (c *ColorScheme) color(name string, a ...interface{}) string {
	if cached, ok := c.cache[name]; ok {
		return cached(a...)
	}

	var colors []color.Attribute
	if v, ok := c.colors[name+"_fg"].(string); ok {
		if fg, ok := toFgAttr(v); ok {
			colors = append(colors, fg)
		}
	}
	if v, ok := c.colors[name+"_bg"].(string); ok {
		if bg, ok := toBgAttr(v); ok {
			colors = append(colors, bg)
		}
	}
	if v, ok := c.colors[name+"_bold"].(bool); ok {
		if bold, ok := toBoldAttr(v); ok {
			colors = append(colors, bold)
		}
	}
	if v, ok := c.colors[name+"_underline"].(bool); ok {
		if underline, ok := toUnderlineAttr(v); ok {
			colors = append(colors, underline)
		}
	}

	c.cache[name] = color.New(colors...).SprintFunc()
	return c.cache[name](a...)
}

var fgColorsMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"blue":    color.FgBlue,
	"cyan":    color.FgCyan,
	"green":   color.FgGreen,
	"magenta": color.FgMagenta,
	"red":     color.FgRed,
	"white":   color.FgWhite,
	"yellow":  color.FgYellow,
}

var bgColorsMap = map[string]color.Attribute{
	"black":   color.BgBlack,
	"blue":    color.BgBlue,
	"cyan":    color.BgCyan,
	"green":   color.BgGreen,
	"magenta": color.BgMagenta,
	"red":     color.BgRed,
	"white":   color.BgWhite,
	"yellow":  color.BgYellow,
}

func toFgAttr(c string) (color.Attribute, bool) {
	attr, ok := fgColorsMap[c]
	return attr, ok
}

func toBgAttr(c string) (color.Attribute, bool) {
	attr, ok := bgColorsMap[c]
	return attr, ok
}

func toBoldAttr(v bool) (color.Attribute, bool) {
	return color.Bold, v
}

func toUnderlineAttr(v bool) (color.Attribute, bool) {
	return color.Underline, v
}

// CointopColorscheme ...
var CointopColorscheme = `
row_text_fg = "white"
row_text_bg = ""
row_text_bold = false
`
