package open

import (
	"os/exec"
)

var openCmd string
var possibleCmds = []string{
	"Start-Process", // windows
}

var possibleShells = []string{
	"powershell.exe",
	"explorer.exe",
}

var mainShell string

func init() {
	for _, sh := range possibleShells {
		shell, err := exec.LookPath(sh)
		if err != nil {
			continue
		}

		mainShell = shell
		break
	}

	for _, cmd := range possibleCmds {
		err := exec.Command(mainShell, "Get-Command", cmd).Run()
		if err != nil {
			continue
		}

		openCmd = cmd
		break
	}
}

// URL open url
func URL(s string) error {
	if openCmd != "" {
		return exec.Command(mainShell, openCmd, s).Run()
	}
	return nil
}

// CommandExists returns true if an 'open' command exists
func CommandExists() bool {
	return openCmd != ""
}
