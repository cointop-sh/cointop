package open

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var openCmd string
var possibleCmds = []string{
	"gvfs-open",          // gnome linux
	"gnome-open",         // gnome linux
	"kde-open",           // kde linux
	"exo-open",           // xfce linux
	"enlightenment_open", // enlightenment linux
	"xdg-open",           // generic linux
	"open",               // mac
	"start",              // windows
	"cygstart",           // windows
}

func init() {
	for _, cmd := range possibleCmds {
		out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s %s %s", "command", "-v", cmd)).Output()
		if err != nil {
			log.Println("err ", err)
			continue
		}

		log.Println("out ", string(out))
		bin := strings.TrimSpace(string(out))
		if bin != "" {
			openCmd = bin
			break
		}
	}
	log.Println("cmd ", openCmd)
	time.Sleep(2 * time.Second)
}

// URL open url
func URL(s string) error {
	if openCmd != "" {
		exec.Command(openCmd, s).Output()
	}
	return nil
}
