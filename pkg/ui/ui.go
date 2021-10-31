package ui

import (
	"github.com/cointop-sh/cointop/pkg/gocui"
	"github.com/cointop-sh/cointop/pkg/termbox"
)

// UI is the UI view struct
type UI struct {
	g *gocui.Gui
}

// NewUI returns a new UI instance
func NewUI() (*UI, error) {
	g, err := gocui.NewGui()
	if err != nil {
		return nil, err
	}

	return &UI{
		g: g,
	}, nil
}

// GetGocui returns the underlying gocui instance
func (ui *UI) GetGocui() *gocui.Gui {
	return ui.g
}

// SetFgColor sets the foreground color
func (ui *UI) SetFgColor(fgColor termbox.Attribute) {
	ui.g.Style = ui.g.Style.Foreground(termbox.MkColor(fgColor))
}

// SetBgColor sets the background color
func (ui *UI) SetBgColor(bgColor termbox.Attribute) {
	ui.g.Style = ui.g.Style.Background(termbox.MkColor(bgColor))
}

// SetInputEsc enables the escape key
func (ui *UI) SetInputEsc(enabled bool) {
	ui.g.InputEsc = enabled
}

// SetMouse enables the mouse
func (ui *UI) SetMouse(enabled bool) {
	ui.g.Mouse = enabled
}

// SetCursor enables the input field cursor
func (ui *UI) SetCursor(enabled bool) {
	ui.g.Cursor = enabled
}

// SetHighlight enables the highlight active state
func (ui *UI) SetHighlight(enabled bool) {
	ui.g.Highlight = enabled
}

// SetManagerFunc sets the function to call for rendering UI
func (ui *UI) SetManagerFunc(fn func() error) {
	ui.g.SetManagerFunc(func(*gocui.Gui) error {
		return fn()
	})
}

// MainLoop starts the UI render loop
func (ui *UI) MainLoop() error {
	return ui.g.MainLoop()
}

// Close ...
func (ui *UI) Close() {
	ui.g.Close()
}

// SetView sets the view layout
func (ui *UI) SetView(view interface{}, x, y, w, h int) error {
	if v, ok := view.(*View); ok {
		gv, err := ui.g.SetView(v.Name(), x, y, w, h)
		v.SetBacking(gv)
		return err
	}
	return nil
}

// SetViewOnBottom sets the view to the bottom layer
func (ui *UI) SetViewOnBottom(view interface{}) error {
	if v, ok := view.(*View); ok {
		if _, err := ui.g.SetViewOnBottom(v.Name()); err != nil {
			return err
		}
	}
	return nil
}

// ActiveViewName returns the name of the active view
func (ui *UI) ActiveViewName() string {
	return ui.g.CurrentView().Name()
}
