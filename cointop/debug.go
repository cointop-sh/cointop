package cointop

import "log"

func (ct *Cointop) debuglog(s string) {
	if ct.debug {
		// TODO: do status debug bar
		log.Println(s)
	}
}
