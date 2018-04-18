package cointop

func (ct *Cointop) save() error {
	ct.setSavingStatus()
	ct.saveConfig()
	return nil
}

func (ct *Cointop) setSavingStatus() {
	go func() {
		ct.loadingTicks("saving", 590)
		ct.updateStatusbar("")
		ct.rowChanged()
	}()
}
