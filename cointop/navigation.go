package cointop

func (ct *Cointop) currentPage() int {
	ct.debuglog("currentPage()")
	return ct.State.page + 1
}

func (ct *Cointop) currentDisplayPage() int {
	ct.debuglog("currentDisplayPage()")
	return ct.State.page + 1
}

func (ct *Cointop) totalPages() int {
	ct.debuglog("totalPages()")
	return ct.getListCount() / ct.State.perPage
}

func (ct *Cointop) totalPagesDisplay() int {
	ct.debuglog("totalPagesDisplay()")
	return ct.totalPages() + 1
}

func (ct *Cointop) totalPerPage() int {
	return ct.State.perPage
}

func (ct *Cointop) setPage(page int) int {
	ct.debuglog("setPage()")
	if (page*ct.State.perPage) < ct.getListCount() && page >= 0 {
		ct.State.page = page
	}
	return ct.State.page
}

func (ct *Cointop) cursorDown() error {
	ct.debuglog("cursorDown()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the bottom
	if ct.isLastRow() {
		return nil
	}

	cx, cy := ct.Views.Table.Backing().Cursor()

	if err := ct.Views.Table.Backing().SetCursor(cx, cy+1); err != nil {
		ox, oy := ct.Views.Table.Backing().Origin()
		// set origin scrolls
		if err := ct.Views.Table.Backing().SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	ct.RowChanged()
	return nil
}

func (ct *Cointop) cursorUp() error {
	ct.debuglog("cursorUp()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the top
	if ct.isFirstRow() {
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
	ct.RowChanged()
	return nil
}

func (ct *Cointop) pageDown() error {
	ct.debuglog("pageDown()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the bottom
	if ct.isLastRow() {
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
	ct.RowChanged()
	return nil
}

func (ct *Cointop) pageUp() error {
	ct.debuglog("pageUp()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the top
	if ct.isFirstRow() {
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
	ct.RowChanged()
	return nil
}

func (ct *Cointop) navigateFirstLine() error {
	ct.debuglog("navigateFirstLine()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the top
	if ct.isFirstRow() {
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

	ct.RowChanged()
	return nil
}

func (ct *Cointop) navigateLastLine() error {
	ct.debuglog("navigateLastLine()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the bottom
	if ct.isLastRow() {
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

	ct.RowChanged()
	return nil
}

func (ct *Cointop) navigatePageFirstLine() error {
	ct.debuglog("navigatePageFirstLine()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the correct line
	if ct.isPageFirstLine() {
		return nil
	}

	cx, _ := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetCursor(cx, 0); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

func (ct *Cointop) navigatePageMiddleLine() error {
	ct.debuglog("navigatePageMiddleLine()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the correct line
	if ct.isPageMiddleLine() {
		return nil
	}

	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	if err := ct.Views.Table.Backing().SetCursor(cx, (sy/2)-1); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

func (ct *Cointop) navigatePageLastLine() error {
	ct.debuglog("navigatePageLastLine()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}

	// NOTE: return if already at the correct line
	if ct.isPageLastLine() {
		return nil
	}

	cx, _ := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	if err := ct.Views.Table.Backing().SetCursor(cx, sy-1); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

func (ct *Cointop) prevPage() error {
	ct.debuglog("prevPage()")

	// NOTE: return if already at the first page
	if ct.isFirstPage() {
		return nil
	}

	ct.setPage(ct.State.page - 1)
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

func (ct *Cointop) nextPage() error {
	ct.debuglog("nextPage()")

	// NOTE: return if already at the last page
	if ct.isLastPage() {
		return nil
	}

	ct.setPage(ct.State.page + 1)
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

func (ct *Cointop) nextPageTop() error {
	ct.debuglog("nextPageTop()")

	ct.nextPage()
	ct.navigateFirstLine()

	return nil
}

func (ct *Cointop) prevPageTop() error {
	ct.debuglog("prevtPageTop()")

	ct.prevPage()
	ct.navigateLastLine()

	return nil
}

func (ct *Cointop) firstPage() error {
	ct.debuglog("firstPage()")

	// NOTE: return if already at the first page
	if ct.isFirstPage() {
		return nil
	}

	ct.State.page = 0
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

func (ct *Cointop) isFirstRow() bool {
	ct.debuglog("isFirstRow()")

	_, y := ct.Views.Table.Backing().Origin()
	_, cy := ct.Views.Table.Backing().Cursor()

	return (cy + y) == 0
}

func (ct *Cointop) isLastRow() bool {
	ct.debuglog("isLastRow()")

	_, y := ct.Views.Table.Backing().Origin()
	_, cy := ct.Views.Table.Backing().Cursor()
	numRows := len(ct.State.coins) - 1

	return (cy + y + 1) > numRows
}

func (ct *Cointop) isFirstPage() bool {
	ct.debuglog("isFirstPage()")
	return ct.State.page == 0
}

func (ct *Cointop) isLastPage() bool {
	ct.debuglog("isLastPage()")
	return ct.State.page == ct.totalPages()-1
}

func (ct *Cointop) isPageFirstLine() bool {
	ct.debuglog("isPageFirstLine()")

	_, cy := ct.Views.Table.Backing().Cursor()
	return cy == 0
}

func (ct *Cointop) isPageMiddleLine() bool {
	ct.debuglog("isPageMiddleLine()")

	_, cy := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	return (sy/2)-1 == cy
}

func (ct *Cointop) isPageLastLine() bool {
	ct.debuglog("isPageLastLine()")

	_, cy := ct.Views.Table.Backing().Cursor()
	_, sy := ct.Views.Table.Backing().Size()
	return cy+1 == sy
}

func (ct *Cointop) lastPage() error {
	ct.debuglog("lastPage()")

	// NOTE: return if already at the last page
	if ct.isLastPage() {
		return nil
	}

	ct.State.page = ct.getListCount() / ct.State.perPage
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

func (ct *Cointop) goToPageRowIndex(idx int) error {
	ct.debuglog("goToPageRowIndex()")
	if ct.Views.Table.Backing() == nil {
		return nil
	}
	cx, _ := ct.Views.Table.Backing().Cursor()
	if err := ct.Views.Table.Backing().SetCursor(cx, idx); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

func (ct *Cointop) goToGlobalIndex(idx int) error {
	ct.debuglog("goToGlobalIndex()")
	perpage := ct.totalPerPage()
	atpage := idx / perpage
	ct.setPage(atpage)
	rowIndex := (idx % perpage)
	ct.highlightRow(rowIndex)
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) highlightRow(idx int) error {
	ct.debuglog("highlightRow()")
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

// CursorDownOrNextPage ...
func (ct *Cointop) CursorDownOrNextPage() error {
	ct.debuglog("CursorDownOrNextPage()")
	if ct.isLastRow() {
		if ct.isLastPage() {
			return nil
		}

		if err := ct.nextPageTop(); err != nil {
			return err
		}

		return nil
	}

	if err := ct.cursorDown(); err != nil {
		return err
	}

	return nil
}

// CursorUpOrPreviousPage ...
func (ct *Cointop) CursorUpOrPreviousPage() error {
	ct.debuglog("CursorUpOrPreviousPage()")
	if ct.isFirstRow() {
		if ct.isFirstPage() {
			return nil
		}

		if err := ct.prevPageTop(); err != nil {
			return err
		}

		return nil
	}

	if err := ct.cursorUp(); err != nil {
		return err
	}

	return nil
}
