package cointop

// CurrentPage returns the current page
func (ct *Cointop) CurrentPage() int {
	ct.debuglog("currentPage()")
	return ct.State.page + 1
}

// CurrentDisplayPage returns the current page in human readable format
func (ct *Cointop) CurrentDisplayPage() int {
	ct.debuglog("currentDisplayPage()")
	return ct.State.page + 1
}

// TotalPages returns the number of total pages
func (ct *Cointop) TotalPages() int {
	ct.debuglog("totalPages()")
	return ct.GetListCount() / ct.State.perPage
}

// TotalPagesDisplay returns the number of total pages in human readable format
func (ct *Cointop) TotalPagesDisplay() int {
	ct.debuglog("totalPagesDisplay()")
	return ct.TotalPages() + 1
}

// TotalPerPage returns the number max rows per page
func (ct *Cointop) TotalPerPage() int {
	return ct.State.perPage
}

// SetPage navigates to the selected page
func (ct *Cointop) SetPage(page int) int {
	ct.debuglog("setPage()")
	if (page*ct.State.perPage) < ct.GetListCount() && page >= 0 {
		ct.State.page = page
	}
	return ct.State.page
}

// CursorDown moves the cursor one row down
func (ct *Cointop) CursorDown() error {
	ct.debuglog("cursorDown()")
	// NOTE: return if already at the bottom
	if ct.IsLastRow() {
		return nil
	}

	cx, cy := ct.Views.Table.Cursor()
	if err := ct.Views.Table.SetCursor(cx, cy+1); err != nil {
		ox, oy := ct.Views.Table.Origin()
		// set origin scrolls
		if err := ct.Views.Table.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	ct.RowChanged()
	return nil
}

// CursorUp moves the cursor one row up
func (ct *Cointop) CursorUp() error {
	ct.debuglog("cursorUp()")
	// NOTE: return if already at the top
	if ct.IsFirstRow() {
		return nil
	}

	ox, oy := ct.Views.Table.Origin()
	cx, cy := ct.Views.Table.Cursor()

	if err := ct.Views.Table.SetCursor(cx, cy-1); err != nil && oy > 0 {
		// set origin scrolls
		if err := ct.Views.Table.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	ct.RowChanged()
	return nil
}

// PageDown moves the cursor one page down
func (ct *Cointop) PageDown() error {
	ct.debuglog("pageDown()")
	// NOTE: return if already at the bottom
	if ct.IsLastRow() {
		return nil
	}

	ox, oy := ct.Views.Table.Origin() // this is prev origin position
	cx := ct.Views.Table.CursorX()    // relative cursor position
	sy := ct.Views.Table.Height()     // rows in visible view
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

	if err := ct.Views.Table.SetOrigin(ox, k); err != nil {
		return err
	}

	// move cursor to last line if can't scroll further
	if k == oy {
		if err := ct.Views.Table.SetCursor(cx, sy-1); err != nil {
			return err
		}
	}
	ct.RowChanged()
	return nil
}

// PageUp moves the cursor one page up
func (ct *Cointop) PageUp() error {
	ct.debuglog("pageUp()")
	// NOTE: return if already at the top
	if ct.IsFirstRow() {
		return nil
	}

	ox, oy := ct.Views.Table.Origin()
	cx := ct.Views.Table.CursorX() // relative cursor position
	sy := ct.Views.Table.Height()  // rows in visible view
	k := oy - sy
	if k < 0 {
		k = 0
	}
	if err := ct.Views.Table.SetOrigin(ox, k); err != nil {
		return err
	}
	// move cursor to first line if can't scroll further
	if k == oy {
		if err := ct.Views.Table.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	ct.RowChanged()
	return nil
}

// NavigateFirstLine moves the cursor to the first row of the table
func (ct *Cointop) NavigateFirstLine() error {
	ct.debuglog("navigateFirstLine()")
	// NOTE: return if already at the top
	if ct.IsFirstRow() {
		return nil
	}

	ox := ct.Views.Table.OriginX()
	cx := ct.Views.Table.CursorX()
	if err := ct.Views.Table.SetOrigin(ox, 0); err != nil {
		return err
	}
	if err := ct.Views.Table.SetCursor(cx, 0); err != nil {
		return err
	}

	ct.RowChanged()
	return nil
}

// NavigateLastLine moves the cursor to the last row of the table
func (ct *Cointop) NavigateLastLine() error {
	ct.debuglog("navigateLastLine()")
	// NOTE: return if already at the bottom
	if ct.IsLastRow() {
		return nil
	}

	ox := ct.Views.Table.OriginX()
	cx := ct.Views.Table.CursorX()
	sy := ct.Views.Table.Height()
	l := len(ct.State.coins)
	k := l - sy
	if err := ct.Views.Table.SetOrigin(ox, k); err != nil {
		return err
	}
	if err := ct.Views.Table.SetCursor(cx, sy-1); err != nil {
		return err
	}

	ct.RowChanged()
	return nil
}

// NavigatePageFirstLine moves the cursor to the visible first row of the table
func (ct *Cointop) NavigatePageFirstLine() error {
	ct.debuglog("navigatePageFirstLine()")
	// NOTE: return if already at the correct line
	if ct.IsPageFirstLine() {
		return nil
	}

	cx := ct.Views.Table.CursorX()
	if err := ct.Views.Table.SetCursor(cx, 0); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

// NavigatePageMiddleLine moves the cursor to the visible middle row of the table
func (ct *Cointop) NavigatePageMiddleLine() error {
	ct.debuglog("navigatePageMiddleLine()")
	// NOTE: return if already at the correct line
	if ct.IsPageMiddleLine() {
		return nil
	}

	cx := ct.Views.Table.CursorX()
	sy := ct.Views.Table.Height()
	if err := ct.Views.Table.SetCursor(cx, (sy/2)-1); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

// NavigatePageLastLine moves the cursor to the visible last row of the table
func (ct *Cointop) navigatePageLastLine() error {
	ct.debuglog("navigatePageLastLine()")
	// NOTE: return if already at the correct line
	if ct.IsPageLastLine() {
		return nil
	}

	cx, _ := ct.Views.Table.Cursor()
	sy := ct.Views.Table.Height()
	if err := ct.Views.Table.SetCursor(cx, sy-1); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

// NextPage navigates to the next page
func (ct *Cointop) NextPage() error {
	ct.debuglog("nextPage()")

	// NOTE: return if already at the last page
	if ct.IsLastPage() {
		return nil
	}

	ct.SetPage(ct.State.page + 1)
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

// PrevPage navigates to the previous page
func (ct *Cointop) PrevPage() error {
	ct.debuglog("prevPage()")

	// NOTE: return if already at the first page
	if ct.IsFirstPage() {
		return nil
	}

	ct.SetPage(ct.State.page - 1)
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

// NextPageTop navigates to the first row of the next page
func (ct *Cointop) nextPageTop() error {
	ct.debuglog("nextPageTop()")

	ct.NextPage()
	ct.NavigateFirstLine()

	return nil
}

// PrevPageTop navigates to the first row of the previous page
func (ct *Cointop) PrevPageTop() error {
	ct.debuglog("prevtPageTop()")

	ct.PrevPage()
	ct.NavigateLastLine()

	return nil
}

// FirstPage navigates to the first page
func (ct *Cointop) FirstPage() error {
	ct.debuglog("firstPage()")

	// NOTE: return if already at the first page
	if ct.IsFirstPage() {
		return nil
	}

	ct.State.page = 0
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

// LastPage navigates to the last page
func (ct *Cointop) LastPage() error {
	ct.debuglog("lastPage()")

	// NOTE: return if already at the last page
	if ct.IsLastPage() {
		return nil
	}

	ct.State.page = ct.GetListCount() / ct.State.perPage
	ct.UpdateTable()
	ct.RowChanged()
	return nil
}

// IsFirstRow returns true if cursor is on first row
func (ct *Cointop) IsFirstRow() bool {
	ct.debuglog("isFirstRow()")
	oy := ct.Views.Table.OriginY()
	cy := ct.Views.Table.CursorY()

	return (cy + oy) == 0
}

// IsLastRow returns true if cursor is on last row
func (ct *Cointop) IsLastRow() bool {
	ct.debuglog("isLastRow()")
	oy := ct.Views.Table.OriginY()
	cy := ct.Views.Table.CursorY()
	numRows := len(ct.State.coins) - 1

	return (cy + oy + 1) > numRows
}

// IsFirstPage returns true if cursor is on the first page
func (ct *Cointop) IsFirstPage() bool {
	ct.debuglog("isFirstPage()")
	return ct.State.page == 0
}

// IsLastPage returns true if cursor is on the last page
func (ct *Cointop) IsLastPage() bool {
	ct.debuglog("isLastPage()")
	return ct.State.page == ct.TotalPages()-1
}

// IsPageFirstLine returns true if the cursor is on the visible first row
func (ct *Cointop) IsPageFirstLine() bool {
	ct.debuglog("isPageFirstLine()")

	cy := ct.Views.Table.CursorY()
	return cy == 0
}

// IsPageMiddleLine returns true if the cursor is on the visible middle row
func (ct *Cointop) IsPageMiddleLine() bool {
	ct.debuglog("isPageMiddleLine()")
	cy := ct.Views.Table.CursorY()
	sy := ct.Views.Table.Height()
	return (sy/2)-1 == cy
}

// IsPageLastLine returns true if the cursor is on the visible last row
func (ct *Cointop) IsPageLastLine() bool {
	ct.debuglog("isPageLastLine()")

	cy := ct.Views.Table.CursorY()
	sy := ct.Views.Table.Height()
	return cy+1 == sy
}

// GoToPageRowIndex navigates to the selected row index of the page
func (ct *Cointop) GoToPageRowIndex(idx int) error {
	ct.debuglog("goToPageRowIndex()")
	cx := ct.Views.Table.CursorX()
	if err := ct.Views.Table.SetCursor(cx, idx); err != nil {
		return err
	}
	ct.RowChanged()
	return nil
}

// GoToGlobalIndex navigates to the selected row index of all page rows
func (ct *Cointop) GoToGlobalIndex(idx int) error {
	ct.debuglog("goToGlobalIndex()")
	perpage := ct.TotalPerPage()
	atpage := idx / perpage
	ct.SetPage(atpage)
	rowIndex := (idx % perpage)
	ct.HighlightRow(rowIndex)
	ct.UpdateTable()
	return nil
}

// HighlightRow highlights the row at index
func (ct *Cointop) HighlightRow(idx int) error {
	ct.debuglog("highlightRow()")
	ct.Views.Table.SetOrigin(0, 0)
	ct.Views.Table.SetCursor(0, 0)
	ox := ct.Views.Table.OriginX()
	cx := ct.Views.Table.CursorX()
	sy := ct.Views.Table.Height()
	perpage := ct.TotalPerPage()
	p := idx % perpage
	oy := (p / sy) * sy
	cy := p % sy
	if oy > 0 {
		ct.Views.Table.SetOrigin(ox, oy)
	}
	ct.Views.Table.SetCursor(cx, cy)

	return nil
}

// CursorDownOrNextPage moves the cursor down one row or goes to the next page if cursor is on the last row
func (ct *Cointop) CursorDownOrNextPage() error {
	ct.debuglog("CursorDownOrNextPage()")
	if ct.IsLastRow() {
		if ct.IsLastPage() {
			return nil
		}

		if err := ct.nextPageTop(); err != nil {
			return err
		}

		return nil
	}

	if err := ct.CursorDown(); err != nil {
		return err
	}

	return nil
}

// CursorUpOrPreviousPage moves the cursor up one row or goes to the previous page if cursor is on the first row
func (ct *Cointop) CursorUpOrPreviousPage() error {
	ct.debuglog("CursorUpOrPreviousPage()")
	if ct.IsFirstRow() {
		if ct.IsFirstPage() {
			return nil
		}

		if err := ct.PrevPageTop(); err != nil {
			return err
		}

		return nil
	}

	if err := ct.CursorUp(); err != nil {
		return err
	}

	return nil
}
