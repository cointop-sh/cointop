// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

var (
	// ErrQuit is used to decide if the MainLoop finished successfully.
	ErrQuit = errors.New("quit")

	// ErrUnknownView allows to assert if a View must be initialized.
	ErrUnknownView = errors.New("unknown view")
)

// OutputMode represents the terminal's output mode (8 or 256 colors).
// type OutputMode termbox.OutputMode // TODO: die

// const ( // TODO: die
// 	// OutputNormal provides 8-colors terminal mode.
// 	OutputNormal = OutputMode(termbox.OutputNormal)

// 	// Output256 provides 256-colors terminal mode.
// 	Output256 = OutputMode(termbox.Output256)
// )

// Gui represents the whole User Interface, including the views, layouts
// and eventBindings.
type Gui struct {
	tbEvents      chan tcell.Event
	userEvents    chan userEvent
	views         []*View
	currentView   *View
	managers      []Manager
	eventBindings []*eventBinding
	maxX, maxY    int
	// outputMode    OutputMode // TODO: die
	screen tcell.Screen

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the GUI.
	BgColor, FgColor tcell.Color

	// SelBgColor and SelFgColor allow to configure the background and
	// foreground colors of the frame of the current view.
	SelBgColor, SelFgColor tcell.Color

	// If Highlight is true, Sel{Bg,Fg}Colors will be used to draw the
	// frame of the current view.
	Highlight bool

	// If Cursor is true then the cursor is enabled.
	Cursor bool

	// If Mouse is true then mouse events will be enabled.
	Mouse bool

	// If InputEsc is true, when ESC sequence is in the buffer and it doesn't
	// match any known sequence, ESC means KeyEsc.
	InputEsc bool

	// If ASCII is true then use ASCII instead of unicode to draw the
	// interface. Using ASCII is more portable.
	ASCII bool

	// The current event while in the handlers.
	CurrentEvent tcell.Event
}

// NewGui returns a new Gui object with a given output mode.
// func NewGui(mode OutputMode) (*Gui, error) {
func NewGui() (*Gui, error) {
	g := &Gui{}

	// outMode = OutputNormal
	if s, e := tcell.NewScreen(); e != nil {
		return nil, e
	} else if e = s.Init(); e != nil {
		return nil, e
	} else {
		g.screen = s
	}

	// g.outputMode = mode
	// termbox.SetScreen(g.Screen) // ugly global
	// termbox.SetOutputMode(termbox.OutputMode(mode))

	g.tbEvents = make(chan tcell.Event, 20)
	g.userEvents = make(chan userEvent, 20)

	g.maxX, g.maxY = g.screen.Size()

	g.BgColor, g.FgColor = tcell.ColorDefault, tcell.ColorDefault
	g.SelBgColor, g.SelFgColor = tcell.ColorDefault, tcell.ColorDefault

	return g, nil
}

// Close finalizes the library. It should be called after a successful
// initialization and when gocui is not needed anymore.
func (g *Gui) Close() {
	g.screen.Fini()
}

// Size returns the terminal's size.
func (g *Gui) Size() (x, y int) {
	return g.maxX, g.maxY
}

// temporary kludge for the pretty
func (g *Gui) prettyColor(x, y int, st tcell.Style) tcell.Style {
	if true {
		w, h := g.screen.Size()

		// dark blue gradient background
		red := int32(0)
		grn := int32(0)
		blu := int32(50 * float64(y) / float64(h))
		st = st.Background(tcell.NewRGBColor(red, grn, blu))

		// two-axis green-blue gradient
		red = int32(200)
		grn = int32(255 * float64(y) / float64(h))
		blu = int32(255 * float64(x) / float64(w))
		st = st.Foreground(tcell.NewRGBColor(red, grn, blu))
	}
	return st
}

// SetRune writes a rune at the given point, relative to the top-left
// corner of the terminal. It checks if the position is valid and applies
// the given colors.
func (g *Gui) SetRune(x, y int, ch rune, st tcell.Style) error {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return errors.New("invalid point")
	}
	// temporary kludge for the pretty
	// st = g.prettyColor(x, y, st)
	g.screen.SetContent(x, y, ch, nil, st)
	return nil
}

// Rune returns the rune contained in the cell at the given position.
// It checks if the position is valid.
// func (g *Gui) Rune(x, y int) (rune, error) {
// 	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
// 		return ' ', errors.New("invalid point")
// 	}
// 	c := termbox.CellBuffer()[y*g.maxX+x]
// 	return c.Ch, nil
// }

// SetView creates a new view with its top-left corner at (x0, y0)
// and the bottom-right one at (x1, y1). If a view with the same name
// already exists, its dimensions are updated; otherwise, the error
// ErrUnknownView is returned, which allows to assert if the View must
// be initialized. It checks if the position is valid.
func (g *Gui) SetView(name string, x0, y0, x1, y1 int) (*View, error) {
	if x0 >= x1 || y0 >= y1 {
		return nil, errors.New("invalid dimensions")
	}
	if name == "" {
		return nil, errors.New("invalid name")
	}

	if v, err := g.View(name); err == nil {
		v.x0 = x0
		v.y0 = y0
		v.x1 = x1
		v.y1 = y1
		v.tainted = true
		return v, nil
	}

	v := newView(name, x0, y0, x1, y1, g)
	v.BgColor, v.FgColor = g.BgColor, g.FgColor
	v.SelBgColor, v.SelFgColor = g.SelBgColor, g.SelFgColor
	g.views = append(g.views, v)
	return v, ErrUnknownView
}

// SetViewOnTop sets the given view on top of the existing ones.
func (g *Gui) SetViewOnTop(name string) (*View, error) {
	for i, v := range g.views {
		if v.name == name {
			s := append(g.views[:i], g.views[i+1:]...)
			g.views = append(s, v)
			return v, nil
		}
	}
	return nil, ErrUnknownView
}

// SetViewOnBottom sets the given view on bottom of the existing ones.
func (g *Gui) SetViewOnBottom(name string) (*View, error) {
	for i, v := range g.views {
		if v.name == name {
			s := append(g.views[:i], g.views[i+1:]...)
			g.views = append([]*View{v}, s...)
			return v, nil
		}
	}
	return nil, ErrUnknownView
}

// Views returns all the views in the GUI.
func (g *Gui) Views() []*View {
	return g.views
}

// View returns a pointer to the view with the given name, or error
// ErrUnknownView if a view with that name does not exist.
func (g *Gui) View(name string) (*View, error) {
	for _, v := range g.views {
		if v.name == name {
			return v, nil
		}
	}
	return nil, ErrUnknownView
}

// ViewByPosition returns a pointer to a view matching the given position, or
// error ErrUnknownView if a view in that position does not exist.
func (g *Gui) ViewByPosition(x, y int) (*View, error) {
	// traverse views in reverse order checking top views first
	for i := len(g.views); i > 0; i-- {
		v := g.views[i-1]
		if x > v.x0 && x < v.x1 && y > v.y0 && y < v.y1 {
			return v, nil
		}
	}
	return nil, ErrUnknownView
}

// ViewPosition returns the coordinates of the view with the given name, or
// error ErrUnknownView if a view with that name does not exist.
func (g *Gui) ViewPosition(name string) (x0, y0, x1, y1 int, err error) {
	for _, v := range g.views {
		if v.name == name {
			return v.x0, v.y0, v.x1, v.y1, nil
		}
	}
	return 0, 0, 0, 0, ErrUnknownView
}

// DeleteView deletes a view by name.
func (g *Gui) DeleteView(name string) error {
	for i, v := range g.views {
		if v.name == name {
			g.views = append(g.views[:i], g.views[i+1:]...)
			return nil
		}
	}
	return ErrUnknownView
}

// SetCurrentView gives the focus to a given view.
func (g *Gui) SetCurrentView(name string) (*View, error) {
	for _, v := range g.views {
		if v.name == name {
			g.currentView = v
			return v, nil
		}
	}
	return nil, ErrUnknownView
}

// CurrentView returns the currently focused view, or nil if no view
// owns the focus.
func (g *Gui) CurrentView() *View {
	return g.currentView
}

// SetKeybinding creates a new eventBinding. If viewname equals to ""
// (empty string) then the eventBinding will apply to all views. key must
// be a rune or a Key.
// TODO: split into key/mouse bindings?
func (g *Gui) SetKeybinding(viewname string, key tcell.Key, ch rune, mod tcell.ModMask, handler func(*Gui, *View) error) error {
	// var kb *eventBinding

	// k, ch, err := getKey(key)
	// if err != nil {
	// 	return err
	// }
	// TODO: get rid of this ugly mess
	//switch key {
	//case termbox.MouseLeft:
	//	kb = newMouseBinding(viewname, tcell.Button1, mod, handler)
	//case termbox.MouseMiddle:
	//	kb = newMouseBinding(viewname, tcell.Button3, mod, handler)
	//case termbox.MouseRight:
	//	kb = newMouseBinding(viewname, tcell.Button2, mod, handler)
	//case termbox.MouseWheelUp:
	//	kb = newMouseBinding(viewname, tcell.WheelUp, mod, handler)
	//case termbox.MouseWheelDown:
	//	kb = newMouseBinding(viewname, tcell.WheelDown, mod, handler)
	//default:
	//	kb = newKeybinding(viewname, key, ch, mod, handler)
	//}
	kb := newKeybinding(viewname, key, ch, mod, handler)
	g.eventBindings = append(g.eventBindings, kb)
	return nil
}

func (g *Gui) SetMousebinding(viewname string, btn tcell.ButtonMask, mod tcell.ModMask, handler func(*Gui, *View) error) error {
	kb := newMouseBinding(viewname, btn, mod, handler)
	g.eventBindings = append(g.eventBindings, kb)
	return nil
}

// DeleteKeybinding deletes a eventBinding.
func (g *Gui) DeleteKeybinding(viewname string, key tcell.Key, ch rune, mod tcell.ModMask) error {
	// k, ch, err := getKey(key)
	// if err != nil {
	// 	return err
	// }

	for i, kb := range g.eventBindings {
		if kbe, ok := kb.ev.(*tcell.EventKey); ok {
			if kb.viewName == viewname && kbe.Rune() == ch && kbe.Key() == key && kbe.Modifiers() == mod {
				g.eventBindings = append(g.eventBindings[:i], g.eventBindings[i+1:]...)
				return nil
			}
		}
	}
	return errors.New("eventBinding not found")
}

// DeleteKeybindings deletes all eventBindings of view.
func (g *Gui) DeleteKeybindings(viewname string) {
	var s []*eventBinding
	for _, kb := range g.eventBindings {
		if kb.viewName != viewname {
			s = append(s, kb)
		}
	}
	g.eventBindings = s
}

// getKey takes an empty interface with a key and returns the corresponding
// typed Key or rune.
// func getKey(key interface{}) (tcell.Key, rune, error) {
// 	switch t := key.(type) {
// 	case Key:
// 		return t, 0, nil
// 	case rune:
// 		return 0, t, nil
// 	default:
// 		return 0, 0, errors.New("unknown type")
// 	}
// }

// userEvent represents an event triggered by the user.
type userEvent struct {
	f func(*Gui) error
}

// Update executes the passed function. This method can be called safely from a
// goroutine in order to update the GUI. It is important to note that the
// passed function won't be executed immediately, instead it will be added to
// the user events queue. Given that Update spawns a goroutine, the order in
// which the user events will be handled is not guaranteed.
func (g *Gui) Update(f func(*Gui) error) {
	go func() { g.userEvents <- userEvent{f: f} }()
}

// A Manager is in charge of GUI's layout and can be used to build widgets.
type Manager interface {
	// Layout is called every time the GUI is redrawn, it must contain the
	// base views and its initializations.
	Layout(*Gui) error
}

// The ManagerFunc type is an adapter to allow the use of ordinary functions as
// Managers. If f is a function with the appropriate signature, ManagerFunc(f)
// is an Manager object that calls f.
type ManagerFunc func(*Gui) error

// Layout calls f(g)
func (f ManagerFunc) Layout(g *Gui) error {
	return f(g)
}

// SetManager sets the given GUI managers. It deletes all views and
// eventBindings.
func (g *Gui) SetManager(managers ...Manager) {
	g.managers = managers
	g.currentView = nil
	g.views = nil
	g.eventBindings = nil

	go func() { g.tbEvents <- tcell.NewEventResize(0, 0) }()
}

// SetManagerFunc sets the given manager function. It deletes all views and
// eventBindings.
func (g *Gui) SetManagerFunc(manager func(*Gui) error) {
	g.SetManager(ManagerFunc(manager))
}

// MainLoop runs the main loop until an error is returned. A successful
// finish should return ErrQuit.
func (g *Gui) MainLoop() error {
	go func() {
		for {
			g.tbEvents <- g.screen.PollEvent()
		}
	}()

	if g.Mouse {
		g.screen.EnableMouse()
	}
	// s.EnablePaste()

	if err := g.flush(); err != nil {
		return err
	}
	for {
		select {
		case ev := <-g.tbEvents:
			if err := g.handleEvent(ev); err != nil {
				return err
			}
		case ev := <-g.userEvents:
			if err := ev.f(g); err != nil {
				return err
			}
		}
		if err := g.consumeevents(); err != nil {
			return err
		}
		if err := g.flush(); err != nil {
			return err
		}
	}
}

// consumeevents handles the remaining events in the events pool.
func (g *Gui) consumeevents() error {
	for {
		select {
		case ev := <-g.tbEvents:
			if err := g.handleEvent(ev); err != nil {
				return err
			}
		case ev := <-g.userEvents:
			if err := ev.f(g); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// handleEvent handles an event, based on its type (key-press, error,
// etc.)
func (g *Gui) handleEvent(ev tcell.Event) error {
	switch tev := ev.(type) {
	case *tcell.EventMouse, *tcell.EventKey:
		return g.onEvent(tev)
	case *tcell.EventError:
		return errors.New(tev.Error())
	default:
		return nil
	}
}

// TODO: delete termbox compat
func (g *Gui) fixColor(c tcell.Color) tcell.Color {
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

func (g *Gui) Style(fg, bg tcell.Color) tcell.Style {
	st := tcell.StyleDefault.Foreground(fg).Background(bg)
  	return st
}

func (g *Gui) MkColor(color Attribute) tcell.Color {
	if color == ColorDefault {
		return tcell.ColorDefault
	} else {
		return g.fixColor(tcell.PaletteColor(int(color)&0x1ff - 1))
	}
}

// TODO: delete termbox compat
func (g *Gui) MkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault
	if fg != ColorDefault {
		f := tcell.PaletteColor(int(fg)&0x1ff - 1)
		f = g.fixColor(f)
		st = st.Foreground(f)
	}
	if bg != ColorDefault {
		b := tcell.PaletteColor(int(bg)&0x1ff - 1)
		b = g.fixColor(b)
		st = st.Background(b)
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

// flush updates the gui, re-drawing frames and buffers.
func (g *Gui) flush() error {
	// termbox.Clear(termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
	st := g.Style(g.FgColor, g.BgColor)
	w, h := g.screen.Size() // TODO: merge with maxX, maxY below
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			g.screen.SetContent(col, row, ' ', nil, st)
		}
	}

	maxX, maxY := g.screen.Size()
	// if GUI's size has changed, we need to redraw all views
	if maxX != g.maxX || maxY != g.maxY {
		for _, v := range g.views {
			v.tainted = true
		}
	}
	g.maxX, g.maxY = maxX, maxY

	for _, m := range g.managers {
		if err := m.Layout(g); err != nil {
			return err
		}
	}
	for _, v := range g.views {
		if v.Frame {
			st := g.Style(v.FgColor, v.BgColor)
			if g.Highlight && v == g.currentView {
				st = g.Style(g.SelFgColor, g.SelBgColor)
			}

			if err := g.drawFrameEdges(v, st); err != nil {
				return err
			}
			if err := g.drawFrameCorners(v, st); err != nil {
				return err
			}
			if v.Title != "" {
				if err := g.drawTitle(v, st); err != nil {
					return err
				}
			}
		}
		if err := g.draw(v); err != nil {
			return err
		}
	}
	g.screen.Show()
	return nil
}

// drawFrameEdges draws the horizontal and vertical edges of a view.
func (g *Gui) drawFrameEdges(v *View, st tcell.Style) error {
	runeH, runeV := '─', '│'
	if g.ASCII {
		runeH, runeV = '-', '|'
	}

	for x := v.x0 + 1; x < v.x1 && x < g.maxX; x++ {
		if x < 0 {
			continue
		}
		if v.y0 > -1 && v.y0 < g.maxY {
			if err := g.SetRune(x, v.y0, runeH, st); err != nil {
				return err
			}
		}
		if v.y1 > -1 && v.y1 < g.maxY {
			if err := g.SetRune(x, v.y1, runeH, st); err != nil {
				return err
			}
		}
	}
	for y := v.y0 + 1; y < v.y1 && y < g.maxY; y++ {
		if y < 0 {
			continue
		}
		if v.x0 > -1 && v.x0 < g.maxX {
			if err := g.SetRune(v.x0, y, runeV, st); err != nil {
				return err
			}
		}
		if v.x1 > -1 && v.x1 < g.maxX {
			if err := g.SetRune(v.x1, y, runeV, st); err != nil {
				return err
			}
		}
	}
	return nil
}

// drawFrameCorners draws the corners of the view.
func (g *Gui) drawFrameCorners(v *View, st tcell.Style) error {
	runeTL, runeTR, runeBL, runeBR := '┌', '┐', '└', '┘'
	if g.ASCII {
		runeTL, runeTR, runeBL, runeBR = '+', '+', '+', '+'
	}

	corners := []struct {
		x, y int
		ch   rune
	}{{v.x0, v.y0, runeTL}, {v.x1, v.y0, runeTR}, {v.x0, v.y1, runeBL}, {v.x1, v.y1, runeBR}}

	for _, c := range corners {
		if c.x >= 0 && c.y >= 0 && c.x < g.maxX && c.y < g.maxY {
			if err := g.SetRune(c.x, c.y, c.ch, st); err != nil {
				return err
			}
		}
	}
	return nil
}

// drawTitle draws the title of the view.
func (g *Gui) drawTitle(v *View, st tcell.Style) error {
	if v.y0 < 0 || v.y0 >= g.maxY {
		return nil
	}

	for i, ch := range v.Title {
		x := v.x0 + i + 2
		if x < 0 {
			continue
		} else if x > v.x1-2 || x >= g.maxX {
			break
		}
		if err := g.SetRune(x, v.y0, ch, st); err != nil {
			return err
		}
	}
	return nil
}

// draw manages the cursor and calls the draw function of a view.
func (g *Gui) draw(v *View) error {
	if g.Cursor {
		if curview := g.currentView; curview != nil {
			vMaxX, vMaxY := curview.Size()
			if curview.cx < 0 {
				curview.cx = 0
			} else if curview.cx >= vMaxX {
				curview.cx = vMaxX - 1
			}
			if curview.cy < 0 {
				curview.cy = 0
			} else if curview.cy >= vMaxY {
				curview.cy = vMaxY - 1
			}

			gMaxX, gMaxY := g.Size()
			cx, cy := curview.x0+curview.cx+1, curview.y0+curview.cy+1
			if cx >= 0 && cx < gMaxX && cy >= 0 && cy < gMaxY {
				g.screen.ShowCursor(cx, cy)
			} else {
				g.screen.ShowCursor(-1, -1) // HideCursor
			}
		}
	} else {
		g.screen.ShowCursor(-1, -1) // HideCursor
	}

	v.clearRunes()
	if err := v.draw(); err != nil {
		return err
	}
	return nil
}

// onEvent manages key/mouse events. A eventBinding handler is called when
// a key-press or mouse event satisfies a configured eventBinding. Furthermore,
// currentView's internal buffer is modified if currentView.Editable is true.
func (g *Gui) onEvent(ev tcell.Event) error {
	switch tev := ev.(type) {
	case *tcell.EventKey:
		matched, err := g.execEventBindings(g.currentView, ev)
		if err != nil {
			return err
		}
		if matched {
			break
		}
		if g.currentView != nil && g.currentView.Editable && g.currentView.Editor != nil {
			g.currentView.Editor.Edit(g.currentView, tev.Key(), tev.Rune(), tev.Modifiers())
		}
	case *tcell.EventMouse:
		v, _, _, err := g.GetViewRelativeMousePosition(tev)
		if err != nil {
			break
		}
		// If the key-binding wants to move the cursor, it should call SetCursorFromCurrentMouseEvent()
		// Not all mouse events will want to do this (eg: scroll wheel)
		g.CurrentEvent = ev
		if _, err := g.execEventBindings(v, g.CurrentEvent); err != nil {
			return err
		}
	}
	return nil
}

// GetViewRelativeMousePosition returns the View and relative x/y for the provided mouse event.
func (g *Gui) GetViewRelativeMousePosition(ev tcell.Event) (*View, int, int, error) {
	if kbe, ok := ev.(*tcell.EventMouse); ok {
		mx, my := kbe.Position()
		v, err := g.ViewByPosition(mx, my)
		if err != nil {
			return nil, 0, 0, err
		}
		return v, mx - v.x0 - 1, my - v.y0 - 1, nil
	}
	return nil, 0, 0, errors.New("Cannot GetViewRelativeMousePosition on non-mouse event")
}

// SetCursorFromCurrentMouseEvent updates the cursor position based on the mouse coordinates.
func (g *Gui) SetCursorFromCurrentMouseEvent() error {
	v, x, y, err := g.GetViewRelativeMousePosition(g.CurrentEvent)
	if err != nil {
		return err
	}
	if err := v.SetCursor(x, y); err != nil {
		return err
	}
	return nil
}

// execEventBindings executes the handlers that match the passed view
// and event. The value of matched is true if there is a match and no errors.
// TODO: rename to more generic - it's not just keys (incl mouse)
func (g *Gui) execEventBindings(v *View, xev tcell.Event) (matched bool, err error) {
	matched = false
	for _, kb := range g.eventBindings {
		if kb.handler == nil {
			continue
		}
		if kb.matchEvent(xev) && kb.matchView(v) {
			if err := kb.handler(g, v); err != nil {
				return false, err
			}
			matched = true
		}
	}
	return matched, nil
}
