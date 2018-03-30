package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bradfitz/slice"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/table"
	cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/willf/pad"
)

var (
	white = color.New(color.FgWhite).SprintFunc()
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	cyan  = color.New(color.FgCyan).SprintFunc()
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
	sortdesc    bool
	currentsort string
	coins       []*cmc.Coin
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
	graphData, err := cmc.GetGlobalMarketGraphData(start, end)
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
		g: g,
	}
	g.SetManagerFunc(ct.layout)

	if err := ct.keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}

func (ct *Cointop) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ct.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, ct.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ct.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, ct.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, ct.pageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, ct.pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', gocui.ModNone, ct.sort("rank", false)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, ct.sort("name", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 's', gocui.ModNone, ct.sort("symbol", false)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gocui.ModNone, ct.sort("price", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'm', gocui.ModNone, ct.sort("marketcap", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'v', gocui.ModNone, ct.sort("24hvolume", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '1', gocui.ModNone, ct.sort("1hchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '2', gocui.ModNone, ct.sort("24hchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '7', gocui.ModNone, ct.sort("7dchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 't', gocui.ModNone, ct.sort("totalsupply", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'a', gocui.ModNone, ct.sort("availablesupply", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, ct.sort("lastupdated", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ct.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, ct.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, ct.quit); err != nil {
		return err
	}
	return nil
}

func (ct *Cointop) sort(sortby string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if ct.currentsort == sortby {
			ct.sortdesc = !ct.sortdesc
		} else {
			ct.currentsort = sortby
			ct.sortdesc = desc
		}
		slice.Sort(ct.coins[:], func(i, j int) bool {
			if ct.sortdesc {
				i, j = j, i
			}
			a := ct.coins[i]
			b := ct.coins[j]
			switch sortby {
			case "rank":
				return a.Rank < b.Rank
			case "name":
				return a.Name < b.Name
			case "symbol":
				return a.Symbol < b.Symbol
			case "price":
				return a.PriceUSD < b.PriceUSD
			case "marketcap":
				return a.MarketCapUSD < b.MarketCapUSD
			case "24hvolume":
				return a.USD24HVolume < b.USD24HVolume
			case "1hchange":
				return a.PercentChange1H < b.PercentChange1H
			case "24hchange":
				return a.PercentChange24H < b.PercentChange24H
			case "7dchange":
				return a.PercentChange7D < b.PercentChange7D
			case "totalsupply":
				return a.TotalSupply < b.TotalSupply
			case "availablesupply":
				return a.AvailableSupply < b.AvailableSupply
			case "lastupdated":
				return a.LastUpdated < b.LastUpdated
			default:
				return a.Rank < b.Rank
			}
		})
		g.Update(func(g *gocui.Gui) error {
			ct.tableview.Clear()
			ct.setTable()
			return nil
		})
		g.Update(func(g *gocui.Gui) error {
			ct.chartview.Clear()
			maxX, _ := g.Size()
			_, cy := ct.chartview.Cursor()
			coin := "ethereum"
			ct.chartPoints(maxX, coin)
			ct.setChart()
			fmt.Fprint(v, cy)
			return nil
		})

		return nil
	}
}

func (ct *Cointop) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	_, y := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	numRows := len(ct.coins) - 1
	//fmt.Fprint(v, cy)
	if (cy + y + 1) > numRows {
		return nil
	}
	if err := ct.tableview.SetCursor(cx, cy+1); err != nil {
		ox, oy := ct.tableview.Origin()
		if err := ct.tableview.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func (ct *Cointop) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	ox, oy := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	//fmt.Fprint(v, oy)
	if err := ct.tableview.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := ct.tableview.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (ct *Cointop) pageDown(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	_, y := ct.tableview.Origin()
	cx, cy := ct.tableview.Cursor()
	numRows := len(ct.coins) - 1
	_, sy := ct.tableview.Size()
	rows := sy
	if (cy + y + rows) > numRows {
		// go to last row
		ct.tableview.SetCursor(cx, numRows)
		ox, _ := ct.tableview.Origin()
		ct.tableview.SetOrigin(ox, numRows)
		return nil
	}
	if err := ct.tableview.SetCursor(cx, cy+rows); err != nil {
		ox, oy := ct.tableview.Origin()
		if err := ct.tableview.SetOrigin(ox, oy+rows); err != nil {
			return err
		}
	}
	return nil
}

func (ct *Cointop) pageUp(g *gocui.Gui, v *gocui.View) error {
	if ct.tableview == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, sy := v.Size()
	rows := sy
	//fmt.Fprint(v, oy)
	if err := v.SetCursor(cx, cy-rows); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-rows); err != nil {
			return err
		}
	}
	return nil
}

func fetchData() ([]*cmc.Coin, error) {
	limit := 100
	result := []*cmc.Coin{}
	coins, err := cmc.GetAllCoinData(int(limit))
	if err != nil {
		return result, err
	}

	for i := range coins {
		coin := coins[i]
		result = append(result, &coin)
	}

	return result, nil
}

func (ct *Cointop) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	chartHeight := 10
	if v, err := g.SetView("chart", 0, 0, maxX, chartHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		ct.chartview = v
		ct.setChart()
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
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
	}

	if v, err := g.SetView("table", 0, chartHeight+1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		ct.tableview = v
		ct.setTable()
	}

	return nil
}

func (ct *Cointop) setChart() error {
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

func (ct *Cointop) setTable() error {
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
		ct.coins, err = fetchData()
		if err != nil {
			return err
		}
	}
	for _, coin := range ct.coins {
		unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
		lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
		colorprice := cyan
		color1h := white
		color24h := white
		color7d := white
		if coin.PercentChange1H > 0 {
			color1h = green
		}
		if coin.PercentChange1H < 0 {
			color1h = red
		}
		if coin.PercentChange24H > 0 {
			color24h = green
		}
		if coin.PercentChange24H < 0 {
			color24h = red
		}
		if coin.PercentChange7D > 0 {
			color7d = green
		}
		if coin.PercentChange7D < 0 {
			color7d = red
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
