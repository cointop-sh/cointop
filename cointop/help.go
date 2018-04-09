package cointop

import (
	"github.com/miguelmota/cointop/pkg/open"
)

// TODO: create a help menu
func (ct *Cointop) openHelp() error {
	open.URL(ct.helpLink())
	return nil
}

func (ct *Cointop) helpLink() string {
	return "https://github.com/miguelmota/cointop#shortcuts"
}
