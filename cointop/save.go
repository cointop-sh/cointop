package cointop

// Save saves the cointop settings to the config file
func (ct *Cointop) Save() error {
	ct.debuglog("Save()")
	ct.SetSavingStatus()
	if err := ct.saveConfig(); err != nil {
		return err
	}

	ct.CacheAllCoinsSlugMap()

	return nil
}

// SetSavingStatus sets the saving indicator in the statusbar
func (ct *Cointop) SetSavingStatus() {
	ct.debuglog("SetSavingStatus()")
	if ct.g == nil {
		return
	}

	go func() {
		ct.loadingTicks("saving", 590)
		ct.UpdateStatusbar("")
		ct.RowChanged()
	}()
}
