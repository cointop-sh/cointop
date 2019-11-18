package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// Update takes a callback which updates the view
func (ct *Cointop) Update(f func()) {
	ct.debuglog(fmt.Sprintf("Update()"))

	if ct.g == nil {
		log.Fatal("gocui is not initialized")
	}

	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
