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
		base := "[q]uit [?]help [c]hart [/]search"
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("%v [← →]page %v/%v %s", base, currpage, totalpages, s), ct.maxtablewidth, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.update(func() {
		ct.updateStatusbar(fmt.Sprintf("[↵]%s", url))
	})
}
