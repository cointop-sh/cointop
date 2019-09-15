package cointop

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) debuglog(msg string) {
	if !ct.debug {
		return
	}

	filename := "/tmp/cointop.log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	text := fmt.Sprintf("%v %s\n", time.Now().Unix(), msg)
	if _, err = f.WriteString(text); err != nil {
		log.Fatal(err)
	}
}
