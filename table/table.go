package table

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/miguelmota/cointop/table/align"
)

// Table table
type Table struct {
	cols             Cols
	rows             Rows
	sort             []SortBy
	width            int
	HideColumHeaders bool
}

// New new table
func New() *Table {
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
		c.width = len(c.name) + 1
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
				r.strValues[j] = fmt.Sprintf("%s", c.formatFn(v)) + " "
			} else if c.format != "" {
				r.strValues[j] = fmt.Sprintf(c.format, v) + " "
			} else {
				r.strValues[j] = fmt.Sprintf("%v", v) + " "
			}

			if len(r.strValues[j]) > t.cols[j].width {
				t.cols[j].width = len(r.strValues[j])
			}
		}
	}

	t.normalizeColWidthPerc()

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
	t.cols[i].width += t.width - t.colWidth()

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
		fmt.Fprintf(w, "\n")
	}
}
