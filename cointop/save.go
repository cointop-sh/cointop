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
	go func() {
		ct.loadingTicks("saving", 590)
		ct.updateStatusbar("")
		ct.rowChanged()
	}()
}
