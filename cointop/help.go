package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) toggleHelp() error {
	ct.helpvisible = !ct.helpvisible
	if ct.helpvisible {
		return ct.showHelp()
	}
	return ct.hideHelp()
}

func (ct *Cointop) updateHelp() {
	str := fmt.Sprintf(" Help %s\n\n", pad.Left("[c]lose", ct.maxtablewidth-11, " "))
	i := 0
	for k, v := range ct.shortcutkeys {
		nl := " "
		i = i + 1
		if i%3 == 0 {
			i = 0
			nl = "\n"
		}
		str = fmt.Sprintf("%s%10s %-40s%s", str, fmt.Sprintf("[%s]", k), v, nl)
	}

	ct.update(func() {
		ct.helpview.Clear()
		ct.helpview.Frame = true
		fmt.Fprintln(ct.helpview, str)
	})
}

func (ct *Cointop) showHelp() error {
	ct.helpvisible = true
	ct.updateHelp()
	ct.setActiveView("help")
	return nil
}

func (ct *Cointop) hideHelp() error {
	ct.helpvisible = false
	ct.setViewOnBottom("help")
	ct.update(func() {
		ct.helpview.Clear()
		ct.helpview.Frame = false
		fmt.Fprintln(ct.helpview, "")
	})
	return nil
}
