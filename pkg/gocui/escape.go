// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type escapeInterpreter struct {
	state    escapeState
	curch    rune
	csiParam []string
	curStyle tcell.Style
	// mode                   OutputMode
}

type escapeState int

const (
	stateNone escapeState = iota
	stateEscape
	stateCSI
	stateParams
)

var (
	errNotCSI        = errors.New("Not a CSI escape sequence")
	errCSIParseError = errors.New("CSI escape sequence parsing error")
	errCSITooLong    = errors.New("CSI escape sequence is too long")
)

// runes in case of error will output the non-parsed runes as a string.
func (ei *escapeInterpreter) runes() []rune {
	switch ei.state {
	case stateNone:
		return []rune{0x1b}
	case stateEscape:
		return []rune{0x1b, ei.curch}
	case stateCSI:
		return []rune{0x1b, '[', ei.curch}
	case stateParams:
		ret := []rune{0x1b, '['}
		for _, s := range ei.csiParam {
			ret = append(ret, []rune(s)...)
			ret = append(ret, ';')
		}
		return append(ret, ei.curch)
	}
	return nil
}

// newEscapeInterpreter returns an escapeInterpreter that will be able to parse
// terminal escape sequences.
func newEscapeInterpreter() *escapeInterpreter {
	ei := &escapeInterpreter{
		state:    stateNone,
		curStyle: tcell.StyleDefault,
		// mode:       mode,
	}
	return ei
}

// reset sets the escapeInterpreter in initial state.
func (ei *escapeInterpreter) reset() {
	ei.state = stateNone
	ei.curStyle = tcell.StyleDefault
	ei.csiParam = nil
}

// parseOne parses a rune. If isEscape is true, it means that the rune is part
// of an escape sequence, and as such should not be printed verbatim. Otherwise,
// it's not an escape sequence.
func (ei *escapeInterpreter) parseOne(ch rune) (isEscape bool, err error) {
	// Sanity checks
	if len(ei.csiParam) > 20 {
		return false, errCSITooLong
	}
	if len(ei.csiParam) > 0 && len(ei.csiParam[len(ei.csiParam)-1]) > 255 {
		return false, errCSITooLong
	}

	ei.curch = ch

	switch ei.state {
	case stateNone:
		if ch == 0x1b {
			ei.state = stateEscape
			return true, nil
		}
		return false, nil
	case stateEscape:
		if ch == '[' {
			ei.state = stateCSI
			return true, nil
		}
		return false, errNotCSI
	case stateCSI:
		switch {
		case ch >= '0' && ch <= '9':
			ei.csiParam = append(ei.csiParam, "")
		case ch == 'm':
			ei.csiParam = append(ei.csiParam, "0")
		default:
			return false, errCSIParseError
		}
		ei.state = stateParams
		fallthrough
	case stateParams:
		switch {
		case ch >= '0' && ch <= '9':
			ei.csiParam[len(ei.csiParam)-1] += string(ch)
			return true, nil
		case ch == ';':
			ei.csiParam = append(ei.csiParam, "")
			return true, nil
		case ch == 'm':
			var err error
			err = ei.parseEscapeParams()
			// switch ei.mode {
			// case OutputNormal:
			// 	err = ei.outputNormal()
			// case Output256:
			// 	err = ei.output256()
			// }
			if err != nil {
				return false, errCSIParseError
			}

			ei.state = stateNone
			ei.csiParam = nil
			return true, nil
		default:
			return false, errCSIParseError
		}
	}
	return false, nil
}

// parseEscapeParams interprets an escape sequence as a style modifier
// allows you to leverage the 256-colors terminal mode:
//   0x01 - 0x08: the 8 colors as in OutputNormal (black, red, green, yellow, blue, magenta, cyan, white)
//   0x09 - 0x10: Color* | AttrBold
//   0x11 - 0xe8: 216 different colors
//   0xe9 - 0x1ff: 24 different shades of grey
// see https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
// see https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
// 256-colors: ESC[ 38;5;${ID}m   # foreground
// 256-colors: ESC[ 48;5;${ID}m   # background
// 24-bit ESC[ 38;2;⟨r⟩;⟨g⟩;⟨b⟩ m Select RGB foreground color
// 24-bit ESC[ 48;2;⟨r⟩;⟨g⟩;⟨b⟩ m Select RGB background color
func (ei *escapeInterpreter) parseEscapeParams() error {
	// TODO: cache escape -> Style
	// convert params to int
	params := make([]int, len(ei.csiParam))
	for i, param := range ei.csiParam {
		if p, err := strconv.Atoi(param); err == nil {
			params[i] = p
		} else {
			return errCSIParseError
		}
	}

	// consume elements of params until done
	pos := 0
	for ok := true; ok; ok = pos < len(params) {
		p := params[pos]
		switch {
		case p >= 30 && p <= 37:
			ei.curStyle = ei.curStyle.Foreground(tcell.PaletteColor(p - 30))
		case p == 39:
			ei.curStyle = ei.curStyle.Foreground(tcell.ColorDefault)
		case p >= 40 && p <= 47:
			ei.curStyle = ei.curStyle.Background(tcell.PaletteColor(p - 40))
		case p == 49:
			ei.curStyle = ei.curStyle.Background(tcell.ColorDefault)
		case p == 1:
			ei.curStyle = ei.curStyle.Bold(true)
		case p == 4:
			ei.curStyle = ei.curStyle.Underline(true)
		case p == 7:
			ei.curStyle = ei.curStyle.Reverse(true)
		case p == 0:
			ei.curStyle = tcell.StyleDefault
		case p == 38 || p == 48: // 256-color or 24-bit
			// parse mode and additional params to generate a color
			mode := params[pos+1] // second param - 2 or 5
			var x tcell.Color
			if mode == 5 { // 256 color
				x = tcell.PaletteColor(params[pos+2] + 1)
				pos += 2 // two additional (5+index)
			} else if mode == 2 { // 24-bit
				x = tcell.NewRGBColor(int32(params[pos+2]), int32(params[pos+3]), int32(params[pos+4]))
				pos += 4 // four additional (2+r/g/b)
			} else {
				return errCSIParseError // invalid mode
			}
			if p == 38 {
				ei.curStyle = ei.curStyle.Foreground(x)
			} else {
				ei.curStyle = ei.curStyle.Background(x)
			}
		}

		pos += 1 // move along 1 by default
	}
	return nil
}
