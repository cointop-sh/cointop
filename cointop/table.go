package cointop

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	apitypes "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/table"
)

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
		name := coin.Name
		lastchar := len(name)
		if lastchar > 20 {
			lastchar = 20
			name = fmt.Sprintf("%s...", name[0:17])
		}
		ct.table.AddRow(
			pad.Left(fmt.Sprint(coin.Rank), 4, " "),
			pad.Right("  "+string(name), 22, " "),
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
