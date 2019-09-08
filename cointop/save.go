package cointop

import "log"

func (ct *Cointop) save() error {
	ct.setSavingStatus()
	if err := ct.saveConfig(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ct *Cointop) setSavingStatus() {
	if ct.g == nil {
		return
	}

	go func() {
		ct.loadingTicks("saving", 590)
		ct.updateStatusbar("")
		ct.rowChanged()
	}()
}
