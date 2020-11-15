//+build !windows

package open

import (
	"fmt"
	"os/exec"
	"strings"
)

var openCmd string
var possibleCmds = []string{
	"xdg-open",           // generic linux
	"gvfs-open",          // gnome linux
	"gnome-open",         // gnome linux
	"kde-open",           // kde linux
	"exo-open",           // xfce linux
	"enlightenment_open", // enlightenment linux
	"open",               // mac
	"start",              // windows
	"cygstart",           // windows
}

func init() {
	for _, cmd := range possibleCmds {
		out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s %s %s", "command", "-v", cmd)).Output()
		if err != nil {
			continue
		}

		bin := strings.TrimSpace(string(out))
		if bin != "" {
			openCmd = bin
			break
		}
	}
}

// URL open url
func URL(url string) error {
	if openCmd != "" {
		return exec.Command(openCmd, url).Run()
	}
	return nil
}

// CommandExists returns true if an 'open' command exists
func CommandExists() bool {
	return openCmd != ""
}
