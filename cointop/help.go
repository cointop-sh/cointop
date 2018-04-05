package cointop

import (
	"os/exec"

	"github.com/jroimartin/gocui"
)

func (ct *Cointop) openHelp(g *gocui.Gui, v *gocui.View) error {
	exec.Command("open", ct.helpLink()).Output()
	return nil
}

func (ct *Cointop) helpLink() string {
	return "https://github.com/miguelmota/cointop#shortcuts"
}
