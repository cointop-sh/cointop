package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// View is a cointop view
type View struct {
	Backing *gocui.View
	Name    string
}

func (ct *Cointop) setActiveView(v string) error {
	ct.g.SetViewOnTop(v)
	ct.g.SetCurrentView(v)
	if v == ct.Views.SearchField.Name {
		ct.Views.SearchField.Backing.Clear()
		ct.Views.SearchField.Backing.SetCursor(1, 0)
		fmt.Fprintf(ct.Views.SearchField.Backing, "%s", "/")
	} else if v == ct.Views.Table.Name {
		ct.g.SetViewOnTop(ct.Views.Statusbar.Name)
	}
	if v == ct.Views.PortfolioUpdateMenu.Name {
		ct.g.SetViewOnTop(ct.Views.Input.Name)
		ct.g.SetCurrentView(ct.Views.Input.Name)
	}
	return nil
}

func (ct *Cointop) activeViewName() string {
	return ct.g.CurrentView().Name()
}

func (ct *Cointop) setViewOnBottom(v string) error {
	_, err := ct.g.SetViewOnBottom(v)
	return err
}
