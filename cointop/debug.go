package cointop

import (
	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) debuglog(s string) {
	if ct.debug {
		// TODO: do status debug bar
		log.Println(s)
	}
}
