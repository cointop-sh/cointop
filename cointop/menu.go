package cointop

import "github.com/miguelmota/cointop/pkg/ui"

// MenuView is structure for menu view
type MenuView = ui.View

// NewMenuView returns a new menu view
func NewMenuView() *MenuView {
	var view *MenuView = ui.NewView("menu")
	return view
}

// HideMenu hides the menu view
func (ct *Cointop) HideMenu() error {
	ct.debuglog("HideMenu()")
	return nil
}
