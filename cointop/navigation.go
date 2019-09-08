package cointop

func (ct *Cointop) currentPage() int {
	return ct.State.page + 1
}

func (ct *Cointop) currentDisplayPage() int {
	return ct.State.page + 1
}

func (ct *Cointop) totalPages() int {
	return ct.getListCount() / ct.State.perPage
}

func (ct *Cointop) totalPagesDisplay() int {
	return ct.totalPages() + 1
}

func (ct *Cointop) totalPerPage() int {
	return ct.State.perPage
}

func (ct *Cointop) setPage(page int) int {
	if (page*ct.State.perPage) < ct.getListCount() && page >= 0 {
		ct.State.page = page
	}
	return ct.State.page
}

func (ct *Cointop) cursorDown() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	_, y := ct.Views.Table.Backing().Origin()
	cx, cy := ct.Views.Table.Backing().Cursor()
	numRows := len(ct.State.coins) - 1
	if (cy + y + 1) > numRows {
		return nil
	}
	if err := ct.Views.Table.Backing().SetCursor(cx, cy+1); err != nil {
		ox, oy := ct.Views.Table.Backing().Origin()
		// set origin scrolls
		if err := ct.Views.Table.Backing().SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) cursorUp() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	ox, oy := ct.Views.Table.Backing().Origin()
	cx, cy := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetCursor(cx, cy-1); err != nil && oy > 0 {
		// set origin scrolls
		if err := ct.Views.Table.Backing().SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) pageDown() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	ox, oy := ct.Views.Table.Backing().Origin() // this is prev origin position
	cx, _ := ct.Views.Table.Backing().Cursor()  // relative cursor position
	_, sy := ct.Views.Table.Backing().Size()    // rows in visible view
	k := oy + sy
	l := len(ct.State.coins)
	// end of table
	if (oy + sy + sy) > l {
		k = l - sy
	}
	// select last row if next jump is out of bounds
	if k < 0 {
		k = 0
		sy = l
	}

	if err := ct.Views.Table.Backing().SetOrigin(ox, k); err != nil {
		return err
	}
	// move cursor to last line if can't scroll further
	if k == oy {
		if err := ct.Views.Table.Backing().SetCursor(cx, sy-1); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) pageUp() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	ox, oy := ct.Views.Table.Backing().Origin()
	cx, _ := ct.Views.Table.Backing().Cursor() // relative cursor position
	_, sy := ct.Views.Table.Backing().Size()   // rows in visible view
	k := oy - sy
	if k < 0 {
		k = 0
	}
	if err := ct.Views.Table.Backing().SetOrigin(ox, k); err != nil {
		return err
	}
	// move cursor to first line if can't scroll further
	if k == oy {
		if err := ct.Views.Table.Backing().SetCursor(cx, 0); err != nil {
			return err
		}
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) navigateFirstLine() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	ox, _ := ct.Views.Table.Backing().Origin()
	cx, _ := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetOrigin(ox, 0); err != nil {
		return err
	}
	if err := ct.Views.Table.Backing().SetCursor(cx, 0); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) navigateLastLine() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	ox, _ := ct.Views.Table.Backing().Origin()
	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	l := len(ct.State.coins)
	k := l - sy
	if err := ct.Views.Table.Backing().SetOrigin(ox, k); err != nil {
		return err
	}
	if err := ct.Views.Table.Backing().SetCursor(cx, sy-1); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) navigatePageFirstLine() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	cx, _ := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetCursor(cx, 0); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) navigatePageMiddleLine() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	if err := ct.Views.Table.Backing().SetCursor(cx, (sy/2)-1); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) navigatePageLastLine() error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	if err := ct.Views.Table.Backing().SetCursor(cx, sy-1); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) prevPage() error {
	if ct.isFirstPage() {
		return nil
	}
	ct.setPage(ct.State.page - 1)
	ct.updateTable()
	ct.rowChanged()
	return nil
}

func (ct *Cointop) nextPage() error {
	ct.setPage(ct.State.page + 1)
	ct.updateTable()
	ct.rowChanged()
	return nil
}

func (ct *Cointop) firstPage() error {
	ct.State.page = 0
	ct.updateTable()
	ct.rowChanged()
	return nil
}

func (ct *Cointop) isFirstPage() bool {
	return ct.State.page == 0
}

func (ct *Cointop) lastPage() error {
	ct.State.page = ct.getListCount() / ct.State.perPage
	ct.updateTable()
	ct.rowChanged()
	return nil
}

func (ct *Cointop) goToPageRowIndex(idx int) error {
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	cx, _ := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetCursor(cx, idx); err != nil {
		return err
	}
	ct.rowChanged()
	return nil
}

func (ct *Cointop) goToGlobalIndex(idx int) error {
	perpage := ct.totalPerPage()
	atpage := idx / perpage
	ct.setPage(atpage)
	rowIndex := (idx % perpage)
	ct.highlightRow(rowIndex)
	ct.updateTable()
	return nil
}

func (ct *Cointop) highlightRow(idx int) error {
	ct.Views.Table.Backing().SetOrigin(0, 0)
	ct.Views.Table.Backing().SetCursor(0, 0)
	ox, _ := ct.Views.Table.Backing().Origin()
	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	perpage := ct.totalPerPage()
	p := idx % perpage
	oy := (p / sy) * sy
	cy := p % sy
	if oy > 0 {
		ct.Views.Table.Backing().SetOrigin(ox, oy)
	}
	ct.Views.Table.Backing().SetCursor(cx, cy)

	return nil
}
