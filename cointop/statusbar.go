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
		base := fmt.Sprintf("%sQuit %sHelp %sChart %sRange %sSearch %sConvert", "[Q]", "[?]", "[Enter]", "[[ ]]", "[/]", "[C]")
		fmt.Fprintln(ct.statusbarview, pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxtablewidth, " "))
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLink()
	ct.updateStatusbar(fmt.Sprintf("%sOpen %s", "[O]", url))
}
