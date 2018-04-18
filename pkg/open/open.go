package open

import (
	"os/exec"
)

var openCmd string
var possibleCmds = []string{
	"xdg-open", // linux
	"open",     // mac
	"start",    // windows?
	"cygstart", // windows?
}

func init() {
	for i, cmd := range possibleCmds {
		out, err := exec.Command("command", "-v", cmd).Output()
		if err != nil {
			continue
		}

		bin := string(out)
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
