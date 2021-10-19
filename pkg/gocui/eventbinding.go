// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"github.com/gdamore/tcell/v2"
)

// eventBinding are used to link a given key-press event with a handler.
type eventBinding struct {
	viewName string
	ev       tcell.Event // ignore the Time
	handler  func(*Gui, *View) error
}

// newKeybinding returns a new eventBinding object for a key event.
func newKeybinding(viewname string, key tcell.Key, ch rune, mod tcell.ModMask, handler func(*Gui, *View) error) (kb *eventBinding) {
	kb = &eventBinding{
		viewName: viewname,
		ev:       tcell.NewEventKey(key, ch, mod),
		handler:  handler,
	}
	return kb
}

// newKeybinding returns a new eventBinding object for a mouse event.
func newMouseBinding(viewname string, btn tcell.ButtonMask, mod tcell.ModMask, handler func(*Gui, *View) error) (kb *eventBinding) {
	kb = &eventBinding{
		viewName: viewname,
		ev:       tcell.NewEventMouse(0, 0, btn, mod),
		handler:  handler,
	}
	return kb
}

func (kb *eventBinding) matchEvent(e tcell.Event) bool {
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

// matchView returns if the eventBinding matches the current view.
func (kb *eventBinding) matchView(v *View) bool {
	if kb.viewName == "" {
		return true
	}
	return v != nil && kb.viewName == v.name
}
