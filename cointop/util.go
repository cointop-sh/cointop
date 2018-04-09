package cointop

import "github.com/miguelmota/cointop/pkg/open"

func (ct *Cointop) openLink() error {
	open.URL(ct.rowLink())
	return nil
}
