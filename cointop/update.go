package cointop

import (
	"github.com/cointop-sh/cointop/pkg/gocui"
	log "github.com/sirupsen/logrus"
)

// UpdateUI takes a callback which updates the view
func (ct *Cointop) UpdateUI(f func() error) {
	log.Debug("UpdateUI()")

	if ct.g == nil {
		return
	}

	ct.g.Update(func(g *gocui.Gui) error {
		return f()
	})
}
