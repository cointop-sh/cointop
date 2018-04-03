package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateStatusbar(s string) {
	maxX := ct.Width()
	ct.Update(func() {
		ct.statusbarview.Clear()
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("[q]uit [← →]page %s", s), maxX, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.Update(func() {
		ct.updateStatusbar(fmt.Sprintf("[↵]%s", url))
	})
}
