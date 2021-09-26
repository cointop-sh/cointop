package cointop

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/miguelmota/cointop/pkg/chartplot"
	"github.com/miguelmota/cointop/pkg/timedata"
	"github.com/miguelmota/cointop/pkg/timeutil"
	"github.com/miguelmota/cointop/pkg/ui"
	log "github.com/sirupsen/logrus"
)

// PriceData is the time-series data for a Coin used when building a Portfolio view for chart
type PriceData struct {
	coin *Coin
	data [][]float64
}

// ChartView is structure for chart view
type ChartView = ui.View

// NewChartView returns a new chart view
func NewChartView() *ChartView {
	var view *ChartView = ui.NewView("chart")
	return view
}

var chartLock sync.Mutex
var chartPointsLock sync.Mutex

// ChartRanges returns list of chart ranges available
func ChartRanges() []string {
	return []string{
		"24H",
		"3D",
		"7D",
		"1M",
		"3M",
		"6M",
		"1Y",
		"YTD",
		"All Time",
	}
}

// ChartRangesMap returns map of chart range time ranges
func ChartRangesMap() map[string]time.Duration {
	return map[string]time.Duration{
		"All Time": time.Duration(24 * 7 * 4 * 12 * 5 * time.Hour),
		"YTD":      time.Duration(1 * time.Second), // this will be calculated
		"1Y":       time.Duration(24 * 7 * 4 * 12 * time.Hour),
		"6M":       time.Duration(24 * 7 * 4 * 6 * time.Hour),
		"3M":       time.Duration(24 * 7 * 4 * 3 * time.Hour),
		"1M":       time.Duration(24 * 7 * 4 * time.Hour),
		"7D":       time.Duration(24 * 7 * time.Hour),
		"3D":       time.Duration(24 * 3 * time.Hour),
		"24H":      time.Duration(24 * time.Hour),
		"6H":       time.Duration(6 * time.Hour),
		"1H":       time.Duration(1 * time.Hour),
	}
}

// UpdateChart updates the chart view
func (ct *Cointop) UpdateChart() error {
	log.Debug("UpdateChart()")
	chartLock.Lock()
	defer chartLock.Unlock()

	if ct.IsPortfolioVisible() {
		if err := ct.PortfolioChart(); err != nil {
			return err
		}
	} else {
		symbol := ct.SelectedCoinSymbol()
		name := ct.SelectedCoinName()
		ct.ChartPoints(symbol, name)
	}

	var body string
	if len(ct.State.chartPoints) == 0 {
		body = "\n\n\n\n\nnot enough data for chart"
	} else {
		for i := range ct.State.chartPoints {
			var s string
			for j := range ct.State.chartPoints[i] {
				p := ct.State.chartPoints[i][j]
				s = fmt.Sprintf("%s%c", s, p)
			}
			body = fmt.Sprintf("%s%s\n", body, s)

		}
	}

	ct.UpdateUI(func() error {
		ct.Views.Chart.Clear()
		return ct.Views.Chart.Update(ct.colorscheme.Chart(body))
	})

	return nil
}

// ChartPoints calculates the the chart points
func (ct *Cointop) ChartPoints(symbol string, name string) error {
	log.Debug("ChartPoints()")
	maxX := ct.ChartWidth()

	chartPointsLock.Lock()
	defer chartPointsLock.Unlock()

	// TODO: not do this (SoC)
	go ct.UpdateMarketbar()

	chart := chartplot.NewChartPlot()
	chart.SetHeight(ct.State.chartHeight)

	rangeseconds := ct.chartRangesMap[ct.State.selectedChartRange]
	if ct.State.selectedChartRange == "YTD" {
		ytd := time.Now().Unix() - int64(timeutil.BeginningOfYear().Unix())
		rangeseconds = time.Duration(ytd) * time.Second
	}

	now := time.Now()
	nowseconds := now.Unix()
	start := nowseconds - int64(rangeseconds.Seconds())
	end := nowseconds

	var cacheData [][]float64

	keyname := symbol
	if keyname == "" {
		keyname = "globaldata"
	}
	cachekey := ct.CompositeCacheKey(keyname, name, ct.State.currencyConversion, ct.State.selectedChartRange)

	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		cacheData, _ = cached.([][]float64)
		log.Debug("ChartPoints() soft cache hit")
	}

	if len(cacheData) == 0 {
		if symbol == "" {
			convert := ct.State.currencyConversion
			graphData, err := ct.api.GetGlobalMarketGraphData(convert, start, end)
			if err != nil {
				return nil
			}
			cacheData = graphData.MarketCapByAvailableSupply
		} else {
			convert := ct.State.currencyConversion
			graphData, err := ct.api.GetCoinGraphData(convert, symbol, name, start, end)
			if err != nil {
				return nil
			}
			sorted := graphData.Price
			sort.Slice(sorted[:], func(i, j int) bool {
				return sorted[i][0] < sorted[j][0]
			})
			cacheData = sorted
		}

		ct.cache.Set(cachekey, cacheData, 10*time.Second)
		if ct.filecache != nil {
			go func() {
				ct.filecache.Set(cachekey, cacheData, 24*time.Hour)
			}()
		}
	}

	// Resample cachedata
	timeQuantum := timedata.CalculateTimeQuantum(cacheData)
	newStart := time.Unix(start, 0).Add(timeQuantum)
	newEnd := time.Unix(end, 0).Add(-timeQuantum)
	timeData := timedata.ResampleTimeSeriesData(cacheData, float64(newStart.UnixMilli()), float64(newEnd.UnixMilli()), chart.GetChartDataSize(maxX))
	labels := timedata.BuildTimeSeriesLabels(timeData)

	// Extract just the values from the data
	var data []float64
	for i := range timeData {
		value := timeData[i][1]
		if math.IsNaN(value) {
			value = 0.0
		}
		data = append(data, value)
	}

	chart.SetData(data)
	chart.SetDataLabels(labels)
	ct.State.chartPoints = chart.GetChartPoints(maxX)

	return nil
}

// PortfolioChart renders the portfolio chart
func (ct *Cointop) PortfolioChart() error {
	log.Debug("PortfolioChart()")
	maxX := ct.ChartWidth()
	chartPointsLock.Lock()
	defer chartPointsLock.Unlock()

	// TODO: not do this (SoC)
	go ct.UpdateMarketbar()

	chart := chartplot.NewChartPlot()
	chart.SetHeight(ct.State.chartHeight)

	convert := ct.State.currencyConversion            // cache here
	selectedChartRange := ct.State.selectedChartRange // cache here
	rangeseconds := ct.chartRangesMap[selectedChartRange]
	if selectedChartRange == "YTD" {
		ytd := time.Now().Unix() - int64(timeutil.BeginningOfYear().Unix())
		rangeseconds = time.Duration(ytd) * time.Second
	}

	now := time.Now()
	nowseconds := now.Unix()
	start := nowseconds - int64(rangeseconds.Seconds())
	end := nowseconds

	var allCacheData []PriceData
	portfolio := ct.GetPortfolioSlice()
	chartname := ct.SelectedCoinName()
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

		var cacheData [][]float64 // [][time,value]
		cachekey := ct.CompositeCacheKey(p.Symbol, p.Name, convert, selectedChartRange)
		cached, found := ct.cache.Get(cachekey)
		if found {
			// cache hit
			cacheData, _ = cached.([][]float64)
			log.Debug("PortfolioChart() soft cache hit")
		} else {
			if ct.filecache != nil {
				ct.filecache.Get(cachekey, &cacheData)
			}

			if len(cacheData) == 0 {
				time.Sleep(2 * time.Second)

				apiGraphData, err := ct.api.GetCoinGraphData(convert, p.Symbol, p.Name, start, end)
				if err != nil {
					return err
				}

				cacheData = apiGraphData.Price
				sort.Slice(cacheData[:], func(i, j int) bool {
					return cacheData[i][0] < cacheData[j][0]
				})
			}

			ct.cache.Set(cachekey, cacheData, 10*time.Second)
			if ct.filecache != nil {
				go func() {
					ct.filecache.Set(cachekey, cacheData, 24*time.Hour)
				}()
			}
		}

		allCacheData = append(allCacheData, PriceData{p, cacheData})
	}

	// Use the gap between price samples to adjust start/end in by one interval
	var timeQuantum time.Duration
	for _, cacheData := range allCacheData {
		timeQuantum = timedata.CalculateTimeQuantum(cacheData.data)
		if timeQuantum != 0 {
			break // use the first one
		}
	}
	newStart := time.Unix(start, 0).Add(timeQuantum)
	newEnd := time.Unix(end, 0).Add(-timeQuantum)

	// Resample and sum data
	var data []float64
	var labels []string
	for i, cacheData := range allCacheData {
		coinData := timedata.ResampleTimeSeriesData(cacheData.data, float64(newStart.UnixMilli()), float64(newEnd.UnixMilli()), chart.GetChartDataSize(maxX))
		if i == 0 {
			labels = timedata.BuildTimeSeriesLabels(coinData)
		}
		// sum (excluding NaN)
		for i := range coinData {
			price := coinData[i][1]
			if math.IsNaN(price) {
				price = 0.0
			}
			sum := cacheData.coin.Holdings * price
			if i < len(data) {
				data[i] += sum
			} else {
				data = append(data, sum)
			}
		}
	}

	// Scale Portfolio Balances to hide value
	if ct.State.hidePortfolioBalances {
		var lastPrice = data[len(data)-1]
		if lastPrice > 0.0 {
			for i, price := range data {
				data[i] = 100 * price / lastPrice
			}
		}
	}

	chart.SetData(data)
	chart.SetDataLabels(labels)
	ct.State.chartPoints = chart.GetChartPoints(maxX)

	return nil
}

// ShortenChart decreases the chart height by one row
func (ct *Cointop) ShortenChart() error {
	log.Debug("ShortenChart()")
	candidate := ct.State.chartHeight - 1
	if candidate < 5 {
		return nil
	}
	ct.State.chartHeight = candidate
	ct.State.lastChartHeight = ct.State.chartHeight

	go ct.UpdateChart()
	return nil
}

// EnlargeChart increases the chart height by one row
func (ct *Cointop) EnlargeChart() error {
	log.Debug("EnlargeChart()")
	candidate := ct.State.lastChartHeight + 1
	if candidate > 30 {
		return nil
	}
	ct.State.chartHeight = candidate
	ct.State.lastChartHeight = ct.State.chartHeight

	go ct.UpdateChart()
	return nil
}

// NextChartRange sets the chart to the next range option
func (ct *Cointop) NextChartRange() error {
	log.Debug("NextChartRange()")
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

	go ct.UpdateChart()
	return nil
}

// PrevChartRange sets the chart to the prevous range option
func (ct *Cointop) PrevChartRange() error {
	log.Debug("PrevChartRange()")
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
	go ct.UpdateChart()
	return nil
}

// FirstChartRange sets the chart to the first range option
func (ct *Cointop) FirstChartRange() error {
	log.Debug("FirstChartRange()")
	ct.State.selectedChartRange = ct.chartRanges[0]
	go ct.UpdateChart()
	return nil
}

// LastChartRange sets the chart to the last range option
func (ct *Cointop) LastChartRange() error {
	log.Debug("LastChartRange()")
	ct.State.selectedChartRange = ct.chartRanges[len(ct.chartRanges)-1]
	go ct.UpdateChart()
	return nil
}

// ToggleCoinChart toggles between the global chart and the coin chart
func (ct *Cointop) ToggleCoinChart() error {
	log.Debug("ToggleCoinChart()")
	highlightedcoin := ct.HighlightedRowCoin()
	if ct.State.selectedCoin == highlightedcoin {
		ct.State.selectedCoin = nil
	} else {
		ct.State.selectedCoin = highlightedcoin
	}

	go func() {
		// keep these two synchronous to avoid race conditions
		ct.ShowChartLoader()
		ct.UpdateChart()
	}()

	// TODO: not do this (SoC)
	go ct.UpdateMarketbar()

	return nil
}

// ShowChartLoader shows chart loading indicator
func (ct *Cointop) ShowChartLoader() error {
	log.Debug("ShowChartLoader()")
	ct.UpdateUI(func() error {
		content := "\n\nLoading..."
		return ct.Views.Chart.Update(ct.colorscheme.Chart(content))
	})

	return nil
}

// ChartWidth returns the width for chart
func (ct *Cointop) ChartWidth() int {
	log.Debug("ChartWidth()")
	w := ct.Width()
	max := 175
	if w > max {
		return max
	}

	return w
}

// ToggleChartFullscreen toggles the chart fullscreen mode
func (ct *Cointop) ToggleChartFullscreen() error {
	log.Debug("ToggleChartFullscreen()")
	ct.State.onlyChart = !ct.State.onlyChart
	ct.State.onlyTable = false
	if !ct.State.onlyChart {
		// NOTE: cached values are initial config settings.
		// If the only-chart config was set then toggle
		// all other initial hidden views.
		onlyChart, _ := ct.cache.Get("onlyChart")

		if onlyChart.(bool) {
			ct.State.hideMarketbar = false
			ct.State.hideChart = false
			ct.State.hideTable = false
			ct.State.hideStatusbar = false
		} else {
			// NOTE: cached values store initial hidden views preferences.
			hideMarketbar, _ := ct.cache.Get("hideMarketbar")
			ct.State.hideMarketbar = hideMarketbar.(bool)
			hideChart, _ := ct.cache.Get("hideChart")
			ct.State.hideChart = hideChart.(bool)
			hideTable, _ := ct.cache.Get("hideTable")
			ct.State.hideTable = hideTable.(bool)
			hideStatusbar, _ := ct.cache.Get("hideStatusbar")
			ct.State.hideStatusbar = hideStatusbar.(bool)
		}
	}

	go func() {
		ct.UpdateTable()
		ct.UpdateChart()
	}()

	return nil
}
