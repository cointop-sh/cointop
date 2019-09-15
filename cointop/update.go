package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// update takes a callback which updates the view
func (ct *Cointop) update(f func()) {
	ct.debuglog(fmt.Sprintf("update()"))

	if ct.g == nil {
		log.Fatal("gocui is not initialized")
	}

	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
