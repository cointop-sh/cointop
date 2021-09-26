package chartplot

import (
	"math"

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

func (c *ChartPlot) GetChartDataSize(width int) int {
	axisYWidth := 30
	return (width * 2) - axisYWidth
}

// GetChartPoints ...
func (c *ChartPlot) GetChartPoints(width int) [][]rune {
	targetWidth := c.GetChartDataSize(width)
	if len(c.t.Data) != targetWidth {
		// Don't resample data if it's already the right size
		c.t.Data = interpolateData(c.t.Data, targetWidth)
	}
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

func interpolateData(data []float64, width int) []float64 {
	var res []float64
	if len(data) == 0 {
		return res
	}
	stepFactor := float64(len(data)-1) / float64(width-1)
	res = append(res, data[0])
	for i := 1; i < width-1; i++ {
		step := float64(i) * stepFactor
		before := math.Floor(step)
		after := math.Ceil(step)
		atPoint := step - before
		pointBefore := data[int(before)]
		pointAfter := data[int(after)]
		interpolated := pointBefore + (pointAfter-pointBefore)*atPoint
		res = append(res, interpolated)
	}
	res = append(res, data[len(data)-1])
	return res
}
