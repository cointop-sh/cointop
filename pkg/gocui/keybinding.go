// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"github.com/cointop-sh/cointop/pkg/termbox"
	"github.com/gdamore/tcell/v2"
)

// Keybidings are used to link a given key-press event with a handler.
type keybinding struct {
	viewName string
	ev       tcell.Event // ignore the Time
	handler  func(*Gui, *View) error
}

// newKeybinding returns a new Keybinding object.
func newKeybinding(viewname string, key tcell.Key, ch rune, mod tcell.ModMask, handler func(*Gui, *View) error) (kb *keybinding) {
	// TODO: take Event
	kb = &keybinding{
		viewName: viewname,
		ev:       tcell.NewEventKey(key, ch, mod),
		handler:  handler,
	}
	return kb
}

func newMouseBinding(viewname string, btn tcell.ButtonMask, mod tcell.ModMask, handler func(*Gui, *View) error) (kb *keybinding) {
	kb = &keybinding{
		viewName: viewname,
		ev:       tcell.NewEventMouse(0, 0, btn, mod),
		handler:  handler,
	}
	return kb
}

func (kb *keybinding) matchEvent(e tcell.Event) bool {
	// TODO: check mask not ==mod?
	switch tev := e.(type) {
	case *tcell.EventKey:
		if kbe, ok := kb.ev.(*tcell.EventKey); ok {
			if tev.Key() == tcell.KeyRune {
				return tev.Key() == kbe.Key() && tev.Rune() == kbe.Rune() && tev.Modifiers() == kbe.Modifiers()
			}
			return tev.Key() == kbe.Key() && tev.Modifiers() == kbe.Modifiers()
		}

	case *tcell.EventMouse:
		if kbe, ok := kb.ev.(*tcell.EventMouse); ok {
			return kbe.Buttons() == tev.Buttons() && kbe.Modifiers() == tev.Modifiers()
		}

	}
	return false
}

// matchView returns if the keybinding matches the current view.
func (kb *keybinding) matchView(v *View) bool {
	if kb.viewName == "" {
		return true
	}
	return v != nil && kb.viewName == v.name
}

// Key represents special keys or keys combinations.
// type Key tcell.Key

// Special keys.
const (
	MouseLeft      = termbox.MouseLeft
	MouseMiddle    = termbox.MouseMiddle
	MouseRight     = termbox.MouseRight
	MouseRelease   = termbox.MouseRelease
	MouseWheelUp   = termbox.MouseWheelUp
	MouseWheelDown = termbox.MouseWheelDown
)
