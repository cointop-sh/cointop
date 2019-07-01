package cointop

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miguelmota/cointop/cointop/common/filecache"
	"github.com/miguelmota/cointop/cointop/common/gizak/termui"
	"github.com/miguelmota/cointop/cointop/common/timeutil"
)

var chartLock sync.Mutex
var chartPointsLock sync.Mutex

func (ct *Cointop) updateChart() error {
	if ct.Views.Chart.Backing == nil {
		return nil
	}

	chartLock.Lock()
	defer chartLock.Unlock()

	if ct.State.portfolioVisible {
		if err := ct.portfolioChart(); err != nil {
			return err
		}
	} else {
		symbol := ct.selectedCoinSymbol()
		name := ct.selectedCoinName()
		ct.calcChartPoints(symbol, name)
	}

	if len(ct.State.chartPoints) != 0 {
		ct.Views.Chart.Backing.Clear()
	}
	var body string
	if len(ct.State.chartPoints) == 0 {
		body = "\n\n\n\n\nnot enough data for chart"
	} else {
		for i := range ct.State.chartPoints {
			var s string
			for j := range ct.State.chartPoints[i] {
				p := ct.State.chartPoints[i][j]
				s = fmt.Sprintf("%s%c", s, p.Ch)
			}
			body = fmt.Sprintf("%s%s\n", body, s)

		}
	}
	ct.update(func() {
		if ct.Views.Chart.Backing == nil {
			return
		}

		fmt.Fprint(ct.Views.Chart.Backing, ct.colorscheme.Chart(body))
	})

	return nil
}

func (ct *Cointop) calcChartPoints(symbol string, name string) error {
	maxX := ct.maxTableWidth - 3
	chartPointsLock.Lock()
	defer chartPointsLock.Unlock()
	// TODO: not do this (SoC)
	go ct.updateMarketbar()

	chart := termui.NewLineChart()
	chart.Height = 10
	chart.Border = false

	// NOTE: empty list means don't show x-axis labels
	chart.DataLabels = []string{""}

	rangeseconds := ct.chartRangesMap[ct.State.selectedChartRange]
	if ct.State.selectedChartRange == "YTD" {
		ytd := time.Now().Unix() - int64(timeutil.BeginningOfYear().Unix())
		rangeseconds = time.Duration(ytd) * time.Second
	}

	now := time.Now()
	nowseconds := now.Unix()
	start := nowseconds - int64(rangeseconds.Seconds())
	end := nowseconds

	var data []float64

	keyname := symbol
	if keyname == "" {
		keyname = "globaldata"
	}
	cachekey := ct.cacheKey(fmt.Sprintf("%s_%s", keyname, strings.Replace(ct.State.selectedChartRange, " ", "", -1)))

	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		data, _ = cached.([]float64)
		ct.debuglog("soft cache hit")
	}

	if len(data) == 0 {
		if symbol == "" {
			graphData, err := ct.api.GetGlobalMarketGraphData(start, end)
			if err != nil {
				return nil
			}
			for i := range graphData.MarketCapByAvailableSupply {
				price := graphData.MarketCapByAvailableSupply[i][1]
				data = append(data, price/1E9)
			}
		} else {
			graphData, err := ct.api.GetCoinGraphData(symbol, name, start, end)
			if err != nil {
				return nil
			}

			// NOTE: edit `termui.LineChart.shortenFloatVal(float64)` to not
			// use exponential notation.
			for i := range graphData.PriceUSD {
				price := graphData.PriceUSD[i][1]
				data = append(data, price)
			}
		}

		ct.cache.Set(cachekey, data, 10*time.Second)
		go func() {
			filecache.Set(cachekey, data, 24*time.Hour)
		}()
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

	ct.State.chartPoints = points

	return nil
}

func (ct *Cointop) portfolioChart() error {
	maxX := ct.maxTableWidth - 3
	chartPointsLock.Lock()
	defer chartPointsLock.Unlock()
	// TODO: not do this (SoC)
	go ct.updateMarketbar()

	chart := termui.NewLineChart()
	chart.Height = 10
	chart.Border = false

	// NOTE: empty list means don't show x-axis labels
	chart.DataLabels = []string{""}

	rangeseconds := ct.chartRangesMap[ct.State.selectedChartRange]
	if ct.State.selectedChartRange == "YTD" {
		ytd := time.Now().Unix() - int64(timeutil.BeginningOfYear().Unix())
		rangeseconds = time.Duration(ytd) * time.Second
	}

	now := time.Now()
	nowseconds := now.Unix()
	start := nowseconds - int64(rangeseconds.Seconds())
	end := nowseconds

	var data []float64
	portfolio := ct.getPortfolioSlice()
	chartname := ct.selectedCoinName()
	for _, p := range portfolio {
		// filter by selected chart if selected
		if chartname != "" {
			if chartname != p.Name {
				continue
			}
		}

		if p.Holdings <= 0 {
			continue
		}

		var graphData []float64
		cachekey := strings.ToLower(fmt.Sprintf("%s_%s", p.Symbol, strings.Replace(ct.State.selectedChartRange, " ", "", -1)))
		cached, found := ct.cache.Get(cachekey)
		if found {
			// cache hit
			graphData, _ = cached.([]float64)
			ct.debuglog("soft cache hit")
		} else {
			filecache.Get(cachekey, &graphData)

			if len(graphData) == 0 {
				time.Sleep(2 * time.Second)
				apiGraphData, err := ct.api.GetCoinGraphData(p.Symbol, p.Name, start, end)
				if err != nil {
					return err
				}
				for i := range apiGraphData.PriceUSD {
					price := apiGraphData.PriceUSD[i][1]
					graphData = append(graphData, price)
				}
			}

			ct.cache.Set(cachekey, graphData, 10*time.Second)
			go func() {
				filecache.Set(cachekey, graphData, 24*time.Hour)
			}()
		}

		for i := range graphData {
			price := graphData[i]
			sum := p.Holdings * price
			if len(data)-1 >= i {
				data[i] += sum
			}
			data = append(data, sum)
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

	ct.State.chartPoints = points

	return nil
}

func (ct *Cointop) nextChartRange() error {
	sel := 0
	max := len(ct.chartRanges)
	for i, k := range ct.chartRanges {
		if k == ct.State.selectedChartRange {
			sel = i + 1
			break
		}
	}
	if sel > max-1 {
		sel = 0
	}

	ct.State.selectedChartRange = ct.chartRanges[sel]

	go ct.updateChart()
	return nil
}

func (ct *Cointop) prevChartRange() error {
	sel := 0
	for i, k := range ct.chartRanges {
		if k == ct.State.selectedChartRange {
			sel = i - 1
			break
		}
	}
	if sel < 0 {
		sel = len(ct.chartRanges) - 1
	}

	ct.State.selectedChartRange = ct.chartRanges[sel]
	go ct.updateChart()
	return nil
}

func (ct *Cointop) firstChartRange() error {
	ct.State.selectedChartRange = ct.chartRanges[0]
	go ct.updateChart()
	return nil
}

func (ct *Cointop) lastChartRange() error {
	ct.State.selectedChartRange = ct.chartRanges[len(ct.chartRanges)-1]
	go ct.updateChart()
	return nil
}

func (ct *Cointop) toggleCoinChart() error {
	highlightedcoin := ct.highlightedRowCoin()
	if ct.State.selectedCoin == highlightedcoin {
		ct.State.selectedCoin = nil
	} else {
		ct.State.selectedCoin = highlightedcoin
	}

	go ct.updateChart()
	go ct.updateMarketbar()
	return nil
}
