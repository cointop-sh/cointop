package cointop

import (
	"github.com/jroimartin/gocui"
)

// update update view
func (ct *Cointop) update(f func()) {
	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
