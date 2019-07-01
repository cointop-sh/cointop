package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// IView is a cointop view
type IView interface {
	Backing() *gocui.View
	SetBacking(gocuiView *gocui.View)
	Name() string
}

// View is a cointop view
type View struct {
	backing *gocui.View
	name    string
}

// NewView creates a new view
func NewView(name string) *View {
	return &View{
		name: name,
	}
}

// Backing returns the backing gocui view
func (view *View) Backing() *gocui.View {
	return view.backing
}

// SetBacking sets the backing gocui view
func (view *View) SetBacking(gocuiView *gocui.View) {
	view.backing = gocuiView
}

// Height returns thejview height
func (view *View) Height() int {
	_, h := view.backing.Size()
	return h
}

// Name returns the view's name
func (view *View) Name() string {
	return view.name
}

// SetActiveView sets the active view
func (ct *Cointop) SetActiveView(v string) error {
	ct.g.SetViewOnTop(v)
	ct.g.SetCurrentView(v)
	if v == ct.Views.SearchField.Name() {
		ct.Views.SearchField.Backing().Clear()
		ct.Views.SearchField.Backing().SetCursor(1, 0)
		fmt.Fprintf(ct.Views.SearchField.Backing(), "%s", "/")
	} else if v == ct.Views.Table.Name() {
		ct.g.SetViewOnTop(ct.Views.Statusbar.Name())
	}
	if v == ct.Views.PortfolioUpdateMenu.Name() {
		ct.g.SetViewOnTop(ct.Views.Input.Name())
		ct.g.SetCurrentView(ct.Views.Input.Name())
	}
	return nil
}

// ActiveViewName returns the name of the active view
func (ct *Cointop) ActiveViewName() string {
	return ct.g.CurrentView().Name()
}

// SetViewOnBottom sets the view to the bottom layer
func (ct *Cointop) SetViewOnBottom(v string) error {
	_, err := ct.g.SetViewOnBottom(v)
	return err
}
