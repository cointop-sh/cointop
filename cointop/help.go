package cointop

import (
	"fmt"
	"sort"

	"github.com/miguelmota/cointop/pkg/color"
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
	keys := make([]string, 0, len(ct.shortcutkeys))
	for k := range ct.shortcutkeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := color.GreenBg(fmt.Sprintf(" Help %s\n\n", pad.Left("[q] close ", ct.maxtablewidth-10, " ")))
	cnt := 0
	h := ct.viewHeight(ct.helpviewname)
	percol := h - 6
	cols := make([][]string, percol)
	for i := range cols {
		cols[i] = make([]string, 20)
	}
	for _, k := range keys {
		v := ct.shortcutkeys[k]
		if cnt%percol == 0 {
			cnt = 0
		}
		item := fmt.Sprintf("%10s %-40s", k, color.Yellow(v))
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
	versionline := pad.Left(fmt.Sprintf("v%s", ct.version()), ct.maxtablewidth-5, " ")
	content := header + infoline + body + versionline

	ct.update(func() {
		ct.helpview.Clear()
		ct.helpview.Frame = true
		fmt.Fprintln(ct.helpview, content)
	})
}

func (ct *Cointop) showHelp() error {
	ct.helpvisible = true
	ct.updateHelp()
	ct.setActiveView(ct.helpviewname)
	return nil
}

func (ct *Cointop) hideHelp() error {
	ct.helpvisible = false
	ct.setViewOnBottom(ct.helpviewname)
	ct.setActiveView(ct.tableviewname)
	ct.update(func() {
		ct.helpview.Clear()
		ct.helpview.Frame = false
		fmt.Fprintln(ct.helpview, "")
	})
	return nil
}
