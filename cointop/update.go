package cointop

import (
	"fmt"

	"github.com/miguelmota/gocui"
)

// Update takes a callback which updates the view
func (ct *Cointop) Update(f func()) {
	ct.debuglog(fmt.Sprintf("Update()"))

	if ct.g == nil {
		panic("gocui is not initialized")
	}

	ct.g.Update(func(g *gocui.Gui) error {
		f()

		return nil
	})
}
