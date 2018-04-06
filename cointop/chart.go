package cointop

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/color"
)

func (ct *Cointop) updateChart() error {
	maxX := ct.Width()
	if maxX > ct.maxtablewidth {
		maxX = ct.maxtablewidth
	}
	coin := ct.selectedCoinName()
	ct.chartPoints(maxX, coin)
	for i := range ct.chartpoints {
		var s string
		for j := range ct.chartpoints[i] {
			p := ct.chartpoints[i][j]
			s = fmt.Sprintf("%s%c", s, p.Ch)
		}
		fmt.Fprintln(ct.chartview, color.White(s))
	}
	return nil
}

func (ct *Cointop) chartPoints(maxX int, coin string) error {
	chart := termui.NewLineChart()
	chart.Height = 10
	chart.AxesColor = termui.ColorWhite
	chart.LineColor = termui.ColorCyan
	chart.Border = false

	now := time.Now()
	secs := now.Unix()
	start := secs - oneWeek
	end := secs

	var data []float64
	if coin == "" {
		graphData, err := ct.api.GetGlobalMarketGraphData(start, end)
		if err != nil {
			return nil
		}
		for i := range graphData.MarketCapByAvailableSupply {
			data = append(data, graphData.MarketCapByAvailableSupply[i][1]/1E9)
		}
	} else {
		graphData, err := ct.api.GetCoinGraphData(coin, start, end)
		if err != nil {
			return nil
		}
		for i := range graphData.PriceUSD {
			data = append(data, graphData.PriceUSD[i][1])
		}
	}

	chart.Data = data
	termui.Body = termui.NewGrid()
	termui.Body.Width = maxX
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, chart),
		),
	)

	var points [][]termui.Cell
	// calculate layout
	termui.Body.Align()
	w := termui.Body.Width
	h := 10
	row := termui.Body.Rows[0]
	b := row.Buffer()
	for i := 0; i < h; i = i + 1 {
		var rowpoints []termui.Cell
		for j := 0; j < w; j = j + 1 {
			p := b.At(j, i)
			rowpoints = append(rowpoints, p)
		}
		points = append(points, rowpoints)
	}

	ct.chartpoints = points
	return nil
}

func (ct *Cointop) selectedCoinName() string {
	coin := ct.selectedcoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

func (ct *Cointop) toggleCoinChart(g *gocui.Gui, v *gocui.View) error {
	highlightedcoin := ct.highlightedRowCoin()
	if ct.selectedcoin == highlightedcoin {
		ct.selectedcoin = nil
	} else {
		ct.selectedcoin = highlightedcoin
	}
	ct.Update(func() {
		ct.chartview.Clear()
		ct.updateMarketbar()
		ct.updateChart()
	})
	return nil
}
