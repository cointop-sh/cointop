package cointop

import (
	"fmt"
	"sort"

	"github.com/miguelmota/cointop/pkg/pad"
	log "github.com/sirupsen/logrus"
)

// UpdateHelp updates the help views
func (ct *Cointop) UpdateHelp() {
	log.Debug("UpdateHelp()")
	keys := make([]string, 0, len(ct.State.shortcutKeys))
	for k := range ct.State.shortcutKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" Help %s\n\n", pad.Left("[q] close ", ct.Width()-9, " ")))
	cnt := 0
	h := ct.Views.Menu.Height()
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
		item := fmt.Sprintf("%10s %-45s", k, ct.colorscheme.MenuLabel(v))
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

	versionLine := fmt.Sprintf("cointop %s - (C) 2017-2021 Miguel Mota", ct.Version())
	licenseLine := "Released under the Apache 2.0 License."
	instructionsLine := "List of keyboard shortcuts"
	infoLine := "See git.io/cointop for more info.\n Press ESC to return."
	content := fmt.Sprintf("%s %s\n %s\n\n %s\n\n%s\n %s", header, versionLine, licenseLine, instructionsLine, body, infoLine)

	ct.UpdateUI(func() error {
		ct.Views.Menu.SetFrame(true)
		return ct.Views.Menu.Update(content)
	})
}

// ShowHelp shows the help view
func (ct *Cointop) ShowHelp() error {
	log.Debug("ShowHelp()")
	ct.State.helpVisible = true
	ct.UpdateHelp()
	ct.SetActiveView(ct.Views.Menu.Name())
	return nil
}

// HideHelp hides the help view
func (ct *Cointop) HideHelp() error {
	log.Debug("HideHelp()")
	ct.State.helpVisible = false
	ct.ui.SetViewOnBottom(ct.Views.Menu)
	ct.SetActiveView(ct.Views.Table.Name())
	ct.UpdateUI(func() error {
		ct.Views.Menu.SetFrame(false)
		return ct.Views.Menu.Update("")
	})
	return nil
}

// ToggleHelp toggles the help view
func (ct *Cointop) ToggleHelp() error {
	log.Debug("ToggleHelp()")
	ct.State.helpVisible = !ct.State.helpVisible
	if ct.State.helpVisible {
		return ct.ShowHelp()
	}
	return ct.HideHelp()
}
