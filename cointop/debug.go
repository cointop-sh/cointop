package cointop

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) initlog() {
	filename := "/tmp/cointop.log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	ct.logfile = f
}
