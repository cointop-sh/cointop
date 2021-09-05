package cointop

import log "github.com/sirupsen/logrus"

// Save saves the cointop settings to the config file
func (ct *Cointop) Save() error {
	log.Debug("Save()")
	ct.SetSavingStatus()
	if err := ct.SaveConfig(); err != nil {
		return err
	}

	ct.CacheAllCoinsSlugMap()

	return nil
}

// SetSavingStatus sets the saving indicator in the statusbar
func (ct *Cointop) SetSavingStatus() {
	log.Debug("SetSavingStatus()")
	if ct.g == nil {
		return
	}

	go func() {
		ct.loadingTicks("saving", 590)
		ct.UpdateStatusbar("")
		ct.RowChanged()
	}()
}
