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
	key      tcell.Key
	ch       rune
	mod      tcell.ModMask
	handler  func(*Gui, *View) error
}

// newKeybinding returns a new Keybinding object.
func newKeybinding(viewname string, key tcell.Key, ch rune, mod tcell.ModMask, handler func(*Gui, *View) error) (kb *keybinding) {
	kb = &keybinding{
		viewName: viewname,
		key:      key,
		ch:       ch,
		mod:      mod,
		handler:  handler,
	}
	return kb
}

// matchKeypress returns if the keybinding matches the keypress.
func (kb *keybinding) matchKeypress(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	// TODO: check mask not ==mod?
	if key == tcell.KeyRune {
		return kb.key == key && kb.ch == ch && kb.mod == mod
	}
	return kb.key == key && kb.mod == mod
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
	// KeyF1         Key = tcell.KeyF1
	// KeyF2             = tcell.KeyF2
	// KeyF3             = tcell.KeyF3
	// KeyF4             = tcell.KeyF4
	// KeyF5             = tcell.KeyF5
	// KeyF6             = tcell.KeyF6
	// KeyF7             = tcell.KeyF7
	// KeyF8             = tcell.KeyF8
	// KeyF9             = tcell.KeyF9
	// KeyF10            = tcell.KeyF10
	// KeyF11            = tcell.KeyF11
	// KeyF12            = tcell.KeyF12
	// KeyInsert         = tcell.KeyInsert
	// KeyDelete         = tcell.KeyDelete
	// KeyHome           = tcell.KeyHome
	// KeyEnd            = tcell.KeyEnd
	// KeyPgup           = tcell.KeyPgup
	// KeyPgdn           = tcell.KeyPgdn
	// KeyArrowUp        = tcell.KeyArrowUp
	// KeyArrowDown      = tcell.KeyArrowDown
	// KeyArrowLeft      = tcell.KeyArrowLeft
	// KeyArrowRight     = tcell.KeyArrowRight

	MouseLeft      = termbox.MouseLeft
	MouseMiddle    = termbox.MouseMiddle
	MouseRight     = termbox.MouseRight
	MouseRelease   = termbox.MouseRelease
	MouseWheelUp   = termbox.MouseWheelUp
	MouseWheelDown = termbox.MouseWheelDown
)

// Keys combinations.
// const (
// 	KeyCtrlTilde      tcell.Key = tcell.KeyCtrlTilde
// 	KeyCtrl2              = tcell.KeyCtrl2
// 	KeyCtrlSpace          = tcell.KeyCtrlSpace
// 	KeyCtrlA              = tcell.KeyCtrlA
// 	KeyCtrlB              = tcell.KeyCtrlB
// 	KeyCtrlC              = tcell.KeyCtrlC
// 	KeyCtrlD              = tcell.KeyCtrlD
// 	KeyCtrlE              = tcell.KeyCtrlE
// 	KeyCtrlF              = tcell.KeyCtrlF
// 	KeyCtrlG              = tcell.KeyCtrlG
// 	KeyBackspace          = tcell.KeyBackspace
// 	KeyCtrlH              = tcell.KeyCtrlH
// 	KeyTab                = tcell.KeyTab
// 	KeyCtrlI              = tcell.KeyCtrlI
// 	KeyCtrlJ              = tcell.KeyCtrlJ
// 	KeyCtrlK              = tcell.KeyCtrlK
// 	KeyCtrlL              = tcell.KeyCtrlL
// 	KeyEnter              = tcell.KeyEnter
// 	KeyCtrlM              = tcell.KeyCtrlM
// 	KeyCtrlN              = tcell.KeyCtrlN
// 	KeyCtrlO              = tcell.KeyCtrlO
// 	KeyCtrlP              = tcell.KeyCtrlP
// 	KeyCtrlQ              = tcell.KeyCtrlQ
// 	KeyCtrlR              = tcell.KeyCtrlR
// 	KeyCtrlS              = tcell.KeyCtrlS
// 	KeyCtrlT              = tcell.KeyCtrlT
// 	KeyCtrlU              = tcell.KeyCtrlU
// 	KeyCtrlV              = tcell.KeyCtrlV
// 	KeyCtrlW              = tcell.KeyCtrlW
// 	KeyCtrlX              = tcell.KeyCtrlX
// 	KeyCtrlY              = tcell.KeyCtrlY
// 	KeyCtrlZ              = tcell.KeyCtrlZ
// 	KeyEsc                = tcell.KeyEsc
// 	KeyCtrlLsqBracket     = tcell.KeyCtrlLsqBracket
// 	KeyCtrl3              = tcell.KeyCtrl3
// 	KeyCtrl4              = tcell.KeyCtrl4
// 	KeyCtrlBackslash      = tcell.KeyCtrlBackslash
// 	KeyCtrl5              = tcell.KeyCtrl5
// 	KeyCtrlRsqBracket     = tcell.KeyCtrlRsqBracket
// 	KeyCtrl6              = tcell.KeyCtrl6
// 	KeyCtrl7              = tcell.KeyCtrl7
// 	KeyCtrlSlash          = tcell.KeyCtrlSlash
// 	KeyCtrlUnderscore     = tcell.KeyCtrlUnderscore
// 	KeySpace              = tcell.KeySpace
// 	KeyBackspace2         = tcell.KeyBackspace2
// 	KeyCtrl8              = tcell.KeyCtrl8
// )

// Modifier allows to define special keys combinations. They can be used
// in combination with Keys or Runes when a new keybinding is defined.
// type Modifier tcell.ModMask

// Modifiers.
// const (
// 	ModNone Modifier = Modifier(0)
// 	ModAlt           = Modifier(termbox.ModAlt)
// )
