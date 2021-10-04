package cointop

import (
	"github.com/cointop-sh/cointop/pkg/ui"
	log "github.com/sirupsen/logrus"
)

// MenuView is structure for menu view
type MenuView = ui.View

// NewMenuView returns a new menu view
func NewMenuView() *MenuView {
	return ui.NewView("menu")
}

// HideMenu hides the menu view
func (ct *Cointop) HideMenu() error {
	log.Debug("HideMenu()")
	return nil
}
