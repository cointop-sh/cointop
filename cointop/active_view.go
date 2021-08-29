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

	// TODO: better way to map/unmap shortcut key actions based on active view
	if v == ct.Views.Table.Name() {
		if err := ct.SetKeybindingAction("/", "open_search"); err != nil {
			return err
		}
	} else {
		// deletes binding to allow using "/" key on input fields
		if err := ct.DeleteKeybinding("/"); err != nil {
			return err
		}
	}
	return nil
}
