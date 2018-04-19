package open

import (
	"log"
	"os/exec"
	"strings"
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
	for i, cmd := range possibleCmds {
		out, err := exec.Command("command", "-v", cmd).Output()
		if err != nil {
			continue
		}

		bin := strings.TrimSpace(string(out))
		log.Fatal("bin ", bin)
		if bin != "" {
			openCmd = possibleCmds[i]
		}
	}
}

// URL open url
func URL(s string) error {
	if openCmd != "" {
		exec.Command(openCmd, s).Output()
	}
	return nil
}
