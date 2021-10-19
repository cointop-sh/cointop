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
	"errors"

	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen
var outMode OutputMode

// Init initializes the screen for use.
func Init() error {
	outMode = OutputNormal
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		screen = s
		return nil
	}
}

// Close cleans up the terminal, restoring terminal modes, etc.
func Close() {
	screen.Fini()
}

// Flush updates the screen.
func Flush() error {
	screen.Show()
	return nil
}

// SetCursor displays the terminal cursor at the given location.
func SetCursor(x, y int) {
	screen.ShowCursor(x, y)
}

// HideCursor hides the terminal cursor.
func HideCursor() {
	SetCursor(-1, -1)
}

// Size returns the screen size as width, height in character cells.
func Size() (int, int) {
	return screen.Size()
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

// InputMode is not used.
type InputMode int

// Unused input modes; here for compatibility.
const (
	InputCurrent InputMode = iota
	InputEsc
	InputAlt
	InputMouse
)

// SetInputMode enables mouse if requested
func SetInputMode(mode InputMode) InputMode {
	if mode&InputMouse != 0 {
		screen.EnableMouse()
	}
	// We don't do anything else right now
	return InputEsc
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

// Sync forces a resync of the screen.
func Sync() error {
	screen.Sync()
	return nil
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

// EventType represents the type of event.
type EventType uint8

// Modifier represents the possible modifier keys.
// type Modifier tcell.ModMask

// Key is a key press.
// type Key tcell.Key

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type   EventType
	Mod    tcell.ModMask
	Key    tcell.Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

// Event types.
const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
)

// Keys codes.
const (
	MouseLeft      = tcell.KeyF63 // arbitrary assignments
	MouseRight     = tcell.KeyF62
	MouseMiddle    = tcell.KeyF61
	MouseRelease   = tcell.KeyF60
	MouseWheelUp   = tcell.KeyF59
	MouseWheelDown = tcell.KeyF58
)

// Modifiers.
const (
	ModAlt = tcell.ModAlt
)

func makeEvent(tev tcell.Event) Event {
	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		return Event{Type: EventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return Event{Type: EventResize, Width: w, Height: h}
	case *tcell.EventKey:
		// k := tev.Key()
		// ch := rune(0)
		// Remove space hack
		// if k == tcell.KeyRune {
		// 	ch = tev.Rune()
		// 	if ch == ' ' {
		// 		k = tcell.Key(' ')
		// 	} else {
		// 		k = tcell.Key(0)
		// 	}
		// }
		// mod := tev.Modifiers()
		return Event{
			Type: EventKey,
			Key:  tev.Key(),
			Ch:   tev.Rune(),
			Mod:  tev.Modifiers(),
		}
	case *tcell.EventMouse:
		x, y := tev.Position()
		button := tev.Buttons()
		// Don't worry about combo buttons for now
		key := tcell.KeyNUL
		if button&tcell.Button1 != 0 {
			key = MouseLeft
		} else if button&tcell.Button2 != 0 {
			key = MouseRight
		} else if button&tcell.Button3 != 0 {
			key = MouseMiddle
		} else if button&tcell.WheelUp != 0 {
			key = MouseWheelUp
		} else if button&tcell.WheelDown != 0 {
			key = MouseWheelDown
		} else {
			return Event{Type: EventNone}
		}
		return Event{
			Type:   EventMouse,
			MouseX: x,
			MouseY: y,
			Key:    key,
			Mod:    0, // not currently supported
		}
	default:
		return Event{Type: EventNone}
	}
}

// ParseEvent is not supported.
func ParseEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollRawEvent is not supported.
func PollRawEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollEvent blocks until an event is ready, and then returns it.
func PollEvent() Event {
	ev := screen.PollEvent()
	return makeEvent(ev)
}

// Interrupt posts an interrupt event.
func Interrupt() {
	screen.PostEvent(tcell.NewEventInterrupt(nil))
}

// Cell represents a single character cell on screen.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}
