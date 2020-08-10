package cointop

import (
	"fmt"

	"github.com/miguelmota/gocui"
)

// Update takes a callback which updates the view
func (ct *Cointop) Update(f func() error) {
	ct.debuglog(fmt.Sprintf("Update()"))

	if ct.g == nil {
		return
	}

	ct.g.Update(func(g *gocui.Gui) error {
		return f()
	})
}
