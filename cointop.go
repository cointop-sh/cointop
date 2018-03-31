package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/apis"
	apitypes "github.com/miguelmota/cointop/apis/types"
	"github.com/miguelmota/cointop/color"
	"github.com/miguelmota/cointop/pad"
	"github.com/miguelmota/cointop/table"
)

var (
	oneMinute int64 = 60
	oneHour         = oneMinute * 60
	oneDay          = oneHour * 24
	oneWeek         = oneDay * 7
	oneMonth        = oneDay * 30
	oneYear         = oneDay * 365
)

// Cointop cointop
type Cointop struct {
	g           *gocui.Gui
	chartview   *gocui.View
	chartpoints [][]termui.Cell
	headersview *gocui.View
	tableview   *gocui.View
	table       *table.Table
	statusview  *gocui.View
	sortdesc    bool
	currentsort string
	api         apis.Interface
	coins       []*apitypes.Coin
}

func (ct *Cointop) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	chartHeight := 10
	if v, err := g.SetView("chart", 0, 0, maxX, chartHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.chartview = v
		ct.chartview.Frame = false
		ct.updateChart()
	}

	if v, err := g.SetView("header", 0, chartHeight, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		t := table.New().SetWidth(maxX)
		headers := []string{
			pad.Right("[r]ank", 13, " "),
			pad.Right("[n]ame", 13, " "),
			pad.Right("[s]ymbol", 8, " "),
			pad.Left("[p]rice", 10, " "),
			pad.Left("[m]arket cap", 17, " "),
			pad.Left("24H [v]olume", 15, " "),
			pad.Left("[1]H%", 9, " "),
			pad.Left("[2]4H%", 9, " "),
			pad.Left("[7]D%", 9, " "),
			pad.Left("[t]otal supply", 20, " "),
			pad.Left("[a]vailable supply", 19, " "),
			pad.Left("[l]ast updated", 17, " "),
		}
		for _, h := range headers {
			t.AddCol(h)
		}

		t.Format().Fprint(v)
		ct.headersview = v
		ct.headersview.Frame = false
		ct.headersview.Highlight = true
		ct.headersview.SelBgColor = gocui.ColorGreen
		ct.headersview.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView("table", 0, chartHeight+1, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.tableview = v
		ct.tableview.Frame = false
		ct.tableview.Highlight = true
		ct.tableview.SelBgColor = gocui.ColorCyan
		ct.tableview.SelFgColor = gocui.ColorBlack
		ct.updateTable()
		ct.rowChanged()
	}

	if v, err := g.SetView("status", 0, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.statusview = v
		ct.statusview.Frame = false
		ct.statusview.Highlight = true
		ct.statusview.SelBgColor = gocui.ColorCyan
		ct.statusview.SelFgColor = gocui.ColorBlack
		ct.updateStatus("")
	}

	return nil
}

func (ct *Cointop) rowChanged() {
	ct.showLink()
}

func (ct *Cointop) fetchData() ([]*apitypes.Coin, error) {
	limit := 100
	result := []*apitypes.Coin{}
	coins, err := ct.api.GetAllCoinData(int(limit))
	if err != nil {
		return result, err
	}

	for i := range coins {
		coin := coins[i]
		result = append(result, &coin)
	}

	return result, nil
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

func (ct *Cointop) updateTable() error {
	maxX, _ := ct.g.Size()
	ct.table = table.New().SetWidth(maxX)
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.AddCol("")
	ct.table.HideColumHeaders = true
	var err error
	if len(ct.coins) == 0 {
		ct.coins, err = ct.fetchData()
		if err != nil {
			return err
		}
	}
	for _, coin := range ct.coins {
		unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
		lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
		colorprice := color.Cyan
		color1h := color.White
		color24h := color.White
		color7d := color.White
		if coin.PercentChange1H > 0 {
			color1h = color.Green
		}
		if coin.PercentChange1H < 0 {
			color1h = color.Red
		}
		if coin.PercentChange24H > 0 {
			color24h = color.Green
		}
		if coin.PercentChange24H < 0 {
			color24h = color.Red
		}
		if coin.PercentChange7D > 0 {
			color7d = color.Green
		}
		if coin.PercentChange7D < 0 {
			color7d = color.Red
		}
		ct.table.AddRow(
			pad.Left(fmt.Sprint(coin.Rank), 4, " "),
			pad.Right(coin.Name, 22, " "),
			pad.Right(coin.Symbol, 6, " "),
			colorprice(pad.Left(humanize.Commaf(coin.PriceUSD), 12, " ")),
			pad.Left(humanize.Commaf(coin.MarketCapUSD), 17, " "),
			pad.Left(humanize.Commaf(coin.USD24HVolume), 15, " "),
			color1h(pad.Left(fmt.Sprintf("%.2f%%", coin.PercentChange1H), 9, " ")),
			color24h(pad.Left(fmt.Sprintf("%.2f%%", coin.PercentChange24H), 9, " ")),
			color7d(pad.Left(fmt.Sprintf("%.2f%%", coin.PercentChange7D), 9, " ")),
			pad.Left(humanize.Commaf(coin.TotalSupply), 20, " "),
			pad.Left(humanize.Commaf(coin.AvailableSupply), 18, " "),
			pad.Left(fmt.Sprintf("%s", lastUpdated), 18, " "),
			// add %percent of cap
		)
	}

	ct.table.Format().Fprint(ct.tableview)
	return nil
}

func (ct *Cointop) updateStatus(s string) {
	maxX, _ := ct.g.Size()
	ct.statusview.Clear()
	fmt.Fprintln(ct.statusview, pad.Right(fmt.Sprintf("[q]uit %s", s), maxX, " "))
}

func (ct *Cointop) showLink() {
	url := ct.rowLink()
	ct.g.Update(func(g *gocui.Gui) error {
		ct.updateStatus(fmt.Sprintf("[enter]=%s", url))
		return nil
	})
}

func (ct *Cointop) enter(g *gocui.Gui, v *gocui.View) error {
	exec.Command("open", ct.rowLink()).Output()
	return nil
}

func (ct *Cointop) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ct *Cointop) selectedRowIndex() int {
	_, y := ct.tableview.Origin()
	_, cy := ct.tableview.Cursor()
	return y + cy
}

func (ct *Cointop) selectedCoin() *apitypes.Coin {
	return ct.coins[ct.selectedRowIndex()]
}

func (ct *Cointop) rowLink() string {
	slug := strings.ToLower(strings.Replace(ct.selectedCoin().Name, " ", "-", -1))
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s", slug)
}

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true
	g.Highlight = true

	ct := Cointop{
		g:   g,
		api: apis.NewCMC(),
	}
	g.SetManagerFunc(ct.layout)

	if err := ct.keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}
