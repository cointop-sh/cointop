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
		base := fmt.Sprintf("%sQuit %sHelp %sChart %sRange %sSearch %sConvert %sFavorites %sPortfolio %sSave", "[Q]", "[?]", "[Enter]", "[[ ]]", "[/]", "[C]", "[F]", "[P]", "[CTRL-S]")
		str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxtablewidth, " ")
		v := fmt.Sprintf("v%s", ct.version())
		str = str[:len(str)-len(v)+2] + v
		fmt.Fprintln(ct.statusbarview, str)
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLinkShort()
	ct.updateStatusbar(fmt.Sprintf("%sOpen %s", "[O]", url))
}
