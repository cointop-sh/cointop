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

// TODO: remove compatability layer
func FixColor(c tcell.Color) tcell.Color {
	if c == tcell.ColorDefault {
		return c
	}
	c = tcell.PaletteColor(int(c) & 0xff)
	// switch g.outputMode {
	// case OutputNormal:
	// 	c = tcell.PaletteColor(int(c) & 0xf)
	// case Output256:
	// 	c = tcell.PaletteColor(int(c) & 0xff)
	// case Output216:
	// 	c = tcell.PaletteColor(int(c)%216 + 16)
	// case OutputGrayscale:
	// 	c %= tcell.PaletteColor(int(c)%24 + 232)
	// default:
	// 	c = tcell.ColorDefault
	// }
	return c
}

// TODO: remove compatability layer
func MkColor(color Attribute) tcell.Color {
	if color == ColorDefault {
		return tcell.ColorDefault
	} else {
		return FixColor(tcell.PaletteColor(int(color)&0x1ff - 1))
	}
}

func MkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault
	if fg != ColorDefault {
		// f := tcell.PaletteColor(int(fg)&0x1ff - 1)
		// f = g.fixColor(f)
		st = st.Foreground(MkColor(fg))
	}
	if bg != ColorDefault {
		// b := tcell.PaletteColor(int(bg)&0x1ff - 1)
		// b = g.fixColor(b)
		st = st.Background(MkColor(bg))
	}
	// TODO: fixme
	// if (fg|bg)&AttrBold != 0 {
	// 	st = st.Bold(true)
	// }
	// if (fg|bg)&AttrUnderline != 0 {
	// 	st = st.Underline(true)
	// }
	// if (fg|bg)&AttrReverse != 0 {
	// 	st = st.Reverse(true)
	// }
	return st
}
