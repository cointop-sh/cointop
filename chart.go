package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gizak/termui"
)

func (ct *Cointop) updateChart() error {
	maxX, _ := ct.g.Size()
	if len(ct.chartpoints) == 0 {
		ct.chartPoints(maxX, "bitcoin")
	}

	for i := range ct.chartpoints {
		var s string
		for j := range ct.chartpoints[i] {
			p := ct.chartpoints[i][j]
			s = fmt.Sprintf("%s%c", s, p.Ch)
		}
		fmt.Fprintln(ct.chartview, s)
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
	start := secs - oneDay
	end := secs

	_ = coin
	//graphData, err := cmc.GetCoinGraphData(coin, start, end)
	graphData, err := ct.api.GetGlobalMarketGraphData(start, end)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var data []float64
	/*
		for i := range graphData.PriceUSD {
			data = append(data, graphData.PriceUSD[i][1])
		}
	*/
	for i := range graphData.MarketCapByAvailableSupply {
		data = append(data, graphData.MarketCapByAvailableSupply[i][1])
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
