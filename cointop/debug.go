package cointop

import (
	"fmt"
	"os"

	"github.com/cointop-sh/cointop/pkg/pathutil"
	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) setLogConfiguration() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		ct.setLogOutputFile()
	}
}

func (ct *Cointop) setLogOutputFile() {
	filename := pathutil.NormalizePath(":PREFERRED_TEMP_DIR:/cointop.log")
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
