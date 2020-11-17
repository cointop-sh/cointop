package chartplot

import (
	"github.com/miguelmota/cointop/pkg/termui"
)

// ChartPlot ...
type ChartPlot struct {
	t *termui.LineChart
}

// NewChartPlot ...
func NewChartPlot() *ChartPlot {
	t := termui.NewLineChart()

	// NOTE: empty list means don't show x-axis labels
	t.DataLabels = []string{""}
	t.Border = false

	return &ChartPlot{
		t: t,
	}
}

// Height ...
func (c *ChartPlot) Height() int {
	return c.t.Height
}

// SetHeight ...
func (c *ChartPlot) SetHeight(height int) {
	c.t.Height = height
}

// Width ...
func (c *ChartPlot) Width() int {
	return c.t.Width
}

// SetWidth ...
func (c *ChartPlot) SetWidth(width int) {
	c.t.Width = width
}

// SetBorder ...
func (c *ChartPlot) SetBorder(enabled bool) {
	c.t.Border = enabled
}

// SetData ...
func (c *ChartPlot) SetData(data []float64) {
	// NOTE: edit `termui.LineChart.shortenFloatVal(float64)` to not
	// use exponential notation.
	c.t.Data = data
}

// GetChartPoints ...
func (c *ChartPlot) GetChartPoints(width int) [][]rune {
	termui.Body = termui.NewGrid()
	termui.Body.Width = width
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, c.t),
		),
	)

	var points [][]rune
	// calculate layout
	termui.Body.Align()
	w := termui.Body.Width
	h := c.Height()
	row := termui.Body.Rows[0]
	b := row.Buffer()
	for i := 0; i < h; i = i + 1 {
		var rowpoints []rune
		for j := 0; j < w; j = j + 1 {
			p := b.At(j, i)
			rowpoints = append(rowpoints, p.Ch)
		}
		points = append(points, rowpoints)
	}

	return points
}
