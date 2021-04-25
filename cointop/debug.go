package cointop

import (
	"log"
	"os"
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

// debuglog writes a debug log message to /tmp/cointop.log if the DEBUG environment is set.
func (ct *Cointop) debuglog(format string, args ...interface{}) {
	if !ct.debug {
		return
	}

	log.Printf(format+"\n", args...)
}
