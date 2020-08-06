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

// UpdateHelp updates the help views
func (ct *Cointop) UpdateHelp() {
	ct.debuglog("updateHelp()")
	keys := make([]string, 0, len(ct.State.shortcutKeys))
	for k := range ct.State.shortcutKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" Help %s\n\n", pad.Left("[q] close ", ct.maxTableWidth-10, " ")))
	cnt := 0
	h := ct.Views.Help.Height()
	percol := h - 11
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

	versionLine := fmt.Sprintf("cointop %s - (C) 2017-2020 Miguel Mota", ct.Version())
	licenseLine := "Released under the Apache 2.0 License."
	instructionsLine := "List of keyboard shortcuts"
	infoLine := "See git.io/cointop for more info.\n Press ESC to return."
	content := fmt.Sprintf("%s %s\n %s\n\n %s\n\n%s\n %s", header, versionLine, licenseLine, instructionsLine, body, infoLine)

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

// ShowHelp shows the help view
func (ct *Cointop) ShowHelp() error {
	ct.debuglog("showHelp()")
	ct.State.helpVisible = true
	ct.UpdateHelp()
	ct.SetActiveView(ct.Views.Help.Name())
	return nil
}

// HideHelp hides the help view
func (ct *Cointop) HideHelp() error {
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

// ToggleHelp toggles the help view
func (ct *Cointop) ToggleHelp() error {
	ct.debuglog("toggleHelp()")
	ct.State.helpVisible = !ct.State.helpVisible
	if ct.State.helpVisible {
		return ct.ShowHelp()
	}
	return ct.HideHelp()
}
