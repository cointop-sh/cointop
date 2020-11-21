package ui

import (
	"github.com/miguelmota/gocui"
)

// UI is the UI view struct
type UI struct {
	g *gocui.Gui
}

// NewUI returns a new UI instance
func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.Output256)
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
func (ui *UI) SetFgColor(fgColor gocui.Attribute) {
	ui.g.FgColor = fgColor
}

// SetBgColor sets the background color
func (ui *UI) SetBgColor(bgColor gocui.Attribute) {
	ui.g.BgColor = bgColor
}

// SetInputEsc enables the escape key
func (ui *UI) SetInputEsc(enabled bool) {
	ui.g.InputEsc = true
}

// SetMouse enables the mouse
func (ui *UI) SetMouse(enabled bool) {
	ui.g.Mouse = true
}

// SetHighlight enables the highlight active state
func (ui *UI) SetHighlight(enabled bool) {
	ui.g.Highlight = true
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
