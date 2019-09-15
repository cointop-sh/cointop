package cointop

import "log"

func (ct *Cointop) save() error {
	ct.debuglog("save()")
	ct.setSavingStatus()
	if err := ct.saveConfig(); err != nil {
		log.Fatal(err)
	}

	ct.cacheAllCoinsSlugMap()

	return nil
}

func (ct *Cointop) setSavingStatus() {
	ct.debuglog("setSavingStatus()")
	if ct.g == nil {
		return
	}

	go func() {
		ct.loadingTicks("saving", 590)
		ct.updateStatusbar("")
		ct.rowChanged()
	}()
}
