package ui

import (
	"github.com/miguelmota/gocui"
)

// UI ...
type UI struct {
	g *gocui.Gui
}

// NewUI ...
func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	return &UI{
		g: g,
	}, nil
}

// GetGocui ...
func (ui *UI) GetGocui() *gocui.Gui {
	return ui.g
}

// SetFgColor ...
func (ui *UI) SetFgColor(fgColor gocui.Attribute) {
	ui.g.FgColor = fgColor
}

// SetBgColor ...
func (ui *UI) SetBgColor(bgColor gocui.Attribute) {
	ui.g.BgColor = bgColor
}

// SetInputEsc ...
func (ui *UI) SetInputEsc(enabled bool) {
	ui.g.InputEsc = true
}

// SetMouse ...
func (ui *UI) SetMouse(enabled bool) {
	ui.g.Mouse = true
}

// SetHighlight ...
func (ui *UI) SetHighlight(enabled bool) {
	ui.g.Highlight = true
}

// SetManagerFunc ...
func (ui *UI) SetManagerFunc(fn func() error) {
	ui.g.SetManagerFunc(func(*gocui.Gui) error {
		return fn()
	})
}

// MainLoop ...
func (ui *UI) MainLoop() error {
	return ui.g.MainLoop()
}

// Close ...
func (ui *UI) Close() {
	ui.g.Close()
}

// SetView ...
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
