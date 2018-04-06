package cointop

import (
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/open"
)

func (ct *Cointop) openHelp(g *gocui.Gui, v *gocui.View) error {
	open.URL(ct.helpLink())
	return nil
}

func (ct *Cointop) helpLink() string {
	return "https://github.com/miguelmota/cointop#shortcuts"
}
