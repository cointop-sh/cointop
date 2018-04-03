package cointop

import "github.com/jroimartin/gocui"

// Update update view
func (ct *Cointop) Update(f func()) {
	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
