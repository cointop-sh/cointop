package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateStatus(s string) {
	maxX, _ := ct.g.Size()
	ct.g.Update(func(g *gocui.Gui) error {
		ct.statusview.Clear()
		fmt.Fprintln(ct.statusview, pad.Right(fmt.Sprintf("[q]uit [← →]page %s", s), maxX, " "))
		return nil
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.g.Update(func(g *gocui.Gui) error {
		ct.updateStatus(fmt.Sprintf("[↵]%s", url))
		return nil
	})
}
