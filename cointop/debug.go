package cointop

import (
	"fmt"
	"os"

	"github.com/cointop-sh/cointop/pkg/pathutil"
	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) initlog() {
	filename := "/tmp/cointop.log"
	debugFile := os.Getenv("DEBUG_FILE")
	if debugFile != "" {
		filename = pathutil.NormalizePath(debugFile)
		if filename != debugFile && os.Getenv("DEBUG") != "" {
			fmt.Printf("Writing debug log to %s\n", filename)
		}
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	ct.logfile = f
}
