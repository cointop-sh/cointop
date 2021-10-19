// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package termbox is a compatibility layer to allow tcell to emulate
// the github.com/nsf/termbox package.
package termbox

import (
	"github.com/gdamore/tcell/v2"
)

// Ugly globals
var screen tcell.Screen
var outMode OutputMode

func SetScreen(s tcell.Screen) {
	screen = s
}

// Attribute affects the presentation of characters, such as color, boldness,
// and so forth.
type Attribute uint16

// Colors first.  The order here is significant.
const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Other attributes.
const (
	AttrBold Attribute = 1 << (9 + iota)
	AttrUnderline
	AttrReverse
)

func fixColor(c tcell.Color) tcell.Color {
	if c == tcell.ColorDefault {
		return c
	}
	switch outMode {
	case OutputNormal:
		c = tcell.PaletteColor(int(c) & 0xf)
	case Output256:
		c = tcell.PaletteColor(int(c) & 0xff)
	case Output216:
		c = tcell.PaletteColor(int(c)%216 + 16)
	case OutputGrayscale:
		c %= tcell.PaletteColor(int(c)%24 + 232)
	default:
		c = tcell.ColorDefault
	}
	return c
}

func mkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault

	f := tcell.PaletteColor(int(fg)&0x1ff - 1)
	b := tcell.PaletteColor(int(bg)&0x1ff - 1)

	f = fixColor(f)
	b = fixColor(b)
	st = st.Foreground(f).Background(b)
	if (fg|bg)&AttrBold != 0 {
		st = st.Bold(true)
	}
	if (fg|bg)&AttrUnderline != 0 {
		st = st.Underline(true)
	}
	if (fg|bg)&AttrReverse != 0 {
		st = st.Reverse(true)
	}
	return st
}

// Clear clears the screen with the given attributes.
func Clear(fg, bg Attribute) {
	st := mkStyle(fg, bg)
	w, h := screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

// OutputMode represents an output mode, which determines how colors
// are used.  See the termbox documentation for an explanation.
type OutputMode int

// OutputMode values.
const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
)

// SetOutputMode is used to set the color palette used.
func SetOutputMode(mode OutputMode) OutputMode {
	if screen.Colors() < 256 {
		mode = OutputNormal
	}
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale:
		outMode = mode
		return mode
	default:
		return outMode
	}
}

// scaledColor returns a Color that is proportional to the x/y coordinates
func scaledColor(x, y int) tcell.Color {
	w, h := screen.Size()
	blu := int32(255 * float64(x) / float64(w))
	grn := int32(255 * float64(y) / float64(h))
	red := int32(200)
	return tcell.NewRGBColor(red, grn, blu)
}

// SetCell sets the character cell at a given location to the given
// content (rune) and attributes.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	st := mkStyle(fg, bg)
	// Set the foreground color to a scaled version of the coordinates
	st = st.Foreground(scaledColor(x, y))
	screen.SetContent(x, y, ch, nil, st)
}

// PollEvent blocks until an event is ready, and then returns it.
func PollEvent() tcell.Event {
	return screen.PollEvent()
}
