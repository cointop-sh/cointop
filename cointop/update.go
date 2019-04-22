package cointop

import (
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// update update view
func (ct *Cointop) update(f func()) {
	if ct.g == nil {
		log.Fatal("gocui is not initialized")
	}

	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
