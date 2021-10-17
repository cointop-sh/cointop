package cointop

import (
	"os"

	"github.com/cointop-sh/cointop/pkg/pathutil"
	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) initlog() {
	filename := "/tmp/cointop.log"
	debugFile := os.Getenv("DEBUG_FILE")
	if debugFile != "" {
		filename = pathutil.NormalizePath(debugFile)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	ct.logfile = f
}
