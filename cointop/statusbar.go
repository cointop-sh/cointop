package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateStatusbar(s string) {
	ct.update(func() {
		ct.statusbarview.Clear()
		currpage := ct.currentDisplayPage()
		totalpages := ct.totalPages()
		base := "[q]Quit [?]Help [c]Chart [/]Search"
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("%v [← →]Page %v/%v %s", base, currpage, totalpages, s), ct.maxtablewidth, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.updateStatusbar(fmt.Sprintf("[o]Open %s", url))
}
