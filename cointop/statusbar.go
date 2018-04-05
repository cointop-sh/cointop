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
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("[q]uit [h]elp [← →]page %v/%v %s", currpage, totalpages, s), maxX, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.Update(func() {
		ct.updateStatusbar(fmt.Sprintf("[↵]%s", url))
	})
}
