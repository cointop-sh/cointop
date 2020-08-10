package asciitable

import (
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Input ...
type Input struct {
	Data      [][]string
	Headers   []string
	Alignment []int
}

// AsciiTable ...
type AsciiTable struct {
	table       *tablewriter.Table
	tableString *strings.Builder
}

// NewAsciiTable ...
func NewAsciiTable(input *Input) *AsciiTable {
	tableString := &strings.Builder{}
	alignment := make([]int, len(input.Alignment))
	for i, value := range input.Alignment {
		switch value {
		case -1:
			alignment[i] = 3
		case 0:
			alignment[i] = 1
		case 1:
			alignment[i] = 2
		}
	}

	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)
	table.SetColumnAlignment(alignment)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.SetHeader(input.Headers)
	table.AppendBulk(input.Data)

	return &AsciiTable{
		table:       table,
		tableString: tableString,
	}
}

// String ...
func (t *AsciiTable) String() string {
	t.table.Render()
	return t.tableString.String()
}
