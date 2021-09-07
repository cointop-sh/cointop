package cointop

import log "github.com/sirupsen/logrus"

// RowChanged is called when the row is updated
func (ct *Cointop) RowChanged() {
	log.Debug("RowChanged()")
	ct.RefreshRowLink()
}
