package main

import (
	"fmt"
	"log"
	"math"
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

var white = color.New(color.FgWhite).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func fmtTime(v interface{}) string {
	t := v.(time.Time)
	return t.Format("2006-01-02 15:04:05")
}

func ltTime(a interface{}, b interface{}) bool {
	return a.(time.Time).Before(b.(time.Time))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func chartCells(maxX int, coin string) [][]termui.Cell {
	//_ = termui.Init()
	//defer termui.Close()

	// quit on Ctrl-c
	/*
		termui.Handle("/sys/kbd/C-c", func(termui.Event) {
			termui.StopLoop()
		})
	*/
	chart := termui.NewLineChart()
	//chart.DataLabels = []string{""}
	chart.Height = 10
	chart.AxesColor = termui.ColorWhite
	chart.LineColor = termui.ColorCyan //| termui.AttrBold
	chart.Border = false

	var (
		oneMinute int64 = 60
		oneHour         = oneMinute * 60
		oneDay          = oneHour * 24
		//oneWeek         = oneDay * 7
		//oneMonth        = oneDay * 30
		//oneYear         = oneDay * 365
	)

	now := time.Now()
	secs := now.Unix()
	start := secs - oneDay
	end := secs

	_ = coin
	//graphData, err := cmc.GetCoinGraphData(coin, start, end)
	graphData, err := cmc.GetGlobalMarketGraphData(start, end)
	if err != nil {
		log.Fatal(err)
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
	//h := chart.InnerHeight()
	//w := chart.InnerWidth()

	// add grid rows and columns
	termui.Body = termui.NewGrid()
	//termui.Body.Rows = termui.Body.Rows[:0]
	termui.Body.Width = maxX
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, chart),
		),
	)

	var points [][]termui.Cell

	// calculate layout
	termui.Body.Align()
	// render to terminal
	//termui.Render(termui.Body)
	w := termui.Body.Width
	h := 10
	row := termui.Body.Rows[0]
	b := row.Buffer()
	for i := 0; i < h; i = i + 1 {
		var rowpoints []termui.Cell
		for j := 0; j < w; j = j + 1 {
			_ = b
			p := b.At(j, i)
			c := p
			rowpoints = append(rowpoints, c)
		}
		points = append(points, rowpoints)
	}
	//termui.Loop()

	return points
}

var points [][]termui.Cell

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.Highlight = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, pageDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, pageUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'r', gocui.ModNone, sort("rank", false)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, sort("name", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 's', gocui.ModNone, sort("symbol", false)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gocui.ModNone, sort("price", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'm', gocui.ModNone, sort("marketcap", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'v', gocui.ModNone, sort("24hvolume", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '1', gocui.ModNone, sort("1hchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '2', gocui.ModNone, sort("24hchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '7', gocui.ModNone, sort("7dchange", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 't', gocui.ModNone, sort("totalsupply", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'a', gocui.ModNone, sort("availablesupply", true)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, sort("lastupdated", true)); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

type strategy func(g *gocui.Gui, v *gocui.View) error

var sortDesc bool
var currentsort string

func sort(sortBy string, desc bool) strategy {
	return func(g *gocui.Gui, v *gocui.View) error {

		if currentsort == sortBy {
			sortDesc = !sortDesc
		} else {
			currentsort = sortBy
			sortDesc = desc
		}

		slice.Sort(coins[:], func(i, j int) bool {
			if sortDesc == true {
				i, j = j, i
			}
			a := coins[i]
			b := coins[j]
			switch sortBy {
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
			v, _ = g.View("table")
			v.Clear()
			setTable(g)
			return nil
		})

		g.Update(func(g *gocui.Gui) error {
			v, _ = g.View("chart")
			v.Clear()

			maxX, _ := g.Size()
			_, cy := v.Cursor()
			coin := "ethereum"
			points = chartCells(maxX, coin)
			setChart(g, v)
			fmt.Fprint(v, cy)
			return nil
		})

		return nil
	}
}

var coins []*cmc.Coin

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.View("table")
	if v != nil {
		_, y := v.Origin()
		cx, cy := v.Cursor()
		numRows := len(coins) - 1
		//fmt.Fprint(v, cy)
		if (cy + y + 1) > numRows {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.View("table")
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		//fmt.Fprint(v, oy)
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func pageDown(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.View("table")
	if v != nil {
		_, y := v.Origin()
		cx, cy := v.Cursor()
		numRows := len(coins) - 1
		_, sy := v.Size()
		rows := sy
		if (cy + y + rows) > numRows {
			// go to last row
			v.SetCursor(cx, numRows)
			ox, _ := v.Origin()
			v.SetOrigin(ox, numRows)
			return nil
		}
		if err := v.SetCursor(cx, cy+rows); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+rows); err != nil {
				return err
			}
		}
	}
	return nil
}

func pageUp(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.View("table")
	if v != nil {
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

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	chartHeight := 10
	if v, err := g.SetView("chart", 0, 0, maxX, chartHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//v.Frame = false
		setChart(g, v)
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
		_ = v
		_ = maxY
		setTable(g)
	}
	return nil
}

func setChart(g *gocui.Gui, v *gocui.View) error {
	v.Frame = false
	/*
		tm.Clear()
		tm.MoveCursor(0, 0)
		chart := tm.NewLineChart(maxX/2, 9)
		//chart.Flags = tm.DRAW_INDEPENDENT
		chart.Flags = tm.DRAW_RELATIVE
		data := new(tm.DataTable)
		data.AddColumn("")
		data.AddColumn("")

		var (
			oneMinute int64 = 60
			oneHour         = oneMinute * 60
			oneDay          = oneHour * 24
			//oneWeek         = oneDay * 7
			oneMonth = oneDay * 30
			//oneYear         = oneDay * 365
		)

		now := time.Now()
		secs := now.Unix()
		//start := secs - oneDay
		start := secs - oneMonth
		end := secs

		coin := "ethereum"
		graphData, err := cmc.GetCoinGraphData(coin, start, end)
		if err != nil {
			log.Fatal(err)
		}

		var dt []float64
		for i := range graphData.PriceUSD {
			dt = append(dt, graphData.PriceUSD[i][1])
		}

		for i, d := range dt {
			_ = d
			_ = i
			data.AddRow(float64(i), d)
		}
		_ = dt

		out := chart.Draw(data, []int{0, 6})
		fmt.Fprint(v, out)
		//tm.Println(out)
		var buf bytes.Buffer
		tm.Output = bufio.NewWriter(&buf)
		//tm.Output = bufio.NewWriter(v)
		_ = buf
		tm.Flush()
		_ = v
		//buf.WriteTo(v)
		//buf.WriteTo(os.Stdout)
	*/

	maxX, _ := g.Size()
	if len(points) == 0 {
		points = chartCells(maxX, "bitcoin")
	}

	for i := range points {
		var s string
		for j := range points[i] {
			p := points[i][j]
			s = fmt.Sprintf("%s%c", s, p.Ch)
		}
		fmt.Fprintln(v, s)
	}
	return nil
}

func setTable(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, _ := g.View("table")
	t := table.New().SetWidth(maxX)

	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.AddCol("")
	t.HideColumHeaders = true

	var err error
	if len(coins) == 0 {
		coins, err = fetchData()
		if err != nil {
			return err
		}
	}
	for _, coin := range coins {
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
		_ = color7d
		t.AddRow(
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

	//t.SortAsc("Name").SortDesc("Age").Sort().Format().Fprint(v)
	t.Format().Fprint(v)
	//v.Editable = true
	//v.SetCursor(0, 2)
	v.Highlight = true
	v.SelBgColor = gocui.ColorCyan
	v.SelFgColor = gocui.ColorBlack
	v.Frame = false
	//v.Autoscroll = true

	//buffer := ln.Buffer()
	//buffer.

	// Sort by time
	// t.SortAscFn("Created", ltTime).Sort().Format().Fprint(v)
	return nil
}
