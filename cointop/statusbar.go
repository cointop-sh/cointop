package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateStatusbar(s string) {
	maxX := ct.Width()
	ct.Update(func() {
		ct.statusbarview.Clear()
		currpage := ct.getCurrentPage()
		totalpages := ct.getTotalPages()
		base := "[q]uit [?]help [c]hart"
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("%v [← →]page %v/%v %s", base, currpage, totalpages, s), maxX, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.Update(func() {
		ct.updateStatusbar(fmt.Sprintf("[↵]%s", url))
	})
}
