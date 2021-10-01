package table

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
	"github.com/cointop-sh/cointop/pkg/pad"
	"github.com/cointop-sh/cointop/pkg/table/align"
)

// Table table
type Table struct {
	cols             Cols
	rows             Rows
	sort             []SortBy
	width            int
	HideColumHeaders bool
}

// NewTable new table
func NewTable() *Table {
	return &Table{}
}

// SetWidth set table width
func (t *Table) SetWidth(w int) *Table {
	t.width = w
	return t
}

// AddCol add column
func (t *Table) AddCol(n string) *Col {
	c := &Col{name: n}
	t.cols = append(t.cols, c)
	return c
}

// AddRow add row
func (t *Table) AddRow(v ...interface{}) *Row {
	r := &Row{table: t, values: v, strValues: make([]string, len(v))}
	t.rows = append(t.rows, r)
	return r
}

// AddRowCells add row using cells
func (t *Table) AddRowCells(cells ...*RowCell) *Row {
	t.SetNumCol(len(cells))
	v := make([]interface{}, len(cells))
	for i, item := range cells {
		v[i] = item.String()
	}
	return t.AddRow(v...)
}

// SetNumCol sets the number of columns
func (t *Table) SetNumCol(count int) {
	for i := 0; i < count; i++ {
		t.AddCol("")
	}
}

// SortAscFn sort ascending function
func (t *Table) SortAscFn(n string, fn SortFn) *Table {
	i := t.cols.Index(n)
	s := SortBy{index: i, order: SortAsc, sortFn: fn}
	t.sort = append(t.sort, s)
	return t
}

// SortAsc sort ascending
func (t *Table) SortAsc(n string) *Table {
	return t.SortAscFn(n, nil)
}

// SortDescFn sort descending function
func (t *Table) SortDescFn(n string, fn SortFn) *Table {
	i := t.cols.Index(n)
	s := SortBy{index: i, order: SortDesc, sortFn: fn}
	t.sort = append(t.sort, s)
	return t
}

// SortDesc sort descending
func (t *Table) SortDesc(n string) *Table {
	return t.SortDescFn(n, nil)
}

// Sort sort
func (t *Table) Sort() *Table {
	if len(t.sort) > 0 {
		sort.Sort(t.rows)
	}
	return t
}

func (t *Table) colWidth() int {
	width := 0
	for _, c := range t.cols {
		if c.hide {
			continue
		}

		width += c.width
	}
	return width
}

func (t *Table) normalizeColWidthPerc() {
	perc := 0
	for _, c := range t.cols {
		if c.hide {
			continue
		}

		perc += c.minWidthPerc
	}

	for _, c := range t.cols {
		if c.hide {
			continue
		}

		c.perc = float32(c.minWidthPerc) / float32(perc)
	}
}

// Format format table
func (t *Table) Format() *Table {
	for _, c := range t.cols {
		s := stripansi.Strip(c.name)
		c.width = utf8.RuneCountInString(s) + 1
		if c.minWidth > c.width {
			c.width = c.minWidth
		}
	}

	for _, r := range t.rows {
		for j, v := range r.values {
			c := t.cols[j]

			if c.hide {
				continue
			}

			if c.formatFn != nil {
				r.strValues[j] = fmt.Sprintf("%s", c.formatFn(v))
			} else if c.format != "" {
				r.strValues[j] = fmt.Sprintf(c.format, v)
			} else {
				r.strValues[j] = fmt.Sprintf("%v", v)
			}

			s := stripansi.Strip(r.strValues[j])
			runeCount := utf8.RuneCountInString(s)
			if runeCount > t.cols[j].width {
				t.cols[j].width = runeCount
			}
		}
	}

	//t.normalizeColWidthPerc()
	unused := t.width - t.colWidth()
	if unused <= 0 {
		return t
	}

	for _, c := range t.cols {
		if c.hide {
			continue
		}

		if c.perc > 0 {
			c.width += int(float32(unused) * c.perc)
		}
	}

	var i int
	for i = len(t.cols) - 1; i >= 0; i-- {
		if t.cols[i].hide {
			continue
		}

		break
	}

	if len(t.cols) > 0 {
		t.cols[i].width += t.width - t.colWidth()
	}

	return t
}

// Fprint write
func (t *Table) Fprint(w io.Writer) {
	if !t.HideColumHeaders {
		for _, c := range t.cols {
			if c.hide {
				continue
			}

			var s string
			switch c.align {
			case AlignLeft:
				s = align.AlignLeft(c.name+" ", c.width)
			case AlignRight:
				s = align.AlignRight(c.name+" ", c.width)
			case AlignCenter:
				s = align.AlignCenter(c.name+" ", c.width)
			}

			fmt.Fprintf(w, "%s", s)
		}
		fmt.Fprintf(w, "\n")

		for _, c := range t.cols {
			if c.hide {
				continue
			}

			fmt.Fprintf(w, strings.Repeat("─", c.width))
		}
		fmt.Fprintf(w, "\n")
	}

	for _, r := range t.rows {
		for i, v := range r.strValues {
			c := t.cols[i]

			if c.hide {
				continue
			}

			var s string
			switch c.align {
			case AlignLeft:
				s = align.AlignLeft(v, c.width)
			case AlignRight:
				s = align.AlignRight(v, c.width)
			case AlignCenter:
				s = align.AlignCenter(v, c.width)
			}

			fmt.Fprintf(w, "%s", s)
		}
		// fill in rest of row with empty spaces to highlight all of row
		fmt.Fprintf(w, strings.Repeat(" ", t.width)+"\n")
	}
}

// RowCount returns the number of rows
func (t *Table) RowCount() int {
	return len(t.rows)
}

// RowCell is a row cell struct
type RowCell struct {
	LeftMargin  int
	RightMargin int
	Width       int
	LeftAlign   bool
	Color       func(a ...interface{}) string
	Text        string
}

// String returns row cell as string
func (rc *RowCell) String() string {
	t := rc.Text
	if rc.LeftAlign {
		t = pad.Right(t, rc.Width, " ")
	} else {
		t = fmt.Sprintf("%"+fmt.Sprintf("%v", rc.Width)+"s", t)
	}
	t = strings.Repeat(" ", rc.LeftMargin) + t + strings.Repeat(" ", rc.RightMargin)
	return rc.Color(t)
}
