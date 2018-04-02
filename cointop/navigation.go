package cointop

import (
	"github.com/jroimartin/gocui"
)

func (ct *Cointop) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	_, y := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	numRows := len(ct.coins) - 1
	//fmt.Fprint(v, cy)
	if (cy + y + 1) > numRows {
		return nil
	}
	if err := ct.tableview.SetCursor(cx, cy+1); err != nil {
		ox, oy := ct.tableview.Origin()
		if err := ct.tableview.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	ox, oy := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	//fmt.Fprint(v, oy)
	if err := ct.tableview.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := ct.tableview.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) pageDown(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	_, y := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	numRows := len(ct.coins) - 1
	_, sy := ct.tableview.Size()
	rows := sy
	if (cy + +y + rows) > numRows {
		// go to last row
		ct.tableview.SetCursor(cx, numRows)
		ox, _ := ct.tableview.Origin()
		ct.tableview.SetOrigin(ox, numRows)
		ct.rowChanged()
		return nil
	}
	if err := ct.tableview.SetCursor(cx, cy+rows); err != nil {
		ox, oy := ct.tableview.Origin()
		if err := ct.tableview.SetOrigin(ox, oy+rows); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) pageUp(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	ox, oy := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	_, sy := ct.tableview.Size()
	rows := sy
	if err := ct.tableview.SetCursor(cx, cy-rows); err != nil && oy > 0 {
		k := oy - rows
		if k < 0 {
			k = 0
		}
		if err := ct.tableview.SetOrigin(ox, k); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) prevPage(g *gocui.Gui, v *gocui.View) error {
	if (ct.page - 1) >= 0 {
		ct.page = ct.page - 1
	}
	ct.updateTable()
	ct.rowChanged()
	return nil
}

func (ct *Cointop) nextPage(g *gocui.Gui, v *gocui.View) error {
	if ((ct.page + 1) * ct.perpage) <= len(ct.allcoins) {
		ct.page = ct.page + 1
	}
	ct.updateTable()
	ct.rowChanged()
	return nil
}
