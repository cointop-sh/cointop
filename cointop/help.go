package cointop

import (
	"fmt"
	"sort"

	"github.com/miguelmota/cointop/cointop/common/pad"
)

// HelpView is structure for help view
type HelpView struct {
	*View
}

// NewHelpView returns a new help view
func NewHelpView() *HelpView {
	return &HelpView{NewView("help")}
}

func (ct *Cointop) updateHelp() {
	ct.debuglog("updateHelp()")
	keys := make([]string, 0, len(ct.State.shortcutKeys))
	for k := range ct.State.shortcutKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" Help %s\n\n", pad.Left("[q] close ", ct.maxTableWidth-10, " ")))
	cnt := 0
	h := ct.Views.Help.Height()
	percol := h - 6
	cols := make([][]string, percol)
	for i := range cols {
		cols[i] = make([]string, 20)
	}
	for _, k := range keys {
		v := ct.State.shortcutKeys[k]
		if cnt%percol == 0 {
			cnt = 0
		}
		item := fmt.Sprintf("%10s %-40s", k, ct.colorscheme.MenuLabel(v))
		cols[cnt] = append(cols[cnt], item)
		cnt = cnt + 1
	}

	var body string
	for i := 0; i < percol; i++ {
		var row string
		for j := 0; j < len(cols[i]); j++ {
			item := cols[i][j]
			row = fmt.Sprintf("%s%s", row, item)
		}
		body = fmt.Sprintf("%s%s\n", body, row)
	}
	body = fmt.Sprintf("%s\n", body)

	infoline := " List of keyboard shortcuts\n\n"
	versionline := pad.Left(fmt.Sprintf("v%s", ct.Version()), ct.maxTableWidth-5, " ")
	content := header + infoline + body + versionline

	ct.Update(func() error {
		if ct.Views.Help.Backing() == nil {
			return nil
		}

		ct.Views.Help.Backing().Clear()
		ct.Views.Help.Backing().Frame = true
		fmt.Fprintln(ct.Views.Help.Backing(), content)
		return nil
	})
}

func (ct *Cointop) showHelp() error {
	ct.debuglog("showHelp()")
	ct.State.helpVisible = true
	ct.updateHelp()
	ct.SetActiveView(ct.Views.Help.Name())
	return nil
}

func (ct *Cointop) hideHelp() error {
	ct.debuglog("hideHelp()")
	ct.State.helpVisible = false
	ct.SetViewOnBottom(ct.Views.Help.Name())
	ct.SetActiveView(ct.Views.Table.Name())
	ct.Update(func() error {
		if ct.Views.Help.Backing() == nil {
			return nil
		}

		ct.Views.Help.Backing().Clear()
		ct.Views.Help.Backing().Frame = false
		fmt.Fprintln(ct.Views.Help.Backing(), "")
		return nil
	})
	return nil
}

func (ct *Cointop) toggleHelp() error {
	ct.debuglog("toggleHelp()")
	ct.State.helpVisible = !ct.State.helpVisible
	if ct.State.helpVisible {
		return ct.showHelp()
	}
	return ct.hideHelp()
}
