package cointop

// SetActiveView sets the active view
func (ct *Cointop) SetActiveView(v string) error {
	ct.g.SetViewOnTop(v)
	ct.g.SetCurrentView(v)
	if v == ct.Views.SearchField.Name() {
		ct.Views.SearchField.SetCursor(1, 0)
		ct.Views.SearchField.Update("/")
	} else if v == ct.Views.Table.Name() {
		ct.g.SetViewOnTop(ct.Views.Statusbar.Name())
	}
	return nil
}
