package cointop

import (
	"sync"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

var updateMutex sync.Mutex

// update takes a callback which updates the view
func (ct *Cointop) update(f func()) {
	updateMutex.Lock()
	defer updateMutex.Unlock()
	if ct.g == nil {
		log.Fatal("gocui is not initialized")
	}

	ct.g.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}
