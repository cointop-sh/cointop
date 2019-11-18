package cointop

// RowChanged is called when the row is updated
func (ct *Cointop) RowChanged() {
	ct.debuglog("RowChanged()")
	ct.RefreshRowLink()
}
